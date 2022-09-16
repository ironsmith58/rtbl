package tables

import (
	"strconv"
	"testing"
)

func TestNewTable(t *testing.T) {
	tbl := NewTable("test-table")
	if len(tbl.Groups) != 0 {
		t.Log("Length of initial Groups is not 0")
		t.Fail()
	}
	if len(tbl.Variables) != 0 {
		t.Log("Length of Variables is not 0")
		t.Fail()
	}
}

func TestVariables(t *testing.T) {
	tbl := NewTable("test-table")

	if len(tbl.Variables) != 0 {
		t.Log("Length of Variables is not 0")
		t.Fail()
	}
	tbl.AddVariable("Var1", "value1")
	if len(tbl.Variables) != 1 {
		t.Log("Variable was not added")
		t.Fail()
	}
	v, exists := tbl.GetVariable("Var1")
	if !exists {
		t.Log("failed to retrieve variable Var1")
		t.Fail()
	}
	if v != "value1" {
		t.Log("Retrieved variable does not have correct value")
		t.Fail()
	}
}

func TestMakeFaces(t *testing.T) {
	rf := makeFaces(13)
	if len(rf) != 13 {
		t.Log(rf)
		t.Fail()
	}
	for j := 0; j < len(rf); j++ {
		if strconv.Itoa(rf[j].N) != rf[j].Value && rf[j].N != j+1 {
			t.Log("Faces is not in sequence")
			t.Log(rf)
			t.Fail()
		}
	}
}

func TestMakeGroup(t *testing.T) {
	g := NewGroup("group1")
	if g == nil {
		t.Log("Group failed, pointer nil")
		t.Fail()
	}

	if g.Name != "group1" {
		t.Log("Group name is incorrect")
		t.Fail()
	}

	g.AddItem(7, 0, "Item")
	if len(g.table.Items) != 1 {
		t.Logf("Expected 1 items, %d added instead", len(g.table.Items))
		t.Fail()
	}

	g.AppendLastItem("more stuff on item")
	if len(g.table.Items) != 1 {
		t.Logf("Expected 1 items, %d added instead", len(g.table.Items))
		t.Fail()
	}

}
