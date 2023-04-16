package csv

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
 These properties have dedicated properties in the Building object. The
 rest of the CSV is groups three per retrofit (EPC-L2 for example).
*/
var BUILDING_BASE_PROPERTIES = []string{"building_id", "area", "net_needed_epc_index_diff"}

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
		Do header indices (split Building props and []Retrofit columns
	*/
	var areaIndex = 0
	var buildingIDIndex = 0
	var epcDiffIndex = 0
	for idx, header := range bReader.Headers().Cells {
		if header == "area" {
			areaIndex = idx
		} else if header == "building_id" {
			buildingIDIndex = idx
		} else if header == "net_needed_epc_index_diff" {
			epcDiffIndex = idx
		}
	}
	/*
		Read rows
	*/
	for i := curPos; i <= len(csv); i++ {
		if csv[i] == '\n' || i == len(csv) {
			/*
			 We need to handle the last line in this block. Otherwise, we'd
			 recode this for it, outside this loop.
			*/
			if i == len(csv) {
				i -= 1
			}
			/*
				Parse Building then identify all applicable retrofits
			*/
			cells := strings.Split(string(csv[curPos:i]), ",")
			/*
			 Map the Retrofit properties. I did try just creating a Retrofit
			 for each code in a single map but Kept getting "Cannot assign Cost"
			 error.
			*/
			costs := make(map[string]float32)
			savings := make(map[string]float32)
			epcDiffs := make(map[string]float32)
			codes := make([]string, 0)
			for cellID := 0; cellID < len(cells); cellID++ {
				if cellID != areaIndex && cellID != epcDiffIndex && cellID != buildingIDIndex {
					/*
					 Get Reco codes and prepare them
					*/
					header := bReader.Headers().Cells[cellID]
					splitHeader := strings.Split(header, "-")
					code := splitHeader[0] + "-" + splitHeader[1]
					codes = append(codes, code)
					if _, ok := costs[code]; !ok {
						costs[code] = 0
						savings[code] = 0
						epcDiffs[code] = 0
					}
					// Get the cell value as float32
					tempValue, _ := strconv.ParseFloat(cells[cellID], 32)
					value := float32(tempValue)
					if splitHeader[2] == "cost" {
						costs[code] = value
					} else if splitHeader[2] == "savings" {
						savings[code] = value
					} else if splitHeader[2] == "epc_index_diff" {
						epcDiffs[code] = value
					}
				}
			}
			/*
			 Prepare common values
			*/
			area64, _ := strconv.ParseFloat(cells[areaIndex], 32)
			epcDiff64, _ := strconv.ParseFloat(cells[epcDiffIndex], 32)
			buildingID, _ := strconv.Atoi(cells[buildingIDIndex])
			building := CreateBuilding(buildingID, float32(area64), float32(epcDiff64))
			/*
				Add Retrofits
			*/
			// Do nothing
			building.AddRetrofit("AS-IS", 0, 0, 0)
			// The rest
			for _, code := range building.PossibleRetrofits {
				building.AddRetrofit(code, costs[code], savings[code], epcDiffs[code])
			}
			/*
			 Finish up
			*/
			// Add building to data set
			bReader.rows = append(bReader.rows, building)
			// Update csvString position
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
