package tables

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *Table) Roll(gn string) string {
	//return t.OldRoll(gn)
	return t.TryRoll(gn)
}

func (t *Table) TryRoll(gn string) string {

	var gen string

	// There is a table call syntax where the 'random' roll
	// is passed in the reference as follows;
	// [Group=%Number%]
	// This is same as lookin up Group table at row 'Number'
	words := strings.Split(gn, "=")
	pick := -1
	var err error
	if len(words) == 2 {
		gn = words[0]
		pick, err = strconv.Atoi(words[1])
		if err != nil {
			return fmt.Sprintf("%s does not select an integer", gn)
		}
	}
	g := t.Groups[gn]
	if g == nil {
		return gn
	}

	if pick == -1 {
		gen = g.Roll()
	} else {
		gen = g.Select(pick)
	}

	gen = t.Evaluate(gen)

	return gen
}

// Starting from the beginning of s, find the end bracket, allow for nesting
func findEndDelim(s string, begin string, end string) (subStr string, lastIndex int) {
	n := 0
	lst := 0
	ret := ""
	for lst = 0; (n > 0 && s[lst] == ']') || lst < len(s); lst++ {
		c := s[lst : lst+1]
		if c == begin { // found a nested reference
			n += 1
			ret += c
		} else if c == end {
			n -= 1
			if n > 0 { // keep sub references
				ret += c
			} else {
				return ret, lst
			}
		} else {
			ret += c
		}
	}
	return ret, lst
}

/*
2,hexagonal|TempNumber={Ceil~{Calc~(%ValueFactor%*0.09)}}||ValueFactor=%TempNumber%|
1,crescent-shaped|TempNumber={Ceil~{Calc~(%ValueFactor%*0.05)}}||ValueFactor=%TempNumber%|
*/
func (t *Table) Evaluate(orig string) string {

	gen := ""

	for j := 0; j < len(orig); j++ {
		switch orig[j] {
		case '[':
			sub, last := findEndDelim(orig[j+1:], "[", "]")
			j += last
			sub = t.Evaluate(sub)
			gen += t.Roll(sub)
		case '{':
			sub, last := findEndDelim(orig[j+1:], "{", "}")
			j += last
			sub = t.Evaluate(sub)
			words := strings.Split(sub, "~")
			res, err := BuiltinCall(t, words[0], words[1])
			if err != nil {
				return "\n--ERROR Calling Builtin-- " + fmt.Sprintf("%s(%s): %s\n", words[0], words[1], err)
			}
			gen += res
		case '%':
			j += 1
			idx := strings.Index(orig[j:], "%")
			varName := orig[j : j+idx]
			v, ok := t.GetVariable(varName)
			if ok {
				gen += t.Evaluate(v)
			} else {
				return "\n--ERROR Accessing Variable-- %" + orig[j:idx] + "% does not exist"
			}
			j += idx
		default:
			gen += orig[j : j+1]
		}
	}
	return gen
}

func (t *Table) OldRoll(gn string) string {
	g := t.Groups[gn]
	if g == nil {
		return gn
	}

	s := g.Roll()

	var result string
	// Replace all table references with text
	for j := 0; j < len(s); j++ {
		c := s[j]
		if c == '[' { // reference to another table
			refgroup := ""
			j++ // skip past open bracket
			for s[j] != ']' {
				refgroup = refgroup + string(s[j])
				j++
			}
			//j++ // skip past close bracket
			foreigncall := strings.Split(refgroup, ".")
			switch len(foreigncall) {
			case 1:
				result = result + t.Roll(refgroup)
			case 2:
				nt, err := Parse(foreigncall[0])
				if err != nil {
					fmt.Printf("Syntax Error: Error parsing Table, %s in %s\n", refgroup, t.Path)
				} else {
					result = result + nt.Roll(foreigncall[1])
				}
			default:
				fmt.Printf("Syntax Error: Table call, %s in %s\n", refgroup, t.Path)
			}
		} else {
			result = result + string(c)
		}
	}
	s = result
	result = ""
	// Replace all variable references with text
	for j := 0; j < len(s); j++ {
		c := s[j]
		if c == '%' { // reference to another table
			variable := ""
			j++ // skip past open bracket
			for s[j] != '%' {
				variable = variable + string(s[j])
				j++
			}
			//j++ // skip past close bracket
			val, ok := t.GetVariable(variable)
			if ok {
				result = result + val
			} else {
				// variable does not exist, put the name back in place
				result = result + "%" + variable + "%"
			}
		} else {
			result = result + string(c)
		}
	}
	s = result
	result = ""
	for j := 0; j < len(s); j++ {
		c := s[j]
		if c == '{' {
			fcall := ""
			j++ // skip past open bracket
			for j < len(s) && s[j] != '}' {
				fcall = fcall + string(s[j])
				j++
			}
			//j++ // skip past close bracket
			words := strings.Split(fcall, "~")
			// TODO i := strings.Index(fcall, "~")
			if len(words) != 2 {
				fmt.Printf("Syntax Error: Builtin call has no ~, %s in %s\n",
					fcall,
					t.Path)
				continue
			}
			res, err := BuiltinCall(t, words[0], words[1])
			if err != nil {
				fmt.Printf("Error: %s(%s): %s\n", words[0], words[1], err)
				res = "-ERROR-"
			}
			result = result + res
		} else {
			result = result + string(c)
		}
	}
	return result
}
