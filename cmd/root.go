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
	logfilePath = "debug.log"

	rootCmd = &cobra.Command{
		Use:     "gotem",
		Short:   "A glamourous TUI for the BitTorrent client Transmission.",
		Version: "0.1.0",
		// Arg:
		Run: func(cmd *cobra.Command, args []string) {

			if debugFlag, _ := cmd.PersistentFlags().GetBool("debug"); debugFlag {
				logfile, err := tea.LogToFile(logfilePath, "gotem ")
				if err != nil {
					exit(err)
				}
				defer func() {
					if err := logfile.Close(); err != nil {
						log.Fatal(err)
					}
				}()
			}

			config := c.New()
			client, err := context.GetClient(config)

			if err != nil {
				exit(err)
			}

			p := tea.NewProgram(ui.New(client), tea.WithAltScreen())
			if err := p.Start(); err != nil {
				log.Fatal(err)
			}

		},
	}
)

func exit(err error) {
	log.Fatal(err)
	os.Exit(1)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().String("host", "localhost", "Host address")
	rootCmd.PersistentFlags().Uint16("port", 9091, "RPC port")
	rootCmd.PersistentFlags().String("username", "", "Username to connect to the host")
	rootCmd.PersistentFlags().String("password", "", "Password to connect to the host")
	// should i make this a string so it accepts a path for the debug file?
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debugging")
	// rootCmd.PersistentFlags().String("path", "/transmission/rpc", "Path to the RPC")
}
