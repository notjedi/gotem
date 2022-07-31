package cmd

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/notjedi/gotem/internal/config"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui"
	"github.com/spf13/cobra"
)

var (
	debugFlag bool
	username  string
	password  string
	host      string
	port      uint16

	rootCmd = &cobra.Command{
		Use:     "gotem",
		Short:   "A glamourous TUI for the BitTorrent client Transmission.",
		Version: "0.1.0",
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			config := c.New()

			config.Username = returnNonNil(username, config.Username)
			config.Password = returnNonNil(password, config.Password)
			config.Debug = returnNonNil(debugFlag, config.Debug)
			config.Host = returnNonNil(host, config.Host)
			config.Port = returnNonNil(port, config.Port)

			client, err := context.GetClient(config)
			if err != nil {
				log.Fatal(err)
			}

			p := tea.NewProgram(ui.New(client), tea.WithAltScreen())
			if err := p.Start(); err != nil {
				log.Fatal(err)
			}

		},
	}
)

func returnNonNil[T comparable](val1, val2 T) T {
	var zeroVal T
	if val1 == zeroVal {
		return val2
	}
	return val1
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username to connect to the host")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password to connect to the host")
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debugging")
	// rootCmd.PersistentFlags().String("path", "/transmission/rpc", "Path to the RPC")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "Host address")
	rootCmd.PersistentFlags().Uint16Var(&port, "port", 0, "RPC port")
}
