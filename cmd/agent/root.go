package agent

import (
	"context"
	"log"
	"time"

	"github.com/mariadb-operator/agent/pkg/router"
	"github.com/mariadb-operator/agent/pkg/server"
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
		server := server.NewServer(addr, router)
		if err := server.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
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
