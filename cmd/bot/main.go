package main

import (
	"context"
	"errors"
	"log"
	"merger-adapter/internal/app"
	"merger-adapter/internal/common/msgs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println(msgs.ServerStarting)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	cfg := initConfig()
	log.Println(msgs.ConfigInitialized)
	ctx, cancel := context.WithCancel(context.Background())
	go runApplication(ctx, cfg)

	gracefulShutdown(cancel)
}

func runApplication(ctx context.Context, cfg *app.Config) {
	log.Println(msgs.ApplicationStart)
	err := app.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("application: %s", err)
	}
	os.Exit(0)
}

func initConfig() *app.Config {
	cfgFs := app.InitFlagSet()

	cfg, err := cfgFs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("config FlagSet initialization: %s", err)
		if errors.Is(err, app.WrongArgumentError) {
			cfgFs.Usage()
		}
		os.Exit(2)
	}
	return cfg
}

var gracefulShutdownTimeout = 2 * time.Second

func gracefulShutdown(cancel context.CancelFunc) {
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("%s signal was received", <-quit)
	log.Printf("after %v seconds, the program will force exit", gracefulShutdownTimeout.Seconds())
	cancel()
	time.Sleep(gracefulShutdownTimeout)
	os.Exit(0)
}
