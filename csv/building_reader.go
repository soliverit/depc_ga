package csv

import (
	"io/ioutil"
	"os"
	"strings"
)

type BuildingReader struct {
	EasyReader EasyReader
	rows       []Building
}

func (bReader BuildingReader) Length() int {
	return len(bReader.rows)
}
func ParseBuildingCSV(path string, skipCorruptRecords bool) BuildingReader {
	var bReader BuildingReader = BuildingReader{}
	bReader.EasyReader = EasyReader{}
	bReader.rows = make([]Building, 0)
	csv, _ := ioutil.ReadFile(path)
	var curPos int = 0

	/*
		Read headers
	*/
	for i := 0; ; i++ {
		if csv[i] == '\n' {
			bReader.EasyReader.headers = StringToRow(string(csv[curPos:i]))
			bReader.EasyReader.Headers = StringToRow(string(csv[curPos:i]))
			curPos = i + 1
			break
		}
	}

	/*
		Read rows
	*/
	var tempLine string
	for i := curPos; i < len(csv); i++ {
		if csv[i] == '\n' {
			tempLine = string(csv[curPos:i])
			/*
				Lines with -9999 didn't have enough data. Skip them
			*/
			if skipCorruptRecords && strings.Contains(tempLine, "-9") {
				continue
			}

			/*
				Parse Building then identify all applicable retrofits
			*/
			bReader.rows = append(bReader.rows, CreateBuilding(tempLine))

			/*
				Update counter(s)
			*/
			curPos = i + 1
		}
	}
	/*
		Append the last row
	*/
	if curPos < len(csv) {
		bReader.rows = append(bReader.rows, CreateBuilding(string(csv[curPos:len(csv)-curPos])))
	}
	println("????")
	println(bReader.Headers().ToString())
	println("???")
	return bReader
}

/*
	Get Row
*/
func (bReader *BuildingReader) Row(idx int) *Building {
	return &bReader.rows[idx]
}
func (bReader *BuildingReader) Building(idx int) *Building {
	return &bReader.rows[idx]
}
func (bReader *BuildingReader) RemoveCorrupt() {
	var newBuildings []Building = make([]Building, 0)
	var building *Building
	for i := 0; i < bReader.Length(); i++ {
		building = bReader.Building(i)
		if building.cells[0] != -9999 {
			newBuildings = append(newBuildings, *building)
		}
	}
	bReader.rows = newBuildings
}
func (bReader *BuildingReader) Headers() *Row {
	return &bReader.EasyReader.headers
}
func (bReader *BuildingReader) headers() *Row {
	return bReader.Headers()
}
func (bReader *BuildingReader) WriteToFile(path string) {
	var file *os.File
	file, _ = os.Create(path)
	file.WriteString(bReader.headers().ToString())
	for i := 0; i < bReader.Length(); i++ {
		file.WriteString("\n" + bReader.Building(i).ToString())
	}
	file.Close()
}
func (bReader *BuildingReader) Sample(size float32) *BuildingReader {
	var length int = int(float32(bReader.Length()) * size)
	var output BuildingReader
	output.EasyReader = bReader.EasyReader
	output.rows = make([]Building, length)
	for i := 0; i < length; i++ {
		output.rows[i] = bReader.rows[i]
	}
	return &output
}
