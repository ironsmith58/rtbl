package tables

import (
	"fmt"
	"rtbl/stringsext"
	"rtbl/tfs"
	"strconv"
	"strings"
)

// Parsing States
const (
	NO_GROUP = iota
	COLON_GROUP
	SEMI_GROUP
)

func parseColonItem(line string) (int, int, string, error) {
	/* Parse these;
	 * 1-2,Orc
	 * 3,Skeleton
	 * 4-7,Archdaemon
	 * 8,[Goblins.Start]
	 * 9,[Color] Spirit
	 */
	var err error
	var min int64
	var max int64

	var irest int

	idx := strings.Index(line, ",")
	if idx == -1 { // comma not found, lets check for a tab
		idx = strings.Index(line, "\t")
		if idx == -1 {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: no delimiter between range and text; %s", line)
		}
	}
	fields := []string{line[:idx], line[idx+1:]}
	// fields[0] is the number range and needs to be parsed
	//   is there a dash??
	nums := strings.Split(fields[0], "-")
	irest = len(fields[0]) + 1
	if len(nums) == 1 {
		// single number, not a range
		min, err = strconv.ParseInt(nums[0], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: single probability is not a number; %s", line)
		}
	} else if len(nums) == 2 {
		// probability range
		min, err = strconv.ParseInt(nums[0], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: min probability of range is not a number; %s", line)
		}
		max, err = strconv.ParseInt(nums[1], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: max probability of range is not a number; %s", line)
		}
	}
	if max != 0 && max < min {
		return int(0), int(0), "",
			fmt.Errorf("Colon Group: max of probability range is less than min; %s", line)
	}
	// fields[1] is the text

	return int(min), int(max), line[irest:], nil
}

func parseSemiItem(line string) (int, string, error) {
	var err error
	var min int64

	t := strings.TrimSpace(line)
	if len(t) == 0 {
		return 0, "", fmt.Errorf("blank line")
	}
	idx := strings.Index(line, ",")
	if idx == -1 {
		idx = strings.Index(line, "\t")
		if idx == -1 {
			return int(0), "",
				fmt.Errorf("Semi Group: no delimiter between range and text; %s", line)
		}
	}
	fields := []string{line[:idx], line[idx+1:]}
	min, err = strconv.ParseInt(fields[0], 10, 0)
	if err != nil {
		return int(0), "",
			fmt.Errorf("Semi Group: probability is not a number; %s", line)
	}
	return int(min), fields[1], nil
}

func parseVariableDeclaration(line string) (string, string, error) {
	// Variable Format: %VariableName%,x
	words := strings.Split(line, ",")
	words[0] = strings.Replace(words[0], "%", "", -1)
	value := ""
	if len(words) == 2 {
		value = words[1]
	}
	return words[0], value, nil
}

func parseVariableAssignment(line string) (string, string, string, error) {
	// Variable Format: |VariableName?x|
	// ? is the opcode, [+-*/\><&=]
	if line[0] != '|' && line[len(line)-1] != '|' {
		return "", "", "", fmt.Errorf("Variable assignment needs delimiters %s", line)
	}
	line = line[1 : len(line)-1]
	idx := strings.IndexAny(line, "+-*/\\><&=")
	if idx == -1 {
		return "", "", "", fmt.Errorf("No OpCode in %s", line)
	}
	op := string(line[idx])
	name := line[:idx]
	val := line[idx+1:]

	return name, val, op, nil
}

