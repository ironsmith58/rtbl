package datasets

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

/*
 * Dataset Registry maps DS name to actual data
 */
type dsRegistry map[string]*dataset

var registry dsRegistry = make(map[string]*dataset)

func findDS(name string) (*dataset, error) {
	ds, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("%s is not an current Dataset", name)
	}
	return ds, nil
}

func addDS(name string, ds *dataset) error {
	name = strings.ToLower(name)
	ds, _ = findDS(name)
	if ds != nil {
		return fmt.Errorf("%s is an existing DS", ds.name)
	}
	registry[name] = ds
	return nil
}

// A single row of cells in a dataset
type row []string

// a dataset, with
//   header names
//   a row of default value
//   and an array of rows/values
type dataset struct {
	name     string
	headers  row
	defaults row
	rows     []row
}

func newDS(dsname string) *dataset {
	return &dataset{
		name: dsname,
		rows: make([]row, 0, 10),
	}
}

func (d *dataset) NewRow() row {
	return append([]string{}, d.defaults...)
}

func (d *dataset) AddRow(nrow row) {
	d.rows = append(d.rows, nrow)
}

func (d *dataset) IsColumn(c string) int {
	for j := 0; j < len(d.headers); j++ {
		if c == d.headers[j] {
			return j
		}
	}
	return -1
}

func DSAdd(s string) (string, error) {
	//DSAdd~VarName,Field1,Value1,Field2,Value2,...
	fields := strings.Split(s, ",")
	ds, err := findDS(fields[0])
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name", fields[0])
	}
	newrow := ds.NewRow() // get new row to defaults
	// go thru each field, check to see if it exists
	// set values that are provided
	for j := 1; j < len(fields); j = j + 2 {
		fld := fields[j]
		val := fields[j+1]
		// find if this is a column, by name
		// and its location as an index
		idx := ds.IsColumn(fld)
		if idx == -1 {
			return "", fmt.Errorf("%s is not a column in dataset %s", fld, fields[0])
		}
		newrow[idx] = val
	}
	// values set, so lets add the new row to the dataset
	ds.AddRow(newrow)
	// return index of added row
	return strconv.Itoa(len(ds.rows) - 1), nil
}

func DSAddNR(s string) (string, error) {
	_, err := DSAdd(s) // throw away index of row
	return "", err
}

func DSCalc(s string) (string, error) {
	//DSCalc~VarName,Operation,Field
	fields := strings.Split(s, ",")
	ds, err := findDS(fields[0])
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name", fields[0])
	}

	idx := ds.IsColumn(fields[2])
	if idx == -1 {
		return "", fmt.Errorf("%s is not a column in dataset %s", fields[2], fields[0])
	}
	// Sum all the columns values
	acc := 0.0
	for j := 0; j < len(ds.rows); j++ {
		rn, err := strconv.ParseFloat(ds.rows[j][idx], 32)
		if err != nil {
			return "", fmt.Errorf("Column %s has non-numeric value %s", fields[2], ds.rows[j][idx])
		}
		acc = acc + rn
	}
	// if user called for an Average ....
	if strings.ToLower(fields[1]) == "avg" {
		acc = acc / float64(len(ds.rows))
	}
	return strconv.FormatFloat(acc, 'f', 3, 64), nil
}

func DSCount(s string) (string, error) {
	//DSCount~VarName
	ds, err := findDS(s)
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name", s)
	}
	return strconv.Itoa(len(ds.rows)), nil
}

func DSCreate(s string) (string, error) {
	//DSCreate~VarName,Field1,Default1,Field2,Default2,...Fieldx,Defaultx
	fields := strings.Split(s, ",")
	dsname := fields[0]
	defrow := make(map[string]string)
	ds := newDS(dsname)
	for j := 0; j < len(fields); j = j + 2 {
		ds.headers = append(ds.headers, fields[j])
		defrow[fields[j]] = fields[j+1]
	}
	addDS(dsname, ds)
	return "", nil
}
func DSFind(s string) (string, error) {
	//DSFind~VarName,Index,Expr1,Expr2,...
	/*
		Starting at the item with index "Index", searches through each item until it finds
		one where the given "Exprx" are true and returns the index of that item (if not match
		is found, "-1" is returned). "Exprx" follows the format "fieldname XXX value".
		"XXX" is a comparison operator, similar to what you use for "If/IIf". The
		following operators are recognized: =, !=, <=, >=, >, <, ~, and !~.
		All but the last two work both with numbers and text. The last two work only
		with text. "~" means "like" and "!~" means "not like". If you want to use
		wildcards with your text search, use "~" and "!~".
	*/
	return "", nil

}
func DSGet(s string) (string, error) {
	//DSGet~VarName,Index,Field
	fields := strings.Split(s, ",")
	ds, err := findDS(s)
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name", s)
	}
	irow, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", fmt.Errorf("DSGet~%s is not a valid index", s)
	}
	icol := ds.IsColumn(fields[2])
	if icol == -1 {
		return "", fmt.Errorf("DSGet~%s is not a valid field", s)
	}
	return ds.rows[irow][icol], nil
}

func DSRandomize(s string) (string, error) {
	ds, err := findDS(s)
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name")
	}

	rand.Shuffle(len(ds.rows), func(i, j int) {
		ds.rows[i], ds.rows[j] = ds.rows[j], ds.rows[i]
	})
	return "", nil
}
func DSRead(s string) (string, error) {
	return "", nil

}
func DSRemove(s string) (string, error) {
	//DSRemove~VarName,Index
	return "", nil

}
func DSRoll(s string) (string, error) {
	//DSRoll~VarName,Field@Mod
	return "", nil

}
func DSSet(s string) (string, error) {
	//DSSet~VarName,Index,Field1,Value1,Field2,Value2,...
	return "", nil

}
func DSSort(s string) (string, error) {
	//DSSort~~VarName,Field1,Direction1,Field2,Direction2,...
	return "", nil

}
func DSWrite(s string) (string, error) {
	//DSWrite~VarName,Filename
	args := strings.Split(s, ",")
	dsname := args[0]
	ds, err := findDS(dsname)
	if err != nil {
		return "", fmt.Errorf("%s is not a dataset name", dsname)
	}
	mkdatadir()
	fname := "Data/" + args[1]
	dsfile, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	// remember to close the file
	defer dsfile.Close()

	dsfile.WriteString("# DataSet written by RTBL in RDB Format")
	dsfile.WriteString("# RDB Format is a Tab Seperate Values with a 2 line header")
	// Write Names of columns; header line 1
	dsfile.WriteString(strings.Join(ds.headers, "\t"))
	// Write width of columns; header line 2
	//  but first calculate width
	max := make([]int, len(ds.headers))
	for i := 0; i < len(ds.rows); i++ {
		row := &ds.rows[i]
		for j := 0; j < len_row; j++ {
			if len(row[j]) > max[j] {
				max[j] = len(row[j])
			}
		}
	}
	maxs := make([]string, len(ds.headers))
	for j := 0; j < len(max); j++ {
		maxs[j] = strconv.Itoa(max[j])
	}
	dsfile.WriteString(strings.Join(maxs, "\t"))

	// Write the default values as 1st line
	dsfile.WriteString(strings.Join(ds.defaults, "\t"))
	// Write all the rows
	for j := 0; j < len(ds.rows); j++ {
		dsfile.WriteString(strings.Join(ds.rows[j], "\t"))
	}
}
