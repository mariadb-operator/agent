package agent

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent",
	Long:  `🤖 Sidecar agent for MariaDB that co-operates with mariadb-operator`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello agent!")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
