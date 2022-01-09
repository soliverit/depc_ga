package csv

import (
	"io/ioutil"
	"os"
	"strings"
)

type EasyReader struct {
	path      string
	delimiter string
	rows      []Row
	headers   Row
	Headers   Row
}

/*
	New Table From Headers
*/
func CreateEasyReader(path string, headers []string) *EasyReader {
	var eReader *EasyReader = &EasyReader{}
	eReader.path = path
	eReader.headers = NewRow(headers)
	eReader.Headers = eReader.headers
	eReader.rows = make([]Row, 0)

	return eReader
}

/*
	Parse a CSV to EasyReader
*/
func ParseCSV(path string) *EasyReader {
	var eReader EasyReader = EasyReader{path: path}
	eReader.rows = make([]Row, 0)
	csv, _ := ioutil.ReadFile(path)
	var curPos int = 0

	/*
		Read headers
	*/
	for i := 0; ; i++ {
		if csv[i] == '\n' {
			eReader.headers = StringToRow(string(csv[curPos:i]))
			eReader.Headers = eReader.headers
			//print("..")
			//println(string(csv[curPos:i]))

			curPos = i + 1
			break
		}
	}
	/*
		Read rows
	*/
	for i := curPos; i < len(csv); i++ {
		if csv[i] == '\n' {
			eReader.rows = append(eReader.rows, StringToRow(string(csv[curPos:i])))
			curPos = i + 1
		}
	}
	/*
		Append the last row
	*/
	if curPos < len(csv) {
		eReader.rows = append(eReader.rows, StringToRow(string(csv[curPos:len(csv)-curPos])))
	}
	return &eReader
}
func (easyReader *EasyReader) Save() {
	easyReader.WriteColumnsToFile(easyReader.path, easyReader.headers.cells)
}

func (easyReader *EasyReader) WriteColumnsToFile(path string, columns []string) {
	output, _ := os.Create(path)
	var headers []int = make([]int, len(columns))
	var foundCount int = 0
	for i := 0; i < len(columns); i++ {
		for x := 0; x < len(easyReader.headers.cells); x++ {
			if easyReader.headers.cells[x] == columns[i] {
				headers[foundCount] = x
				foundCount++
				break
			}
		}
	}
	output.WriteString(strings.Join(columns, ",") + "\n")
	var cells []string = make([]string, len(columns))
	for i := 0; i < easyReader.Length(); i++ {
		for x := 0; x < len(columns); x++ {
			cells[x] = easyReader.rows[i].cells[headers[x]]
		}
		output.WriteString(strings.Join(cells, ",") + "\n")
	}
	output.Close()
}
func (easyReader *EasyReader) ColumnNameToIndex(name string) int {
	for i := 0; i < len(easyReader.headers.cells); i++ {
		if easyReader.headers.cells[i] == name {
			return i
		}
	}
	return -1
}
func (easyReader *EasyReader) Join(joinER *EasyReader) {
	easyReader.headers.cells = append(easyReader.headers.cells, joinER.headers.cells...)
	var temp []string
	var rowCount int = easyReader.Length()
	var joinCount int = len(joinER.headers.cells)
	for rowID := 0; rowID < rowCount; rowID++ {
		temp = make([]string, joinCount)
		for cellID := 0; cellID < len(joinER.headers.cells); cellID++ {
			temp[cellID] = joinER.rows[rowID].cells[cellID]
		}
		easyReader.rows[rowID].cells = append(easyReader.rows[rowID].cells, temp...)
	}
}
func (easyReader *EasyReader) CellCount() int {
	return len(easyReader.headers.cells)
}

/*
	Get Row
*/
func (easyReader *EasyReader) Row(idx int) *Row {
	return &easyReader.rows[idx]
}

/*
	Add Row
*/
func (easyReader *EasyReader) AddRow(row Row) {
	easyReader.rows = append(easyReader.rows, row)
}

/*
	Number of Rows
*/
func (easyReader *EasyReader) Length() int {
	return len(easyReader.rows)
}
