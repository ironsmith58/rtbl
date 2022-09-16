package tables

/*
 * Test the Builtin Functions defined in new.go
 *
 * Functions are in the array 'FunctionRegistry'
 */
import (
	"strconv"
	"testing"
)

func TestDice(t *testing.T) {
	tests := []struct {
		input    string
		min, max int
	}{
		{input: "3d6", min: 3, max: 18},
		{input: "2d12+4", min: 6, max: 28},
		{input: "4d6Dl1+10", min: 13, max: 28},
		{input: "4d10Kh3Dl1", min: 2, max: 20},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Dice", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			rd, err := strconv.Atoi(res)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if rd < tt.min || rd > tt.max {
				t.Logf("Case %d: %s %d not in range %d-%d",
					tcase, tt.input, rd, tt.min, tt.max)
				t.Fail()
			}
		})
	}
}

func TestAorAn(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "apple", expected: "an apple"},
		{input: "Orc", expected: "an Orc"},
		{input: "finkle", expected: "a finkle"},
		{input: "laser pistol", expected: "a laser pistol"},
		{input: "an laser pistol", expected: "a laser pistol"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "AorAn", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestCapEachWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "sword of ancient might", expected: "Sword Of Ancient Might"},
		{input: "ORC", expected: "ORC"},
		{input: "laser Pistol", expected: "Laser Pistol"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "capeachword", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "1", expected: "1"},
		{input: "1200.4", expected: "1201"},
		{input: "1200.8", expected: "1201"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "ceil", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "1", expected: "1"},
		{input: "1200.4", expected: "1200"},
		{input: "1200.8", expected: "1200"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "FLOOR", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   error
	}{
		{input: "3.14", expected: "1"},
		{input: "0.12", expected: "1"},
		{input: "+782", expected: "1"},
		{input: "-0.0", expected: "1"},
		{input: "tree", expected: "0"},
		{input: "42 skidoo", expected: "0"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "isnumber", tt.input)
			if tt.experr != err {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestLCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "Mighty Giant", expected: "mighty giant"},
		{input: "1200.4", expected: "1200.4"},
		{input: "space faring shuttle", expected: "space faring shuttle"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "lCase", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestLeft(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "6,Mighty Giant", expected: "Mighty"},
		{input: "42,1200.4", expected: "1200.4"},
		{input: "0,space faring shuttle", expected: ""},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "LEFT", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestLength(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "A Cargo hold filled with gems", expected: "29"},
		{input: "", expected: "0"},
		{input: "gem", expected: "3"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "length", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestMid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "5,2,A Cargo hold filled with gems", expected: "Cargo"},
		{input: "4,6", expected: "", experr: true},
		{input: "3,0,gem", expected: "gem", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Mid", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if !tt.experr && res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestOrderAsc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "\"|\",sword|dagger|gem|rope|apple", expected: "apple|dagger|gem|rope|sword"},
		{input: "\"|\",sword|dagger||gem|rope|apple", expected: "|apple|dagger|gem|rope|sword"},
		{input: "\"|\",|sword|dagger|gem|rope|apple|", expected: "||apple|dagger|gem|rope|sword"},

		{input: "", expected: "", experr: false},
		{input: ",,", expected: "", experr: false},
		{input: ",this is missing a delmiter", expected: "", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "OrderAsc", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d:%s failed: wanted err have nil", tcase, tt.input)
				t.Fail()
			} else if !tt.experr && err != nil {
				t.Logf("Case %d:%s failed: unexpected error %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestOrderDesc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "\"|\",sword|dagger|gem|rope|apple", expected: "sword|rope|gem|dagger|apple"},
		{input: "\"|\",sword|dagger||gem|rope|apple", expected: "sword|rope|gem|dagger|apple|"},
		{input: "\"|\",|sword|dagger|gem|rope|apple|", expected: "sword|rope|gem|dagger|apple||"},

		{input: "", expected: "", experr: false},
		{input: ",,", expected: "", experr: false},
		{input: ",this is missing a delmiter", expected: "", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "OrderDesc", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d:%s failed: wanted err have nil", tcase, tt.input)
				t.Fail()
			} else if !tt.experr && err != nil {
				t.Logf("Case %d:%s failed: unexpected error %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestOrdinal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},

		{input: "1", expected: "1st"},
		{input: "2", expected: "2nd"},
		{input: "3", expected: "3rd"},
		{input: "4", expected: "4th"},

		{input: "11", expected: "11th"},
		{input: "12", expected: "12th"},
		{input: "13", expected: "13th"},
		{input: "14", expected: "14th"},

		{input: "21", expected: "21st"},
		{input: "32", expected: "32nd"},
		{input: "43", expected: "43rd"},
		{input: "54", expected: "54th"},

		{input: "111", expected: "111th"},
		{input: "112", expected: "112th"},
		{input: "113", expected: "113th"},
		{input: "114", expected: "114th"},

		{input: "121", expected: "121st"},
		{input: "132", expected: "132nd"},
		{input: "143", expected: "143rd"},
		{input: "154", expected: "154th"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "ordinal", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestPlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "money", expected: "moneys"}, // strictly, plural of money is money
		{input: "lady", expected: "ladies"},
		{input: "dwarf", expected: "dwarves"},
		{input: "egg", expected: "eggs"},
		{input: "", expected: ""},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "pluRAL", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestPluralIf(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "2,money", expected: "moneys"}, // strictly, plural of money is money
		{input: "11,lady", expected: "ladies"},
		{input: "14,dwarf", expected: "dwarves"},
		{input: "0,egg", expected: "eggs"},
		{input: "42,", expected: ""},

		{input: "1,money", expected: "money"},
		{input: "1,lady", expected: "lady"},
		{input: "1,dwarf", expected: "dwarf"},
		{input: "1,egg", expected: "egg"},
		{input: "1,", expected: ""},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "pluralif", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestRight(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "5,Mighty Giant", expected: "Giant"},
		{input: "42,1200.4", expected: "1200.4"},
		{input: "0,space faring shuttle", expected: ""},
		{input: "4,", expected: ""},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "right", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "1,5.32", expected: "5.3"},
		{input: "1,5.36", expected: "5.4"},
		{input: "6,12.1234567890", expected: "12.123457"},
		{input: "", expected: "", experr: true},
		{input: "n,3.1456", expected: "", experr: true},
		{input: "4,Apple", expected: "", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Round", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d:%s failed: wanted err have nil", tcase, tt.input)
				t.Fail()
			} else if !tt.experr && err != nil {
				t.Logf("Case %d:%s failed: unexpected error %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "@@,Sword,the Giant is armed with a @@", expected: "the Giant is armed with a Sword"},
		{input: "@,HI,--@@--", expected: "--HIHI--"},
		{input: "////,lots,//// of goblins", expected: "lots of goblins"},
		{input: "giant,,giant goblins", expected: " goblins"},
		{input: "", expected: "", experr: true},
		{input: ",,", expected: "", experr: true},
		{input: ",Bye,Forty Four Skidoo", expected: "", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Replace", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d:%s failed: wanted err have nil", tcase, tt.input)
				t.Fail()
			} else if !tt.experr && err != nil {
				t.Logf("Case %d:%s failed: unexpected error %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "5", expected: "     "},
		{input: "0", expected: ""},
		{input: "1", expected: " "},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Space", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
			// Test Alias
			res, err = BuiltinCall(nil, "Spc", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "Author: Eric F. Wolcott", expected: "Author: Eric F. Wolcott"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Status", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "4", expected: "2"},
		{input: "71", expected: "8.42615"},
		{input: "45", expected: "6.708204"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "sqrt", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   error
	}{
		{input: "a red book", expected: "A Red Book"},
		{input: "when below stairs", expected: "When Below Stairs"},
		{input: "over, under, around time", expected: "Over, Under, Around Time"},
		{input: "a brief history of constantinople", expected: "A Brief History of Constantinople"},
		{input: "a roLLicking trip through the lower reaches of rome", expected: "A Rollicking Trip Through the Lower Reaches of Rome"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Title", tt.input)
			if tt.experr != err {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: " \t \n", expected: ""},
		{input: "\tgiant Goblin  ", expected: "giant Goblin"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Trim", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestTrunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		experr   bool
	}{
		{input: "3.14", expected: "3"},
		{input: "0.12", expected: "0"},
		{input: "+782", expected: "782"},
		{input: "-0.0", expected: "0"},
		{input: "tree", expected: "", experr: true},
		{input: "42 skidoo", expected: "", experr: true},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "Trunc", tt.input)
			if tt.experr && err == nil {
				t.Logf("Case %d: wanted an error have: nil", tcase)
				t.Fail()
			} else if !tt.experr && err != nil {
				t.Logf("Case %d: did not expect error have: %s", tcase, err)
				t.Fail()
			} else if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}

func TestVowelStart(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: "0"},
		{input: "goblin", expected: "0"},
		{input: "apples", expected: "1"},
		{input: "elves", expected: "1"},
		{input: "ichor", expected: "1"},
		{input: "orc", expected: "1"},
		{input: "umber hulk", expected: "1"},
	}

	for tcase, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := BuiltinCall(nil, "VowelStart", tt.input)
			if err != nil {
				t.Logf("Case %d:%s failed: %s", tcase, tt.input, err)
				t.Fail()
			}
			if res != tt.expected {
				t.Logf("Case %d: wanted %s, have %s", tcase, tt.expected, res)
				t.Fail()
			}
		})
	}
}
