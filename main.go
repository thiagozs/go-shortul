package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/thiagozs/go-shorturl/initialize"
)

func main() {
	// Allow the user to specify a port via a command-line flag or environment variable
	portFlag := flag.String("port", "8080", "port to listen on")
	secretTokenFlag := flag.String("token", "5ecr3tT0k3n", "secret token for authentication")
	domainFlag := flag.String("domain", "localhost", "domain name")
	useHttpsFlag := flag.Bool("https", false, "use https")
	useLocalFlag := flag.Bool("local", true, "use local")

	flag.Parse()

	// Initialize structured logger with JSON handler and default options
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	init, err := initialize.NewInitialize(
		initialize.WithLogger(logger),
	)
	if err != nil {
		logger.Error("Failed to create initialize", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Bootstrap the application
	if err := init.Init(); err != nil {
		logger.Error("Failed to initialize", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Get configuration from environment variables or use default values
	cfg := init.GetParams().GetConfig()

	// Validation of configuration over command-line flags
	if port := cfg.GetPort(); port == "" {
		cfg.SetPort(*portFlag)
	}

	if token := cfg.GetToken(); token == "" {
		cfg.SetToken(*secretTokenFlag)
	}

	if domain := cfg.GetDomain(); domain == "" {
		cfg.SetDomain(*domainFlag)
	}

	if https := cfg.GetHTTPS(); !https {
		cfg.SetHTTPS(*useHttpsFlag)
	}

	if local := cfg.GetLocal(); !local {
		cfg.SetLocal(*useLocalFlag)
	}

	// Set the configuration
	init.SetConfigByFlags(cfg)

	if err := init.ReloadInit(); err != nil {
		logger.Error("Failed to reload configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Shorten the variable name
	api := init.GetParams().GetAPI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	infoStart := fmt.Sprintf("Starting server on host:%s port:%s...", cfg.GetHost(), cfg.GetPort())
	logger.Info(infoStart)

	// Register the server and endpoints
	api.RegisterEndPoints()
	api.RegisterServer()

	// Start the sever
	api.Start()

	<-sigChan

	logger.Info("Shutting down server...")
	// Gracefully shutdown the server
	api.Shutdown()

	os.Exit(0)

}
