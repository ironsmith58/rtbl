/*
Copyright © 2022 Eric F. Wolcott <efwolcott@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rtbl",
	Short: "create random results from tables in files, similar to TableSmith",
	Long:  "This is a command line tool that will read TableSmith tables.  TableSmith was written by Bruce Gulke and has many more features that are present here.  See the tablesmith homepage; http://www.mythosa.net/p/tablesmith.html",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	env_root := os.Getenv("RTBL_ROOT")
	help := "root directory directory to start table lookup. \n"
	help += "Environment variable RTBL_ROOT may be set\n"
	help += "instead of using this flag.\n"
	if env_root != "" {
		help += "RTBL_ROOT=" + env_root
	}
	rootCmd.PersistentFlags().StringP("root", "r", "", help)
}
