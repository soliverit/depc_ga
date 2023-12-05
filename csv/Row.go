package csv

import (
	"strings"
)

type Row struct {
	cells []string
	Cells []string
}

func NewRow(strs []string) Row {
	var row Row
	row.cells = strs
	return row
}

func StringToRow(str string) Row {
	var row Row
	var cells []string = make([]string, 0)
	var curPos uint32 = 0
	var inString bool = false

	//Sanitise the input
	str = strings.TrimSpace(str)
	/*
		Parse cells

		TODO: Why'd you use uint32 for i. Must be something to do with the inter but can't remember. First
		TODO: script so don't sweat it.
	*/
	for i := uint32(0); i < uint32(len(str)); i++ {
		/*
			Deal with comma encapsulated
		*/
		if inString {
			for str[i] != '"' {
				i += 1
			}
			inString = false
		} else {
			if str[i] == ',' {
				cells = append(cells, str[curPos:i])
				curPos = 1 + i
			} else if i+1 == uint32(len(str)) {
				cells = append(cells, str[curPos:i+1])
			} else {
				if str[i] == '"' {
					inString = true
				}
			}
		}
	}
	row.cells = cells
	return row
}

/*
Get Cell value
*/
func (row *Row) Cell(idx int) string {
	return row.cells[idx]
}

/*
To String!
*/
func (row *Row) ToString() string {
	return strings.Join(row.cells, ",")
}
