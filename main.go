package main

import (
	"context"
	"log"
	"server/http"
	"server/service"
	"time"
)

func main() {
	_, cancel := start()
	defer shutdown(cancel)
	service.WaitShutdown()
}

func start() (ctx context.Context, cancel context.CancelFunc) {
	// This is the main context for the service. When it is canceled it means the service is going down.
	// All the tasks must be canceled
	ctx, cancel = context.WithCancel(context.Background())
	http.Start()
	return
}

func shutdown(cancel context.CancelFunc) {
	cancel()
	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelTimeout()
	doneHTTP := http.Shutdown(ctx)
	err := service.WaitUntilIsDoneOrCanceled(ctx, doneHTTP)
	if err != nil {
		log.Printf("service stopped by timeout %s\n", err)
	}
	log.Println("Shutting down service...")
}
