package csv

import (
	"strconv"
	"strings"
)

type Building struct {
	Row
	cells     []float32
	epcIndex  float32
	co2       float32
	retrofits []Retrofit
}

func CreateBuilding(lineString string) Building {
	var building Building
	var tempFloat float64
	var tempString string
	/*
		Iterate values converting to float32. Fills blanks with -1
	*/
	var cells []float32 = make([]float32, 0)
	var curPos uint32 = 0
	var inString bool = false

	/*
		Parse cells
	*/
	var stringLength uint32 = uint32(len(lineString))
	for i := uint32(0); i < stringLength; i++ {
		/*
			Deal with comma encapsulated
		*/
		if inString {
			for lineString[i] != '"' {
				i += 1
			}
			inString = false
		} else {
			if lineString[i] == ',' || i+1 == stringLength {
				tempString = lineString[curPos:i]
				if tempString == "" {
					cells = append(cells, -9999) //Think this injects -9999 for blank (corrupt) cells
				} else {
					tempFloat, _ = strconv.ParseFloat(tempString, 32)
					cells = append(cells, float32(tempFloat))
				}
				curPos = 1 + i
			} else {
				if lineString[i] == '"' {
					inString = true
				}
			}
		}
	}
	/*
		Do last value
	*/

	/*
		Option headers array: This will hold the retrofit options aliases that
		have an impact on the Building.
	*/
	building.retrofits = make([]Retrofit, 0)
	building.AddRetrofit(CreateRetrofit("as-built", 0, 0, 0))
	/*
		return point and... set cells
	*/
	building.cells = cells
	return building
}

/*
Option Aliases: Filter out the retrofits that do nothing for the building, leaving
only the labels that have an impact. Uses *-cost.
*/
func (building *Building) AddRetrofit(retrofit Retrofit) {
	building.retrofits = append(building.retrofits, retrofit)
}

/*
Getters
*/
func (building *Building) EPCIndex() float32         { return building.epcIndex }
func (building *Building) CO2() float32              { return building.co2 }
func (building *Building) NumberOfRetrofits() int    { return len(building.retrofits) }
func (building *Building) Retrofit(id int) *Retrofit { return &building.retrofits[id] }

/*
Set the original co2, energy efficiency and energy consumption (stored in a separate file from retrofits.csv)
*/
func (building *Building) SetOriginalProperties(index, co2 float32) {
	building.epcIndex = index
	building.co2 = co2
}
func (building *Building) Cell(idx int) float32 {
	return building.cells[idx]
}
func (building *Building) Subtract(value float32, cellIDXs []int) {
	for i := 0; i < len(cellIDXs); i++ {
		if building.cells[cellIDXs[i]] <= 0 {
			building.cells[cellIDXs[i]] = building.cells[cellIDXs[i]] - value
		}
	}
}
func (building *Building) ToString() string {
	var strCells []string = make([]string, len(building.cells))
	for i := 0; i < len(building.cells); i++ {
		strCells[i] = strconv.FormatFloat(float64(building.cells[i]), 'f', 3, 32)
	}
	return strings.Join(strCells, ",")
}
