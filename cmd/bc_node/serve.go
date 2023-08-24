package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdServe() *cobra.Command {
	var genesis bool
	var connectTo string
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run a node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("serve")
		},
	}
	serveCmd.Flags().BoolVar(&genesis, "genesis", false, "to create the genesis block")
	serveCmd.Flags().StringVar(&connectTo, "connect-to", "", "connect to a node as the first neighbour (either this flag or --genesis should be set)")
	serveCmd.MarkFlagsOneRequired("genesis", "connect-to")
	return serveCmd
}
