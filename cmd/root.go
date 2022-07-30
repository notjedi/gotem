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

	logfilePath = "debug.log"

	rootCmd = &cobra.Command{
		Use:     "gotem",
		Short:   "A glamourous TUI for the BitTorrent client Transmission.",
		Version: "0.1.0",
		// Arg: TODO
		Run: func(cmd *cobra.Command, args []string) {

			config := c.New()

			if debugFlag {
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

			config.Username = returnNonNil(username, config.Username)
			config.Password = returnNonNil(password, config.Password)
			config.Host = returnNonNil(host, config.Host)
			config.Port = returnNonNil(port, config.Port)

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

func returnNonNil[T comparable](val1, val2 T) T {
	var temp T
	if val1 == temp {
		return val2
	}
	return val1
}

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

	rootCmd.PersistentFlags().StringVar(&host, "host", "", "Host address")
	rootCmd.PersistentFlags().Uint16Var(&port, "port", 0, "RPC port")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username to connect to the host")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password to connect to the host")
	// should i make this a string so it accepts a path for the debug file?
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debugging")
	// rootCmd.PersistentFlags().String("path", "/transmission/rpc", "Path to the RPC")
}
