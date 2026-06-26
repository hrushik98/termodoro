// Termodoro is a terminal-based Pomodoro timer.
package main

import (
	"fmt"
	"os"

	"termodoro/app"

	tea "github.com/charmbracelet/bubbletea"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

var version = "0.1.1"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "termodoro",
	Short: "Termodoro is a terminal Pomodoro timer",
	Run: func(cmd *cobra.Command, args []string) {
		runTUI()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Termodoro",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var genManCmd = &cobra.Command{
	Use:   "genman",
	Short: "Use this command to generate the man page",
	Run: func(cmd *cobra.Command, args []string) {
		manPage, err := mcobra.NewManPage(1, rootCmd)
		if err != nil {
			panic(err)
		}

		manPage = manPage.WithSection("Copyright", "(C) 2026 Sairash Sharma Gautam. \n"+"Released under MIT license.")
		fmt.Println(manPage.Build(roff.NewDocument()))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genManCmd)
}

func runTUI() {
	p := tea.NewProgram(
		app.NewModel(),
		tea.WithAltScreen(),
	)

	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if model, ok := m.(app.Model); ok && model.Quitting {
		if model.Err != nil {
			fmt.Printf("Error Occurred: %s\n", model.Err.Error())
			return
		}
		fmt.Println(app.EndInfo)
	}
}
