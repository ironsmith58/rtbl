package tables

import (
	"bufio"
	_ "embed"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/Knetic/govaluate"
	"github.com/nboughton/go-roll"
)

//go:embed version.txt
var Version string

type BuiltInFunc func(*Table, string) (string, error)

type Builtin struct {
	Name  string
	BFunc BuiltInFunc
}

func BuiltinCall(t *Table, fname, args string) (string, error) {
	for j := 0; j < len(FunctionRegistry); j++ {
		if strings.EqualFold(fname, FunctionRegistry[j].Name) {
			defer func() { // recovers panic
				if e := recover(); e != nil {
					fmt.Println("Recovered from panic", e)
				}
			}()
			res, err := FunctionRegistry[j].BFunc(t, args)
			return res, err
		}
	}
	return "", fmt.Errorf("No builtin function named %s", fname)
}

var (
	lexMath = regexp.MustCompile(`[+\-*/]\s*\d+`)
)

func isNumber(s string) bool {
	dotFound := false

	// + or - are allowed as 1st char only
	// if see allow it, and check all other chars
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	for _, v := range s {
		if v == '.' {
			if dotFound {
				return false
			}
			dotFound = true
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func strtoi(s string) (int, int) {
	j := 0
	i := 0
	for ; j < len(s); j++ {
		c := s[j]
		if c <= '0' || c >= '9' {
			break
		}
		i = i*10 + int((byte(c) - byte('0')))
	}
	return i, j
}

func isVowel(c byte) bool {
	if strings.Contains("aeiouAEIOU", string(c)) {
		return true
	}
	return false
}

func stripInsignificantDigits(s string) string {
	//strip off insignificant trailing zeros
	for j := len(s) - 1; j > 0; j-- {
		if s[j] == '0' { // trailing zero, strip and continue
			s = s[:j]
		} else if s[j] == '.' { // trailing decimal, strip and stop
			s = s[:j]
			break
		} else {
			break // not a trailing zero or decimal, stop
		}
	}
	return s
}

func helperInputList(t *Table, s string) (string, error) {
	options := strings.Split(s, ",")
	def, err := strconv.Atoi(options[0])
	if err != nil {
		return "", fmt.Errorf("InpuList~Def,Prompt,Option,... %s is not a number", options[0])
	}
	fmt.Println(options[1])
	prefix := ""
	for num, p := range options[2:] {
		if num == def {
			prefix = "*"
		} else {
			prefix = " "
		}
		prefix += fmt.Sprintf("%d)", num)
		fmt.Println(prefix, p)
	}
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("An error occured while reading input. %s", err)
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	var choice int
	if len(input) == 0 {
		choice = def
	} else {
		choice, err = strconv.Atoi(input)
	}
	return options[choice+2], nil
}

//
// Array of BuiltIn Functions
//
var FunctionRegistry []Builtin = []Builtin{
	{
		Name: "Abs",
		BFunc: func(t *Table, s string) (string, error) {
			n, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if n < 0 {
				n = -n
			}
			return strconv.Itoa(n), nil
		},
	}, {
		Name: "AorAn",
		BFunc: func(t *Table, s string) (string, error) {
			first := 0
			if s[0:2] == "a " {
				first = 2
			} else if s[0:3] == "an " {
				first = 3
			}
			if len(s[first:]) == 0 {
				return s, nil
			}
			c := s[first]
			if isVowel(c) {
				return "an " + s[first:], nil
			}
			return "a " + s[first:], nil
		},
	},
	{
		Name: "Cap",
		BFunc: func(t *Table, s string) (string, error) {
			return strings.ToUpper(s), nil
		},
	},
	{
		Name: "CapEachWord",
		BFunc: func(t *Table, s string) (string, error) {
			return strings.Title(s), nil
		},
	},
	{
		Name: "Ceil",
		BFunc: func(t *Table, s string) (string, error) {
			var f float64
			_, err := fmt.Sscanf(s, "%f", &f)
			if err != nil {
				return "", err
			}
			f = math.Ceil(f)
			ret := strconv.FormatFloat(f, 'f', 3, 64)
			ret = stripInsignificantDigits(ret)
			return ret, nil
		},
	},
	{
		Name: "CharRet",
		BFunc: func(t *Table, s string) (string, error) {
			return "\n", nil
		},
	},
	{
		Name: "CR",
		BFunc: func(t *Table, s string) (string, error) {
			return "\n", nil
		},
	},
	{
		Name: "Char",
		BFunc: func(t *Table, s string) (string, error) {
			args := strings.Split(s, ",")
			j, err := strconv.Atoi(args[0])
			if err != nil {
				return "", err
			}
			if j > len(args[1]) {
				return s, fmt.Errorf("Char~%d is greater than the the length of '%s'", j, args[1])
			}
			return strings.ToUpper(s), nil
		},
	},
	{
		Name: "Color",
		BFunc: func(t *Table, s string) (string, error) {
			ic := strings.Index(s, ",")
			if ic == -1 {
				return "", fmt.Errorf("Color~Color,Text: No color supplied in %s", s)
			}
			color := s[:ic]
			text := s[ic+1:]
			return fmt.Sprintf("<font color=\"%s\">%s</font>", color, text), nil
		},
	},
	{
		Name: "Dice",
		BFunc: func(t *Table, s string) (string, error) {
			// roll.FromString does not support math operators
			// so i will removethem first and process them later
			maths := lexMath.FindAllString(s, -1)
			if len(maths) > 0 {
				j := strings.Index(s, maths[0]) // get index of first match
				s = s[:j]
			}
			// get random roll
			res, err := roll.FromString(s)
			if err != nil {
				return "", err
			}
			// sum all the rolls
			sum := res.Sum()
			var fact int
			// perform post RNG math
			if len(maths) > 0 {
				for _, m := range maths {
					op := m[0]
					// skip 1st char and whitespace
					j := 1
					for ; j < len(m); j++ {
						if m[j] != ' ' {
							break
						}
					}
					m = m[j:]
					// convert remaining chars as an int
					fact, err = strconv.Atoi(m)
					if err != nil {
						return "", err
					}
					// now it later, lets perform the operation
					switch op {
					case '+':
						sum = sum + fact
					case '-':
						sum = sum - fact
					case '*':
						sum = sum * fact
					case '/':
						sum = int(sum / fact)
					}
				}
			}
			return strconv.Itoa(sum), nil
		},
	},
	{
		Name: "Floor",
		BFunc: func(t *Table, s string) (string, error) {
			var f float64
			_, err := fmt.Sscanf(s, "%f", &f)
			if err != nil {
				return "", err
			}
			f = math.Floor(f)
			ret := strconv.FormatFloat(f, 'f', 3, 64)
			//strip off insignificant trailing zeros
			for j := len(ret) - 1; j > 0; j-- {
				if ret[j] == '0' { // trailing zero, strip and continue
					ret = ret[:j]
				} else if ret[j] == '.' { // trailing decimal, strip and stop
					ret = ret[:j]
					break
				} else {
					break // not a trailing zero or decimal, stop
				}
			}
			return ret, nil
		},
	},
	{
		Name: "If",
		BFunc: func(t *Table, s string) (string, error) {
			args := strings.Split(s, ",")
			args[0] = strings.Replace(args[0], "%", "", -1)
			expression, err := govaluate.NewEvaluableExpression(args[0])
			if err != nil {
				return "", fmt.Errorf("If~%s, %s", s, err)
			}
			parameters := make(map[string]interface{}, len(t.Variables))
			for n, v := range t.Variables {
				ival, end := strtoi(v)
				if end != len(v) {
					continue // this is not a number....
				}
				parameters[n] = ival
			}
			result, err := expression.Evaluate(parameters)
			args = strings.Split(args[1], "/")
			if result == true {
				return fmt.Sprintf("%v", args[0]), err
			}
			return fmt.Sprintf("%v", args[1]), err
		},
	},
	{
		Name:  "InputList",
		BFunc: helperInputList,
	},
	{
		Name: "IsNumber",
		BFunc: func(t *Table, s string) (string, error) {
			nm := isNumber(s)
			if nm {
				return "1", nil
			}
			return "0", nil
		},
	},
	{
		Name: "LCase",
		BFunc: func(t *Table, s string) (string, error) {
			return strings.ToLower(s), nil
		},
	},
	{
		Name: "Left",
		BFunc: func(t *Table, s string) (string, error) {
			args := strings.Split(s, ",")
			if len(args) != 2 {
				return "", fmt.Errorf("No offset in Left~%s", s)
			}
			j, _ := strconv.Atoi(args[0])
			if j > len(args[1]) {
				return args[1], nil
			}
			return args[1][:j], nil
		},
	},
	{
		Name: "Length",
		BFunc: func(t *Table, s string) (string, error) {
			l := len(s)
			return strconv.Itoa(l), nil
		},
	},
	{
		Name: "Mid",
		BFunc: func(t *Table, s string) (string, error) {
			lenret, comma1 := strtoi(s)
			start, comma2 := strtoi(s[comma1+1:])
			if s[comma2] != ',' {
				return "", fmt.Errorf("Mid~Len,Start,String: bad arguments %s", s)
			}
			if start > len(s) || comma1+1+comma2+1 > len(s) {
				return "", fmt.Errorf("Mid~Len,Start,String: Start is past end %s", s)
			}
			s = s[comma1+1+comma2+1:]
			return s[start : start+lenret], nil
		},
	},
	{
		Name: "OrderAsc",
		BFunc: func(t *Table, s string) (string, error) {
			//OrderAsc~"X",Text
			delim := string(s[1]) // delimiter must be 1 char, in quotes
			// s[0] and s[2] are the quotes
			if s[3] != ',' {
				return "", fmt.Errorf("OrderAsc~%s is missing a delimiter", s)
			}
			s = s[4:]
			words := strings.Split(s, delim)

			sort.Strings(words)
			return strings.Join(words, delim), nil
		},
	},
	{
		Name: "OrderDesc",
		BFunc: func(t *Table, s string) (string, error) {
			//OrderAsc~"X",Text
			delim := string(s[1]) // delimiter must be 1 char, in quotes
			if s[3] != ',' {
				return "", fmt.Errorf("OrderDesc~%s is missing a delimiter", s)
			}
			s = s[4:]
			words := strings.Split(s, delim)

			sort.Sort(sort.Reverse(sort.StringSlice(words)))
			return strings.Join(words, delim), nil
		},
	},
	{
		Name: "Ordinal",
		BFunc: func(t *Table, s string) (string, error) {
			// get right most digits, 'ones' position
			if s == "" {
				return "", nil
			}
			one := len(s) - 1
			suff := "th"
			// is this a teens number?
			if len(s) >= 2 && s[one-1] == '1' {
				return s + "th", nil
			}
			switch s[one] {
			case '1':
				suff = "st"
			case '2':
				suff = "nd"
			case '3':
				suff = "rd"
			}
			return s + suff, nil
		},
	},
	{
		Name: "Plural",
		BFunc: func(t *Table, s string) (string, error) {

			if s == "" {
				return "", nil
			}
			if strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "o") || strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "x") {
				return s + "es", nil
			}
			if strings.HasSuffix(s, "fe") {
				return s[:len(s)-2] + "ves", nil
			}
			if strings.HasSuffix(s, "f") {
				return s[:len(s)-1] + "ves", nil
			}
			if strings.HasSuffix(s, "y") && !isVowel(s[len(s)-2]) {
				return s[:len(s)-1] + "ies", nil
			}
			return s + "s", nil
		},
	},
	{
		Name: "Pluralif",
		BFunc: func(t *Table, s string) (string, error) {
			//{PluralIf~X,Text}
			//Description
			//This function will return "Text" in its plural form (see Plural for criteria)
			//if "X" does not equal 1.
			if s == "" {
				return "", nil
			}
			idx := strings.Index(s, ",")
			if idx == -1 {
				return "", fmt.Errorf("Pluralif~%s does not have a number", s)
			}
			n, err := strconv.ParseFloat(s[:idx], 32)
			if err != nil {
				return "", fmt.Errorf("Pluralif~%s 1st argument is not a number (%s)", s, s[:idx])
			}
			s = s[idx+1:]
			if n == 1 {
				return s, nil
			}
			if s == "" {
				return "", nil
			}
			if strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "o") || strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "x") {
				return s + "es", nil
			}
			if strings.HasSuffix(s, "fe") {
				return s[:len(s)-2] + "ves", nil
			}
			if strings.HasSuffix(s, "f") {
				return s[:len(s)-1] + "ves", nil
			}
			if strings.HasSuffix(s, "y") && !isVowel(s[len(s)-2]) {
				return s[:len(s)-1] + "ies", nil
			}
			return s + "s", nil
		},
	},
	{
		Name: "Replace",
		BFunc: func(t *Table, s string) (string, error) {
			//Replace~~SearchFor,~ReplaceWith,Text}
			//Description
			//Replaces each instance of "SearchFor" in "Text" with "ReplaceWith".
			args := strings.Split(s, ",")
			if len(args) != 3 {
				return "", fmt.Errorf("Replace~%s not enough arguments", s)
			}
			if len(args[0]) == 0 {
				return "", fmt.Errorf("Replace~%s missing SearchFor word", s)
			}
			// ReplaceWith may be blank when deleting
			if len(args[2]) == 0 {
				return "", nil
			}

			ret := strings.Replace(args[2], args[0], args[1], -1)
			return ret, nil
		},
	},
	{
		Name: "Reset",
		BFunc: func(t *Table, s string) (string, error) {
			// s is the GroupName to reset
			g, err := t.GetGroup(s)
			if err != nil {
				return "", fmt.Errorf("Reset~%s: nonexistent group", s)
			}
			g.Reset()
			return "", nil
		},
	},
	{
		Name: "Right",
		BFunc: func(t *Table, s string) (string, error) {
			args := strings.Split(s, ",")
			if len(args) != 2 {
				return "", fmt.Errorf("No offset in Right~%s", s)
			}
			if args[1] == "" {
				return "", nil
			}
			j, _ := strconv.Atoi(args[0])
			j = len(args[1]) - j
			if j < 0 {
				return args[1], nil
			}
			return args[1][j:], nil
		},
	},
	{
		Name: "Round",
		BFunc: func(t *Table, s string) (string, error) {
			// {Round~X,Value}
			args := strings.Split(s, ",")
			if len(args) != 2 {
				return "", fmt.Errorf("No decimal places in Round~%s", s)
			}
			if args[1] == "" {
				return "", nil
			}
			//roundFloat(val float64, precision uint) float64
			precision, err := strconv.ParseFloat(args[0], 32)
			if err != nil {
				return "", fmt.Errorf("Round~%s precision is not a number", s)
			}
			ratio := math.Pow(10, float64(precision))
			val, err := strconv.ParseFloat(args[1], 32)
			if err != nil {
				return "", fmt.Errorf("Round~%s value is not a number", s)
			}
			ret := fmt.Sprintf("%f", math.Round(val*ratio)/ratio)
			ret = stripInsignificantDigits(ret)
			return ret, nil
		},
	},
	{
		Name: "Space",
		BFunc: func(t *Table, s string) (string, error) {
			i, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Space~%s is not a number", s)
			}
			pad := ""
			for j := 0; j < i; j++ {
				pad += " "
			}
			return pad, nil
		},
	},
	{
		Name: "Spc",
		BFunc: func(t *Table, s string) (string, error) {
			i, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Space~%s is not a number", s)
			}
			pad := ""
			for j := 0; j < i; j++ {
				pad += " "
			}
			return pad, nil
		},
	},
	{
		Name: "Sqrt",
		BFunc: func(t *Table, s string) (string, error) {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return "", fmt.Errorf("Sqrt~%s is not a number", s)
			}
			f = math.Sqrt(f)
			ret := fmt.Sprintf("%f", f)
			ret = stripInsignificantDigits(ret)
			return ret, nil
		},
	},
	{
		Name: "Status",
		BFunc: func(t *Table, s string) (string, error) {
			return s, nil
		},
	},
	{
		Name: "Title",
		BFunc: func(t *Table, s string) (string, error) {
			res := strings.ToLower(s)
			res = strings.Title(res)
			// count the words in the title
			// if the title is a short title, 4 words or less
			// return it as is, without lower case some
			// prepositions, etc
			nwords := len(strings.Split(res, " "))
			if nwords <= 4 {
				return res, nil
			}
			lowerwords := []string{"After", "As Of", "From", "Ago", "At", "Before",
				"By", "During", "For", "From", "Until", "Till", "To", "In", "On", "Past", "Since",
				"To", "Of",
				"Above", "Across", "Against", "Behind", "Below", "Between", "Next To", "Beside",
				"Into", "In Front Of", "Near", "Onto", "Over", "Under",
				"About", "With", "If", "Else", "And", "Or", "The", " A ", " An "}
			for _, w := range lowerwords {
				res = strings.Replace(res, w, strings.ToLower(w), -1)
			}
			// Captitalize first char
			r := []rune(res)
			r[0] = unicode.ToUpper(r[0])
			res = string(r)
			return res, nil
		},
	},
	{
		Name: "Trim",
		BFunc: func(t *Table, s string) (string, error) {
			return strings.TrimSpace(s), nil
		},
	},
	{
		Name: "Trunc",
		BFunc: func(t *Table, s string) (string, error) {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return "", fmt.Errorf("Trunc~%s is not a number", s)
			}
			return fmt.Sprintf("%d", int(f)), nil
		},
	},
	{
		Name: "UCase",
		BFunc: func(t *Table, s string) (string, error) {
			return strings.ToUpper(s), nil
		},
	},
	{
		Name: "Version",
		BFunc: func(t *Table, s string) (string, error) {
			return Version, nil
		},
	},
	{
		Name: "VowelStart",
		BFunc: func(t *Table, s string) (string, error) {
			if s == "" {
				return "0", nil
			}
			//VowelStart~Text
			//~VowelStart~%var1%
			if s[0] == '%' {
				v, exists := t.GetVariable(s[1 : len(s)-1])
				if !exists {
					return "", fmt.Errorf("VowelStart~%%%s%% is not an existing variable", s)
				}
				s = v
			}
			// this code does not support variables, only text
			if isVowel(s[0]) {
				return "1", nil
			}
			return "0", nil
		},
	},
}
