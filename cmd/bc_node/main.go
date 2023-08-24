package main

import (
	"context"
	"os"
	"path"

	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	var DataDir string
	var rootCmd = &cobra.Command{
		Use: "bc_node",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configPath := path.Join(DataDir, "config.yaml")
			conf := utils.ParseConfig(configPath)
			ctx := context.WithValue(cmd.Context(), &utils.ConfigKey, conf)
			cmd.SetContext(ctx)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", "./", "data directory")

	rootCmd.AddCommand(NewCmdConfig(), NewCmdServe())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
