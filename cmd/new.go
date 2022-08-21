/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"rtbl/tables"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a new result from a table",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		env_root := os.Getenv("RTBL_TABLE_ROOT")
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(root) == 0 {
			root = env_root
		}

		tablenames := args

		for _, tn := range tablenames {
			path := root + "/" + tn + ".tab"
			parsedTable, err := tables.Parse(path)
			if err != nil {
				fmt.Println(path, ":", err)
				return
			}
			r := parsedTable.Roll("Start")
			fmt.Println(r)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