func Parse(tableName string) (*Table, error) {
	errorFmt := "Error: %s, line %d\n"
	tableName = strings.ToLower(tableName)

	//check the table registry to see if the table has
	// already been loaded
	loadedTable, ok := TableRegistry[tableName]
	if !ok {
		return nil, fmt.Errorf("No table found")
	}
	if loadedTable.table != nil {
		return loadedTable.table, nil
	}

	// if not already loaded, lets load it and parse  it
	content, err := tfs.ReadFile(loadedTable.path)
	if err != nil {
		return nil, err
	}

	// Let see if we have already parse this table
	table := NewTable(tableName)
	table.Size = len(content)
	table.Path = loadedTable.path // TODO refactor path away. dont need here and in loadedTable
	var group *Group
	var state int

	//  .-------------.
	//  | Parse Lines |
	//  '-------------'
	for lineno, line := range content {
		// Check for comment and strip it
		idx := strings.Index(line, "#")
		if idx > -1 {
			line = line[:idx]
		}
		// trim the line
		line = strings.TrimRight(line, "\t\r\n")
		line = strings.TrimLeft(line, " \t")
		// if there is nothing to parse go to next line
		// blank line also closes any previous group parsing
		if len(line) == 0 {
			//state = NO_GROUP
			//if group != nil {
			//	table.AddGroup(group)
			//	group = nil
			//}
			continue
		}
		// now we parse data lines
		if line[0] == '/' { // this is a special directive
			directive := stringsext.First(line[1:])
			// this is tructures so that, in time, each case
			// can be implemented in someway
			switch directive {
			case "BackColor":
				fmt.Printf("Unknown directive, ignoring %s, line %d\n", line, lineno)
			case "Background":
				fmt.Printf("Unknown directive, ignoring %s, line %d\n", line, lineno)
			case "OutputFooter":
				table.Footer = stringsext.Rest(line)
			case "OutputHeader":
				table.Header = stringsext.Rest(line)
			case "OverrideRolls":
				fmt.Printf("Unknown directive, ignoring %s, line %d\n", line, lineno)
			case "Stylesheet":
				fmt.Printf("Unknown directive, ignoring %s, line %d\n", line, lineno)
			default:
				fmt.Printf("Unknown directive, ignoring %s, line %d\n", line, lineno)
			}
		} else if line[0] == ':' {
			state = COLON_GROUP
			// Save any previous group
			if group != nil {
				table.AddGroup(group)
				group = nil
			}
			group = NewGroup(line)
		} else if line[0] == ';' {
			state = SEMI_GROUP
			// Save any previous group
			if group != nil {
				table.AddGroup(group)
				group = nil
			}
			group = NewGroup(line)
		} else if state == COLON_GROUP {
			if line[0] == '<' {
				// is it a prefix?
				group.Prefix = line[1:]
			} else if line[0] == '>' {
				// is it a suffix?
				group.Suffix = line[1:]
			} else if line[0] == '_' {
				// is it continuation
				err := group.AppendLastItem("<br>" + line[1:])
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				}
			} else {
				start, end, text, err := parseColonItem(line)
				// for single number 'ranges' end will be == 0
				// to complete teh range we wil set it equal to
				// start
				if end == 0 {
					end = start
				}
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				}
				group.AddItem(start, end, text)
			}
		} else if state == SEMI_GROUP {
			if line[0] == '<' {
				// is it a prefix?
				group.Prefix = line[1:]
			} else if line[0] == '>' {
				// is it a suffix?
				group.Suffix = line[1:]
			} else if line[0] == '_' {
				// is it continuation
				err := group.AppendLastItem("<br>" + line[1:])
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				}
			} else {
				num, text, err := parseSemiItem(line)
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				} else {
					group.AddItem(num, 0, text)
				}
			}
		} else if line[0] == '%' {
			// Variable Format: %VariableName%,x
			name, value, err := parseVariableDeclaration(line)
			if err != nil {
				err = fmt.Errorf(errorFmt, err, lineno)
				return nil, err
			}

			err = table.AddVariable(name, value)
			if err != nil {
				err = fmt.Errorf(errorFmt, err, lineno)
				return nil, err
			}
		} else if line[0] == '|' {
			// Variable Format: |VariableName?x|
			// ? is the opcode, [+-*/\><&=]

			name, newstr, op, err := parseVariableAssignment(line)
			if err != nil {
				err = fmt.Errorf(errorFmt, err, lineno)
				return nil, err
			}

			oldstr, ok := table.GetVariable(name)
			var oldval float64
			var newval float64
			if ok {
				oldval, _ = strconv.ParseFloat(oldstr, 32)
			} else {
				oldval = 0.0
			}
			newval, _ = strconv.ParseFloat(newstr, 64)
			// process the op to create the new string value
			// that wiil be stored under the variable 'name'
			switch op {
			case "+":
				newstr = fmt.Sprintf("%f", oldval+newval)
			case "-":
				newstr = fmt.Sprintf("%f", oldval-newval)
			case "*":
				newstr = fmt.Sprintf("%f", oldval*newval)
			case "/":
				newstr = fmt.Sprintf("%f", oldval/newval)
			case "\\":
				newstr = fmt.Sprintf("%d", int(oldval/newval))
			case ">":
				if newval <= oldval {
					newstr = oldstr // re-assign old value to variable
				}
			case "<":
				if newval >= oldval {
					newstr = oldstr
				}
			case "&":
				// string catenation
				newstr = oldstr + newstr
			case "=":
				// noop, this will just assign newstr to the variable
			default:
				return nil, fmt.Errorf("Unknown OpCode in %s", line)
			}
			// will add or update variable
			table.AddVariable(name, newstr)
		}
	}
	if group != nil {
		table.AddGroup(group)
		group = nil
	}
	loadedTable.table = table
	return table, nil
}
