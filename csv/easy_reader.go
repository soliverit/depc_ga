package csv

import (
	"io/ioutil"
	"os"
	"strings"
)

type EasyReader struct{
	path 		string
	delimiter 	string
	rows		[]Row
	headers		Row
}
/*
	Add Row
 */
func(eReader * EasyReader) AddRow(row Row){
	eReader.rows = append(eReader.rows, row)
}

/*
	Get Row
 */
func(eReader *EasyReader) Row(idx int)*Row{
	return &eReader.rows[idx]
}
/*
	Number of Rows
 */
func(eReader *EasyReader)Length()int{
	return len(eReader.rows)
}
/*
	New Table From Headers
 */
func CreateEasyReader(path string,headers []string) *EasyReader{
	var eReader *EasyReader = &EasyReader{}
	eReader.path			= path
	eReader.headers 		= NewRow(headers)
	eReader.rows			= make([]Row, 0)
	return eReader
}

/*
	Parse a CSV to EasyReader
 */
func ParseCSV(path string)*EasyReader{
	var eReader EasyReader 	= EasyReader{path:path}
	eReader.rows 			= make([]Row,0)
	csv, _ 					:= ioutil.ReadFile(path)
	var curPos int 			= 0

	/*
		Read headers
	 */
	for i := 0; ; i++{
		if csv[i] == '\n'{
			eReader.headers = StringToRow(string(csv[curPos : i]))
			curPos = i + 1
			break
		}
	}
	/*
		Read rows
	 */
	for i := curPos; i < len(csv) ; i++{
		if csv[i] == '\n'{
			eReader.rows = append(eReader.rows,StringToRow(string(csv[curPos : i])))
			curPos = i + 1
		}
	}
	/*
		Append the last row
	 */
	if curPos < len(csv){
		eReader.rows = append(eReader.rows, StringToRow(string(csv[curPos : len(csv) - curPos])))
	}
	return &eReader
}
func(easyReader * EasyReader) Save(){
	easyReader.WriteColumnsToFile(easyReader.path, easyReader.headers.cells)
}
func(easyReader * EasyReader)WriteColumnsToFile(path string, columns []string){
	output, _ := os.Create(path)
	var headers []int = make([]int, len(columns))
	var foundCount int = 0
	for i := 0; i < len(columns); i++ {
		for x := 0; x < len(easyReader.headers.cells); x++ {
			if easyReader.headers.cells[x] == columns[i]{
				headers[foundCount] = x
				foundCount++
				break
			}
		}
	}
	output.WriteString(strings.Join(columns, ",") + "\n")
	var cells []string = make([]string,len(columns))
	for i := 0; i < easyReader.Length(); i++{
		for x := 0; x < len(columns); x++{
			cells[x] = easyReader.rows[i].cells[headers[x]]
		}
		output.WriteString(strings.Join(cells,",")+ "\n")
	}
	output.Close()
}
func(eReader * EasyReader)ColumnNameToIndex(name string)int{
	for i := 0; i < len(eReader.headers.cells); i++{
		if eReader.headers.cells[i] == name{
			return i
		}
	}
	return -1
}
func(eReader *EasyReader)Join(joinER *EasyReader){
	eReader.headers.cells = append(eReader.headers.cells, joinER.headers.cells...)
	var temp []string
	var rowCount int 	= eReader.Length()
	var joinCount int	= len(joinER.headers.cells)
	for rowID := 0; rowID < rowCount; rowID++{
		temp = make([]string, joinCount)
		for cellID := 0; cellID < len(joinER.headers.cells); cellID++{
			temp[cellID] = joinER.rows[rowID].cells[cellID]
		}
		eReader.rows[rowID].cells = append(eReader.rows[rowID].cells, temp...)
	}
}
func(eReader *EasyReader)CellCount()int{
	return len(eReader.headers.cells)
}
