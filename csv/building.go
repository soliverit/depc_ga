package csv

type Building struct {
	Row
	Id                   int
	Area                 float32
	MinEPCIndexReduction float32
	Retrofits            map[string]Retrofit
	PossibleRetrofits    []string
}

func CreateBuilding(id int, area float32, minEPCIndexReduction float32) Building {
	var building Building
	building.Id = id
	building.Area = area
	building.MinEPCIndexReduction = minEPCIndexReduction
	return building
}
func (building *Building) AddRetrofit(code string, cost, savings, epcDiff float32) {
	building.Retrofits[code] = CreateRetrofit(code, cost, savings, epcDiff)
	// Add to possible Retrofits if it's not a replacement of the code
	found := false
	for _, val := range building.PossibleRetrofits {
		if val == code {
			found = true
			break
		}
	}
	if !found {
		building.PossibleRetrofits = append(building.PossibleRetrofits, code)
	}
}
