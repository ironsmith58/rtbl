package tables

import (
	"fmt"
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

	var fields []string
	var irest int

	fields = strings.Split(line, ",")
	if len(fields) < 2 {
		fields = strings.Split(line, "\t")
		if len(fields) != 2 {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: no delimiter between range and text; %s", line)
		}
	}
	// fields[0] is the number range and needs to be parsed
	//   is there a dash??
	nums := strings.Split(fields[0], "-")
	irest = len(fields[0]) + 1
	if len(nums) == 1 {
		// single number, not a range
		min, err = strconv.ParseInt(nums[0], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: probablity is not a number; %s", line)
		}
	} else if len(nums) == 2 {
		// probability range
		min, err = strconv.ParseInt(nums[0], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: min probablity is not a number; %s", line)
		}
		max, err = strconv.ParseInt(nums[0], 10, 0)
		if err != nil {
			return int(0), int(0), "",
				fmt.Errorf("Colon Group: max probablity is not a number; %s", line)
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

	fields := strings.Split(line, ",")
	if len(fields) != 2 {
		fields = strings.Split(line, "\t")
		if len(fields) != 2 {
			return int(0), "",
				fmt.Errorf("Colon Group: no delimter between range and text; %s", line)
		}
	}
	min, err = strconv.ParseInt(fields[0], 10, 0)
	if err != nil {
		return int(0), "",
			fmt.Errorf("Colon Group: probablity is not a number; %s", line)
	}
	return int(min), fields[1], nil
}

func parseVariableDeclaration(line string) (string, string, error) {
	// Variable Format: %VariableName%,x
	words := strings.Split(line, ",")
	if len(words) != 2 {
		return "", "", fmt.Errorf("Variable malformed %s", line)
	}
	words[0] = strings.Replace(words[0], "%", "", -1)
	return words[0], words[1], nil
}
func parseVariableAssignment(line string) (*string, int, *string, error) {
	// Variable Format: |VariableName?x|
	// ? is the opcode, [+-*/\><&=]
	if line[0] != '|' && line[len(line)-1] != '|' {
		return nil, 0, nil, fmt.Errorf("Variable assignment needs delimiters %s", line)
	}
	line = line[1 : len(line)-1]
	idx := strings.IndexAny(line, "+-*/\\><&=")
	if idx == -1 {
		return nil, 0, nil, fmt.Errorf("No OpCode in %s", line)
	}
	op := string(line[idx])
	name := line[:idx]
	valstr := line[idx+1:]
	val, err := strconv.ParseInt(valstr, 0, 0)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("Variable value is not an integer %s in %s", valstr, line)
	}
	return &name, int(val), &op, nil
}

func Parse(path string) (*Table, error) {
	errorFmt := "Error: %s, line %d\n"
	content, err := tfs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	name := makeName(path, "")
	table := NewTable(name)
	table.Size = len(content)
	table.Path = path
	var group *Group
	var state int

	for lineno, line := range content {
		// Check for comment and strip it
		idx := strings.Index(line, "#")
		if idx > -1 {
			line = line[:idx]
		}
		// trim the line
		line = strings.TrimRight(line, "\r\n")
		line = strings.TrimLeft(line, " ")
		// if there is nothing to parse go to next line
		// blank line also closes any previous group parsing
		if len(line) == 0 {
			state = NO_GROUP
			if group != nil {
				table.AddGroup(group)
				group = nil
			}
			continue
		}
		// now we parse data lines
		if state == COLON_GROUP {
			if line[0] == '<' {
				// is it a prefix?
				group.Prefix = line[1:]
			} else if line[0] == '>' {
				// is it a suffix?
				group.Suffix = line[1:]
			} else if line[0] == '_' {
				// is it continuation
			} else {
				min, max, text, err := parseColonItem(line)
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				} else {
					group.AddItem(min, max, text)
				}
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
			} else {
				num, text, err := parseSemiItem(line)
				if err != nil {
					return nil, fmt.Errorf(errorFmt, err, lineno)
				} else {
					group.AddItem(1, num, text)
				}
			}
		} else if line[0] == ':' {
			state = COLON_GROUP
			// Save any previous group
			if group != nil {
				table.AddGroup(group)
				group = nil
			}
			group = &Group{}
			group.Name = line[1:]
		} else if line[0] == ';' {
			state = SEMI_GROUP
			// Save any previous group
			if group != nil {
				table.AddGroup(group)
				group = nil
			}
			group = &Group{}
			group.Name = line[1:]
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

			name, val, op, err := parseVariableAssignment(line)
			if err != nil {
				err = fmt.Errorf(errorFmt, err, lineno)
				return nil, err
			}
			i64, err := strconv.ParseInt(table.Variables[*name], 10, 0)
			oval := int(i64)
			switch *op {
			case "+":
				val = oval + val
			}

			table.AddVariable(*name, fmt.Sprintf("%d", oval))
		}
	}
	if group != nil {
		table.AddGroup(group)
		group = nil
	}
	return table, nil
}
