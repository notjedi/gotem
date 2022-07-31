package cmd

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/config"
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

			cfg := config.New()

			cfg.Username = returnNonNil(username, cfg.Username)
			cfg.Password = returnNonNil(password, cfg.Password)
			cfg.Debug = returnNonNil(debugFlag, cfg.Debug)
			cfg.Host = returnNonNil(host, cfg.Host)
			cfg.Port = returnNonNil(port, cfg.Port)

			ctx, err := context.GetContext(cfg)
			if err != nil {
				log.Fatal(err)
			}

			p := tea.NewProgram(ui.New(ctx), tea.WithAltScreen())
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
