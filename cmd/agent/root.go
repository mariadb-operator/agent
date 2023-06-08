package agent

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/handler"
	"github.com/mariadb-operator/agent/pkg/handler/recovery"
	"github.com/mariadb-operator/agent/pkg/logger"
	"github.com/mariadb-operator/agent/pkg/router"
	"github.com/mariadb-operator/agent/pkg/server"
	"github.com/spf13/cobra"
)

var (
	addr      string
	configDir string
	stateDir  string

	compressLevel     int
	rateLimitRequests int
	rateLimitDuration time.Duration

	logLevel       string
	logTimeEncoder string
	logDev         bool

	recoveryRetries   int
	recoveryRetryWait time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent",
	Long:  `🤖 Sidecar agent for MariaDB that co-operates with mariadb-operator`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
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
		handler := handler.NewHandler(fileManager, &handlerLogger,
			handler.WithRecoveryOptions(
				recovery.WithRecovery(&recovery.RecoveryOptions{
					Retries:   recoveryRetries,
					WaitRetry: recoveryRetryWait,
				}),
			),
		)

		router := router.NewRouter(
			handler,
			router.WithCompressLevel(compressLevel),
			router.WithRateLimit(rateLimitRequests, rateLimitDuration),
		)

		serverLogger := logger.WithName("server")
		server := server.NewServer(addr, router, &serverLogger)
		if err := server.Start(context.Background()); err != nil {
			logger.Error(err, "error starting server")
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVar(&addr, "addr", ":5555", "The address that the HTTP server binds to")
	rootCmd.Flags().StringVar(&configDir, "config-dir", "/etc/mysql/mariadb.conf.d", "The directory that contains MariaDB configuration files")
	rootCmd.Flags().StringVar(&stateDir, "state-dir", "/var/lib/mysql", "The directory that contains MariaDB state files")

	rootCmd.Flags().IntVar(&compressLevel, "compress-level", 5, "HTTP compression level")
	rootCmd.Flags().IntVar(&rateLimitRequests, "rate-limit-requests", 100, "Number of requests to be used as rate limit")
	rootCmd.Flags().DurationVar(&rateLimitDuration, "rate-limit-duration", 1*time.Minute, "Duration to be used as rate limit")

	rootCmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level to use, one of: "+
		"debug, info, warn, error, dpanic, panic, fatal.")
	rootCmd.Flags().StringVar(&logTimeEncoder, "log-time-encoder", "epoch", "Log time encoder to use, one of: "+
		"epoch, millis, nano, iso8601, rfc3339 or rfc3339nano")
	rootCmd.Flags().BoolVar(&logDev, "log-dev", false, "Enable development logs.")

	rootCmd.Flags().IntVar(&recoveryRetries, "recovery-retries", 10, "Maximum number of attempts "+
		"to recover the Galera cluster")
	rootCmd.Flags().DurationVar(&recoveryRetryWait, "recovery-retry-wait", 3*time.Second, "Time to wait between "+
		"Galera cluster recover attempts ")
}
