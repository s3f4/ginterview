package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/s3f4/ginterview/api/handlers"
	"github.com/s3f4/ginterview/api/repository"
)

// Run starts the application
func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := flag.String("port", ":3001", " default port is 3001")
	flag.Parse()

	// c is used to graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		oscall := <-c
		log.Printf("system call: %+v\n", oscall)
		cancel()
	}()

	initConnections()
	defer mongoClient.Disconnect(ctx)

	// repository and handler initializations
	mongoRepository := repository.NewMongoRepository(mongoClient)
	mongoHandler := handlers.NewMongoHandler(mongoRepository)
	inMemoryRepository := repository.NewInMemoryRepository()
	inMemoryHandler := handlers.NewInMemoryHandler(inMemoryRepository)

	http.Handle("/mongo", mongoHandler)
	http.Handle("/in-memory", inMemoryHandler)

	// serve and wait
	if err := serve(ctx, port); err != nil {
		log.Printf("API serve failed : %v\n", err)
	}
}

// serve
func serve(ctx context.Context, port *string) (err error) {
	server := &http.Server{
		Addr: *port,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server listen err: %s\n", err)
		}
	}()

	log.Printf("API server started on port %s...\n", *port)
	<-ctx.Done()
	log.Printf("API server stopped.\n")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("API server shutdown failed: %s\n", err)
	}

	log.Printf("API server exited poroperly \n")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
