package tables

/* Groups are named sub-tables within a Table file
 * e.g. ;Name
 *      1,Fred
 *      1,Paul
 *      ....
 */

import (
	"fmt"
	"strconv"

	"github.com/nboughton/go-roll"
)

/*
Group struct holds a single table within a file
it has a unique name allowing it to be referenced
with a link style name, e.g. [Start]
*/
type Group struct {
	Name     string              // unique name within a table
	useOnce  bool                // after an entry is used it is removed
	probType rune                // Relative or Absolute Probability
	Prefix   string              // string placed before all random entries when returned
	Suffix   string              // string placed after all random entries when returned
	maxRoll  int                 // all rolls are essentially 1D{maxRoll}
	table    roll.Table          // table of entries and their percentage chance of appearing
	seen     map[string]struct{} // for useOnce groups, this holds previously seen values
}

const ABS_GROUP = ':' // flag for Absolute Percentage Chance group
const REL_GROUP = ';' // flag for Relative Percentage Chance group

func NewGroup(name string) *Group {
	probType := ABS_GROUP
	if name[0] == ':' { // Absolute range probablity group
		probType = ABS_GROUP
		name = name[1:]
	} else if name[0] == ';' { // Relative probabality group
		probType = REL_GROUP
		name = name[1:]
	}
	once := false
	for {
		if name[0] == '!' { // only return group entries once
			once = true
			name = name[1:]
		} else if name[0] == '~' { // reroll option is available on tableSmith screen
			name = name[1:] // reroll is ignored in CLI rtbl
		} else {
			break
		}
	}

	return &Group{
		Name:     name,
		useOnce:  once,
		probType: probType,
		table: roll.Table{
			Name: name,
			ID:   name,
		},
		seen: make(map[string]struct{}), // make a set using map of empty struct
	}
}

// the nboughton package requires an array of all 'face' values
// since we only use numeric dice; 1d4, 1d6, 1d8 ... 1d77, 1d103
// we must make each Die face bespoke for each table
func makeFaces(n int) roll.Faces {
	var f roll.Faces

	for j := 1; j <= n; j++ {
		f = append(f, roll.Face{N: j, Value: strconv.Itoa(j)})
	}
	return f
}

func (g *Group) Close() {
	last := len(g.table.Items)
	if last == 0 {
		return
	}
	match := g.table.Items[last-1].Match
	g.maxRoll = 0
	for j := 0; j < len(match); j++ {
		if match[j] > g.maxRoll {
			g.maxRoll = match[j]
		}
	}
	g.table.Dice = roll.Dice{N: 1, Die: roll.NewDie(makeFaces(g.maxRoll))}
}

func (g *Group) Len() int { return len(g.table.Items) }
func (g *Group) Min() int { return 1 }
func (g *Group) Max() int { return g.maxRoll }

// add a single item to this group with its matching percentage
func (g *Group) AddItem(start, end int, l string) {

	if g.probType == REL_GROUP {
		// only start has meaning for a relative group
		// these groups have a single int that is relative
		// to all the entries
		last := len(g.table.Items)
		if last == 0 {
			end = start
			start = 1
		} else {
			match := g.table.Items[last-1].Match
			largest := 0
			for j := 0; j < len(match); j++ {
				if match[j] > largest {
					largest = match[j]
				}
			}
			end = largest + start
			start = largest + 1
		}
	}
	g.table.Items = append(g.table.Items,
		roll.TableItem{
			Match: roll.MatchRange(start, end),
			Text:  l,
		})

}

// append string argument to the last Item
// used during parsing
// used to support underscore(_) continuation
func (g *Group) AppendLastItem(l string) error {
	if len(g.table.Items) == 0 {
		return fmt.Errorf("Can not append, no items in Group %s", g.Name)
	}
	last := len(g.table.Items) - 1

	g.table.Items[last].Text = g.table.Items[last].Text + l

	return nil
}

// defined as an empty struct so it takes no space in the map
// which is how we implement sets in golang
var dummy struct{}

// randomly select an entry from the group and apply prefix and suffix
// to returned value
// if this is a useOnce group,loop until a unique value can be
// returned
// this is implementaiton of UseOnce groups, :!Gear
func (g *Group) Roll() string {
	var s string
	// let not loop infinetely
	if g.useOnce {
		if len(g.seen) == len(g.table.Items) {
			return ""
		}
	}
	//repeatedly select a value until done
	for {
		s = g.table.Roll()
		_, alreadyUsed := g.seen[s]
		if !g.useOnce || !alreadyUsed {
			break
		}
	}
	// if this is a useOnce group save the returned value
	// so we can check to see if we used it next time
	if g.useOnce {
		g.seen[s] = dummy
	}
	return g.Prefix + s + g.Suffix
}

// select the Nth item from the table
func (g *Group) Select(n int) string {
	var s string

	//repeatedly select a value until done
	for j := range g.table.Items {
		if g.table.Items[j].Match.Contains(n) {
			s = g.table.Items[j].Text
			break
		}
	}
	return g.Prefix + s + g.Suffix
}

// Reset the state of the Group
// - delete the already used entries, so they can be re-used
func (g *Group) Reset() {
	for k := range g.seen {
		delete(g.seen, k)
	}
}
