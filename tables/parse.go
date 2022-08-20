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
    text := ""

    line = strings.TrimLeft(line, " ")
    words := strings.FieldsFunc(line, func(c rune) bool {
        if c == ',' || c == ' ' || c == '\t' {
            return true
        }
        return false
    })
    if len(words) != 2 {
        return 0, 0, "",
            fmt.Errorf("Colon Group: missing delimiter between range and text '%s'", line)
    }
    rndmm := strings.Split(words[0], "-")
    min, err = strconv.ParseInt(rndmm[0], 10, 0)
    if err != nil {
        return 0, 0, "", err
    }
    // Random range is optional
    // there might not be a 'max'
    if len(rndmm) == 2 {
        max, err = strconv.ParseInt(rndmm[1], 10, 0)
        if err != nil {
            return 0, 0, "", err
        }
    }
    return int(min), int(max), text, nil
}

func parseSemiItem(line string) (int, string, error) {
    var err error
    var num int64
    text := ""

    line = strings.TrimLeft(line, " ")
    words := strings.FieldsFunc(line, func(c rune) bool {
        if c == ',' || c == ' ' || c == '\t' {
            return true
        }
        return false
    })
    if len(words) == 1 {
        return 0, "",
            fmt.Errorf("Semi Group: missing delimiter between range and text '%s'", line)
    }
    num, err = strconv.ParseInt(words[0], 10, 0)
    if err != nil {
        return 0, "", err
    }
    return int(num), text, nil
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
        line = strings.TrimRight(line, " \t\r\n")
        // if there is nothing to parse go to next line
        // blank line also closes any previous group parsing
        if len(line) == 0 {
            state = NO_GROUP
            if group != nil {
                table.AddGroup(group)
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
                    group.AddItem(num, 0, text)
                }
            }
        } else if line[0] == ':' {
            state = COLON_GROUP
            // Save any previous group
            if group != nil {
                table.AddGroup(group)
            }
            group = &Group{}
            group.Name = line[1:]
        } else if line[0] == ';' {
            state = SEMI_GROUP
            // Save any previous group
            if group != nil {
                table.AddGroup(group)
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
    }
    return table, nil
}
