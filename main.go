package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/handler"
	"github.com/mariadb-operator/agent/pkg/logger"
	"github.com/mariadb-operator/agent/pkg/router"
	"github.com/mariadb-operator/agent/pkg/server"
)

var (
	addr      string
	configDir string
	stateDir  string

	compressLevel           int
	rateLimitRequests       int
	rateLimitDuration       time.Duration
	gracefulShutdownTimeout time.Duration

	logLevel       string
	logTimeEncoder string
	logDev         bool

	recoveryTimeout time.Duration
)

func main() {
	flag.StringVar(&addr, "addr", ":5555", "The address that the HTTP server binds to")
	flag.StringVar(&configDir, "config-dir", "/etc/mysql/mariadb.conf.d", "The directory that contains MariaDB configuration files")
	flag.StringVar(&stateDir, "state-dir", "/var/lib/mysql", "The directory that contains MariaDB state files")

	flag.IntVar(&compressLevel, "compress-level", 5, "HTTP compression level")
	flag.IntVar(&rateLimitRequests, "rate-limit-requests", 0, "Number of requests to be used as rate limit")
	flag.DurationVar(&rateLimitDuration, "rate-limit-duration", 0, "Duration to be used as rate limit")
	flag.DurationVar(&gracefulShutdownTimeout, "graceful-shutdown-timeout", 5*time.Second, "Timeout to gracefully terminate "+
		"in-flight requests")

	flag.StringVar(&logLevel, "log-level", "info", "Log level to use, one of: "+
		"debug, info, warn, error, dpanic, panic, fatal.")
	flag.StringVar(&logTimeEncoder, "log-time-encoder", "epoch", "Log time encoder to use, one of: "+
		"epoch, millis, nano, iso8601, rfc3339 or rfc3339nano")
	flag.BoolVar(&logDev, "log-dev", false, "Enable development logs.")

	flag.DurationVar(&recoveryTimeout, "recovery-timeout", 1*time.Minute, "Timeout to obtain sequence number "+
		"during the Galera cluster recovery process")
	flag.Parse()

	logger, err := logger.NewLogger(
		logger.WithLogLevel(logLevel),
		logger.WithTimeEncoder(logTimeEncoder),
		logger.WithDevelopment(logDev),
	)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	fileManager, err := filemanager.NewFileManager(configDir, stateDir)
	if err != nil {
		logger.Error(err, "error creating file manager")
		os.Exit(1)
	}

	handlerLogger := logger.WithName("handler")
	handler := handler.NewHandler(
		fileManager,
		&handlerLogger,
		handler.WithRecoveryTimeout(recoveryTimeout),
	)

	router := router.NewRouter(
		handler,
		router.WithCompressLevel(compressLevel),
		router.WithRateLimit(rateLimitRequests, rateLimitDuration),
	)

	serverLogger := logger.WithName("server")
	server := server.NewServer(
		addr,
		router,
		&serverLogger,
		server.WithGracefulShutdownTimeout(gracefulShutdownTimeout),
	)
	if err := server.Start(context.Background()); err != nil {
		logger.Error(err, "server error")
		os.Exit(1)
	}
}
