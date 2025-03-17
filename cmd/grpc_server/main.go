package main

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/app"
	"log"
)

var configPath string

//func init() {
//	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
//}

func main() {
	//flag.Parse()
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}

}
