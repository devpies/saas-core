package main

import (
	"context"
	"fmt"
	"github.com/ivorscott/devpie-client-backend-go/cmd/api/internal/handlers"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivorscott/devpie-client-backend-go/internal/platform/conf"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	infolog := log.New(os.Stdout, "mic-projects: ", log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// App Starting

	infolog.Printf("main : Started")
	defer infolog.Println("main : Completed")

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Port             string        `conf:"default::4000"`
			Debug            string        `conf:"default:localhost:6060"`
			Production       bool          `conf:"default:false"`
			ReadTimeout      time.Duration `conf:"default:5s"`
			WriteTimeout     time.Duration `conf:"default:5s"`
			ShutdownTimeout  time.Duration `conf:"default:5s"`
			FrontendAddress  string        `conf:"default:https://localhost:3000"`
			AuthDomain       string        `conf:"default:none,noprint"`
			AuthAudience     string        `conf:"default:none,noprint"`
		}
		DB struct {
			User       string `conf:"default:postgres,noprint"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost,noprint"`
			Name       string `conf:"default:postgres,noprint"`
			DisableTLS bool   `conf:"default:false,noprint"`
		}
		//Nats struct {
		//	Url string `conf:"default:nats://"`
		//	ClientId string `conf:"default:client-id"`
		//	ClusterId string `conf:"default:cluster-id"`
		//}
	}

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	infolog.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Enabled Profiler

	go func() {
		log.Printf("main: Debug service listening on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, nil)
		if err != nil {
			log.Printf("main: Debug service listening on %s", cfg.Web.Debug)
		}
	}()

	// =========================================================================
	// Start Database

	repo, err := database.NewRepository(database.Config{
		User:       cfg.DB.User,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer repo.Close()

	// =========================================================================
	// Start API Service

	// Make a channel to listen for shutdown signal from the OS.
	shutdown := make(chan os.Signal, 1)

	api := http.Server{
		Addr: cfg.Web.Port,
		Handler: handlers.API(shutdown, repo, infolog, cfg.Web.FrontendAddress, cfg.Web.AuthAudience,
			cfg.Web.AuthDomain),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	//
	//// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case sig := <-shutdown:
		log.Println("main : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}