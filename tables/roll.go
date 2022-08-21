package tables

func (t *Table) Roll(gn string) string {
	g := t.Groups[gn]
	s := g.Roll()
	return s
}
