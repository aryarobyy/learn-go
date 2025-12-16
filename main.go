package main

import (
	"auth/app"
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("ENV LOAD ERROR:", err)
	}

	app := app.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
