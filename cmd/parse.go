/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"rtbl/tables"
	"strings"

	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse a table, reporting errors",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		xport, err := cmd.Flags().GetBool("export")
		if err != nil {
			fmt.Println(err)
			return
		}
		env_root := os.Getenv("RTBL_TABLE_ROOT")
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(root) == 0 {
			root = env_root
		}
		//if len(root) > 0 && root[0] != '/' {
		//		root = env_root + "/" + root
		//}
		if len(args) > 0 {
			for _, name := range args {
				path := root + "/" + name
				if !strings.HasSuffix(path, ".tab") {
					path = path + ".tab"
				}
				parsedTable, err := tables.Parse(path)
				if err != nil {
					fmt.Println(path, ":", err)
				} else {
					fmt.Println(path, " No errors")
					if xport {
						// Marshalling the structure
						// For now ignoring error
						// but you should handle
						// the error in above function
						jsonF, err := json.MarshalIndent(parsedTable, "", "  ")
						if err != nil {
							fmt.Println("Internal Error:", err)
						}
						// typecasting byte array to string
						fmt.Println(string(jsonF))
					}
				}
			}
		} else {
			parsedTable, err := tables.Parse(root)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(root, " No errors")
			if xport {
				// Marshalling the structure
				// For now ignoring error
				// but you should handle
				// the error in above function
				jsonF, err := json.MarshalIndent(parsedTable, "", "  ")
				if err != nil {
					fmt.Println("Internal Error:", err)
				}

				// typecasting byte array to string
				fmt.Println(string(jsonF))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	parseCmd.Flags().BoolP("export", "x", false, "print table as json")
}
