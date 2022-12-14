/*
Copyright © 2022 Eric F. Wolcott <efwolcott@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"rtbl/tables"
	"strconv"
	"strings"

	"golang.org/x/term"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/spf13/cobra"
	"jaytaylor.com/html2text"
)

type TableCall struct {
	repeat int
	table  string
	group  string
	args   []string
}

func parseCall(s string) (TableCall, error) {
	// [Table.Group(a1,a2)]:n
	var tableCall TableCall

	// Get Table, text within brackets, left of the dot(.) (required)
	// if no brackets, then entire string is parsed
	i := strings.Index(s, "[")
	e := strings.Index(s, "]")
	if i != -1 && e != -1 {
		if e <= i {
			return TableCall{}, fmt.Errorf("Mismatched brackets in group call")
		} else {
			if i == -1 {
				return TableCall{}, fmt.Errorf("No brackets in group call")
			} else {
				// lets get rid of the brackets []
				s = strings.Replace(s[i:], "]", "", 1)
			}
		}
	}

	// Get Starting Group, text within brackets right of the dot(.) (optional)
	w := strings.Split(s, ".")
	tableCall.table = w[0]
	tableCall.group = "Start"
	if len(w) == 2 {
		tableCall.group = w[1]
		// Get Arguments, text within parens (optional)
		i = strings.Index(w[1], "(")
		if i != -1 {
			e := strings.Index(w[1], ")")
			if e <= i {
				return TableCall{}, fmt.Errorf("Mismatched braces in group call")
			}
			tableCall.args = strings.Split(w[1][i+1:e], ",")
			tableCall.group = w[1][:i] // group name stops at open paren [(]
		}
	}

	// get repeat Count, number after colon(:) (optional)
	i = strings.LastIndex(s, ":")
	if i != -1 {
		n, err := strconv.ParseInt(s[i+1:], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		tableCall.repeat = int(n)
		s = s[:i]
	} else {
		tableCall.repeat = 1
	}

	return tableCall, nil
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a new result from a table",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		env_root := os.Getenv("RTBL_ROOT")
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(root) == 0 {
			root = env_root
		}

		// Load all tables to table registry
		var rootpath string
		if root != "" {
			// humans might enter the path with a wildcard that expands to
			// contain the Tables sub-dir
			if !strings.HasSuffix(root, "/Tables") && !strings.HasSuffix(root, "/Tables/") {
				rootpath = root + "/Tables/"
			} else {
				rootpath = root
			}
		} else {
			rootpath = "./Tables"
		}
		err = tables.LoadAllTables(rootpath)
		//paths, err := tables.FindTables(rootpath)
		//tableList := tables.NewTableList(paths)
		tablenames := args
		for _, tn := range tablenames {

			tc, err := parseCall(tn)

			parsedTable, err := tables.Parse(tc.table)
			if err != nil {
				fmt.Println(tn, ":", err)
				return
			}
			// Roll on Table
			html := parsedTable.Roll(tc.group)
			// Handle OutputHeader and OutputFooter directive
			if len(parsedTable.Header) > 0 {
				html = parsedTable.Header + html
			}
			if len(parsedTable.Footer) > 0 {
				html = html + parsedTable.Footer
			}
			xport, err := cmd.Flags().GetString("export")
			switch xport {
			case "html":
				fmt.Println(html)
			case "text":
				colWidth := 70
				widthOpt, err := cmd.Flags().GetInt("width")
				if widthOpt > 0 {
					colWidth = widthOpt
				} else {
					width, _, err := term.GetSize(0)
					if err != nil {
						colWidth = 72
					} else {
						// never make output too small]
						if colWidth > 23 {
							colWidth = int(float64(width) * .75)
						} else {
							colWidth = width
						}
					}
				}

				pto := html2text.NewPrettyTablesOptions()
				pto.ColWidth = colWidth
				text, err := html2text.FromString(html,
					html2text.Options{
						PrettyTables:        true,
						PrettyTablesOptions: pto,
					})
				if err != nil {
					panic(err)
				}
				fmt.Println(text)
			case "md":
				converter := md.NewConverter("", true, nil)
				markdown, err := converter.ConvertString(html)
				if err != nil {
					panic(err)
				}
				fmt.Println(markdown)
			default:
				fmt.Printf("Export format is unsupported; %s\n", xport)
			}
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
	newCmd.Flags().StringP("export", "x", "text", "output format (text,html,md)")
	newCmd.Flags().IntP("width", "w", 0, "width of text output")
}
