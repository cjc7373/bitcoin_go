package main

import (
	"fmt"

	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/spf13/cobra"
)

func NewCmdConfig() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "operate configs",
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show config",
		Run: func(cmd *cobra.Command, args []string) {
			config := utils.GetConfigFromContext(cmd.Context())
			fmt.Println("DB Path ", config.DBPath)
			fmt.Println()
			fmt.Println("Wallets:")
			for k, v := range config.Wallets {
				fmt.Printf("name: %v\nprivate key: %v\n\n", k, v)
			}
		},
	}

	configCmd.AddCommand(showCmd)

	return configCmd
}
