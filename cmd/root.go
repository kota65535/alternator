package cmd

import (
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
)

//go:embed root.tmpl
var rootUsage string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "alternator <command> [options]",
	Long: "SQL database schema management tool.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var (
	debug bool
	width int
)

const MinTermWidth = 80

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug flag")
	rootCmd.SetUsageTemplate(rootUsage)

	w, _, _ := term.GetSize(syscall.Stdin)
	if w < MinTermWidth {
		w = MinTermWidth
	}
	width = w
	if debug || os.Getenv("ALTERNATOR_DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
