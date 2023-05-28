package agent

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mariadb-operator/agent/pkg/router"
	"github.com/spf13/cobra"
)

var (
	addr              string
	compressLevel     int
	rateLimitRequests int
	rateLimitDuration time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent",
	Long:  `🤖 Sidecar agent for MariaDB that co-operates with mariadb-operator`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := router.NewRouter(
			router.WithCompressLevel(compressLevel),
			router.WithRateLimit(rateLimitRequests, rateLimitDuration),
		)
		server := http.Server{
			Addr:    addr,
			Handler: router,
		}

		serverContext, stopServer := context.WithCancel(context.Background())

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-sig
			defer stopServer()

			log.Println("shutting down server")
			if err := server.Shutdown(context.Background()); err != nil {
				log.Fatalf("error shutting down server: %v", err)
			}
		}()

		log.Printf("server listening at %s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}

		<-serverContext.Done()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVar(&addr, "addr", ":5555", "The address that the HTTP server binds to")
	rootCmd.Flags().IntVar(&compressLevel, "compress-level", 5, "HTTP compression level")
	rootCmd.Flags().IntVar(&rateLimitRequests, "rate-limit-requests", 100, "Number of requests to be used as rate limit")
	rootCmd.Flags().DurationVar(&rateLimitDuration, "rate-limit-duration", 1*time.Minute, "Duration to be used as rate limit")
}
