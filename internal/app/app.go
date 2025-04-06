package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/closer"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"github.com/Muvi7z/chat-auth-s/internal/interceptor"
	"github.com/Muvi7z/chat-auth-s/internal/logger"
	"github.com/Muvi7z/chat-auth-s/internal/metrics"
	"github.com/Muvi7z/chat-auth-s/internal/rate_limiter"
	"github.com/Muvi7z/chat-auth-s/internal/tracing"
	_ "github.com/Muvi7z/chat-auth-s/statik"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var logLevel = flag.String("l", "info", "Log level")

const rateLimit = 10

const serviceName = "chat-auth-service"

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil

}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to start swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runPrometheusServer()
		if err != nil {
			log.Fatalf("failed to start prometheus server: %v", err)
		}
	}()
	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initLogger,
		a.initSwagger,
		a.InitTracer,
		metrics.Init,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) InitTracer(_ context.Context) error {
	err := tracing.Init(serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	rateLimiter := rate_limiter.NewTokenBucketLimiter(ctx, rateLimit, time.Second)

	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        "my-service",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Info(fmt.Sprintf("Circuit Breaker: %s, changed from %s, to %s", name, from.String(), to.String()))
		},
	})

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
				interceptor.NewRateLimiterInterceptor(rateLimiter).Unary,
				interceptor.NewCircuitBreakerInterceptor(cb).Unary,
				interceptor.MetricsInterceptor,
				interceptor.SeverTracingInterceptor,
				interceptor.ErrorCodesInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)
	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserServer(ctx))

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	logger.Init(getCore(zap.NewAtomicLevelAt(level)))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwagger(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Println("starting swagger server on " + a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runHTTPServer() error {
	log.Println("starting http server on " + a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := http.Server{
		Addr:    "localhost:2112", //TODO брать из конфига
		Handler: mux,
	}

	log.Printf("Prometheus server is running on :2112")

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
	})

	prodConfig := zap.NewProductionEncoderConfig()
	prodConfig.TimeKey = "timestamp"
	prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentConfig := zap.NewDevelopmentEncoderConfig()
	developmentConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentConfig)
	fileEncoder := zapcore.NewJSONEncoder(prodConfig)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Opening swagger file %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file %s", path)

		contents, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file %s", path)

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(contents)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file %s", path)

	}
}
