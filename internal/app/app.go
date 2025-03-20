package app

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/closer"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"github.com/Muvi7z/chat-auth-s/internal/interceptor"
	_ "github.com/Muvi7z/chat-auth-s/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

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
	wg.Add(3)

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

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
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
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)
	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserServer(ctx))

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
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Println("starting http server on " + a.serviceProvider.SwaggerConfig().Address())

	err := a.httpServer.ListenAndServe()
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

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file %s", path)

		statilFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Opening swagger file %s", path)

		file, err := statilFs.Open(path)
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
