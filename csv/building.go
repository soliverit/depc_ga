package csv

import (
	"strconv"
	"strings"
)

type Building struct{
	Row
	cells []float32
}
func CreateBuilding(lineString string) Building{
	var building 	Building
	var tempFloat 	float64
	var tempString	string
	/*
		Iterate values converting to float32. Fills blanks with -1
	 */
	var cells 		[]float32 = make([]float32,0)
	var curPos 		uint32= 0
	var inString	bool = false

	/*
		Parse cells
	*/
	var stringLength uint32 = uint32(len(lineString))
	for i := uint32(0); i < stringLength; i++{
		/*
			Deal with comma encapsulated
		*/
		if inString {
			for lineString[i] != '"' {
				i += 1
			}
			inString = false
		}else{
			if lineString[i] == ',' {
				tempString 		= lineString[curPos:i]
				if tempString ==""{
					cells 			= append(cells, -9999)
				}else {
					tempFloat, _ 	= strconv.ParseFloat(tempString, 32)
					cells 			= append(cells, float32(tempFloat))
				}
				curPos 			= 1 + i
			}else{
				if lineString[i] == '"'{
					inString = true
				}
			}
		}
	}
	/*
		return point and... set cells
	 */
	building.cells = cells
	return building
}
func(building *Building)Cell(idx int)float32{
	return building.cells[idx]
}
func(building *Building)Subtract(value float32, cellIDXs []int){
	for i := 0; i < len(cellIDXs); i++{
		if building.cells[cellIDXs[i]] <= 0{
			building.cells[cellIDXs[i]] = building.cells[cellIDXs[i]] - value
		}
	}
}
func(building *Building)ToString()string{
	var strCells []string = make([]string, len(building.cells))
	for i := 0; i < len(building.cells); i++{
		strCells[i]	= strconv.FormatFloat(float64(building.cells[i]),'f',3,32)
	}
	return strings.Join(strCells, ",")
}