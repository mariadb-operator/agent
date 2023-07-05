package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/handler"
	"github.com/mariadb-operator/agent/pkg/kubeclientset"
	"github.com/mariadb-operator/agent/pkg/kubernetesauth"
	"github.com/mariadb-operator/agent/pkg/logger"
	"github.com/mariadb-operator/agent/pkg/router"
	"github.com/mariadb-operator/agent/pkg/server"
)

var (
	addr      string
	configDir string
	stateDir  string

	compressLevel              int
	rateLimitRequests          int
	rateLimitDuration          time.Duration
	kubernetesAuth             bool
	kubernetesTrustedName      string
	kubernetesTrustedNamespace string
	recoveryTimeout            time.Duration
	gracefulShutdownTimeout    time.Duration

	logLevel       string
	logTimeEncoder string
	logDev         bool
)

func main() {
	flag.StringVar(&addr, "addr", ":5555", "The address that the HTTP server binds to")
	flag.StringVar(&configDir, "config-dir", "/etc/mysql/mariadb.conf.d", "The directory that contains MariaDB configuration files")
	flag.StringVar(&stateDir, "state-dir", "/var/lib/mysql", "The directory that contains MariaDB state files")

	flag.IntVar(&compressLevel, "compress-level", 5, "HTTP compression level")
	flag.IntVar(&rateLimitRequests, "rate-limit-requests", 0, "Number of requests to be used as rate limit")
	flag.DurationVar(&rateLimitDuration, "rate-limit-duration", 0, "Duration to be used as rate limit")
	flag.BoolVar(&kubernetesAuth, "kubernetes-auth", false, "Enable Kubernetes authentication via the TokenReview API")
	flag.StringVar(&kubernetesTrustedName, "kubernetes-trusted-name", "", "Trusted Kubernetes ServiceAccount name to be verified")
	flag.StringVar(&kubernetesTrustedNamespace, "kubernetes-trusted-namespace", "", "Trusted Kubernetes ServiceAccount "+
		"namespace to be verified")
	flag.DurationVar(&recoveryTimeout, "recovery-timeout", 1*time.Minute, "Timeout to obtain sequence number "+
		"during the Galera cluster recovery process")
	flag.DurationVar(&gracefulShutdownTimeout, "graceful-shutdown-timeout", 5*time.Second, "Timeout to gracefully terminate "+
		"in-flight requests")

	flag.StringVar(&logLevel, "log-level", "info", "Log level to use, one of: "+
		"debug, info, warn, error, dpanic, panic, fatal.")
	flag.StringVar(&logTimeEncoder, "log-time-encoder", "epoch", "Log time encoder to use, one of: "+
		"epoch, millis, nano, iso8601, rfc3339 or rfc3339nano")
	flag.BoolVar(&logDev, "log-dev", false, "Enable development logs.")

	flag.Parse()

	logger, err := logger.NewLogger(
		logger.WithLogLevel(logLevel),
		logger.WithTimeEncoder(logTimeEncoder),
		logger.WithDevelopment(logDev),
	)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	clientset, err := kubeclientset.NewKubeclientSet()
	if err != nil {
		logger.Error(err, "error creating Kubernetes clientset")
		os.Exit(1)
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

	routerOpts := []router.Option{
		router.WithCompressLevel(compressLevel),
		router.WithRateLimit(rateLimitRequests, rateLimitDuration),
	}
	if kubernetesAuth && kubernetesTrustedName != "" && kubernetesTrustedNamespace != "" {
		routerOpts = append(routerOpts, router.WithKubernetesAuth(
			kubernetesAuth,
			&kubernetesauth.Trusted{
				ServiceAccountName:      kubernetesTrustedName,
				ServiceAccountNamespace: kubernetesTrustedNamespace,
			},
		))
	}
	router := router.NewRouter(
		handler,
		clientset,
		logger,
		routerOpts...,
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
