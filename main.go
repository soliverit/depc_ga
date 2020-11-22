package main
import(
	"./csv"
	"./ga"
	"./helpers"
	"math/rand"
	"strconv"
)
const PATH string 		= "C:\\repos\\depc_emulator\\data\\epc_data\\"
const CERTS string 		= "certificates.csv"
const OUTPUT string 	= "indices.csv"
const RETROFITS string	= "retrofits.csv"
const LOCATION string	= PATH + "domestic-E06000002-Middlesbrough\\"
const GA_PATH string	= LOCATION + RETROFITS
const GA_TARGETS string	= LOCATION + "targets.csv"
var RETROFIT_LABELS 	= []string{"envelopes_hotwater_roof_windows",
	"envelopes_hotwater_roof",
	"envelopes_hotwater_windows",
	"envelopes_hotwater",
	"envelopes_roof_windows",
	"envelopes_roof",
	"envelopes_windows",
	"envelopes",
	"hotwater_roof_windows",
	"hotwater_roof",
	"hotwater_windows",
	"hotwater",
	"roof_windows",
	"roof",
	"windows"}
var RETROFIT_TARGET_LABELS = []string{
	"hotwater_envelopes",
	"roof_envelopes",
	"windows_envelopes",
	"envelopes",
	"roof_hotwater",
	"hotwater_windows",
	"hotwater",
	"roof_windows",
	"roof",
	"windows"}

func main() {
	rand.Seed(1)
	/*
		Do helpers
	 */
	var ph *helpers.PrintHelper 	= helpers.CreatePrintHelper(true)
	ph.P("Starting")

	/*
		Parse main data
	 */
	var data	csv.BuildingReader	= csv.ParseBuildingCSV(GA_PATH, false)
	ph.P("Data loaded - " +strconv.Itoa(data.Length()) + " records")

	/*
		Subtract original ratings from all values
	 */
	var targets csv.BuildingReader 	= csv.ParseBuildingCSV(GA_TARGETS,false)
	var targetLength int			= targets.Length()
	var targetsColumnIDX int		= targets.ColumnNameToIndex("energyEfficiency")


	var subtractColumns	[]int		= make([]int, len(RETROFIT_TARGET_LABELS))
	for i := 0; i < len(subtractColumns); i++{
		subtractColumns[i] = data.ColumnNameToIndex(RETROFIT_TARGET_LABELS[i] + "-Eff")
	}
	for i := 0; i < targetLength; i++{
		if data.Building(i).Cell(0) != -9999{
			data.Building(i).Subtract(targets.Row(i).Cell(targetsColumnIDX),subtractColumns)
		}
	}
	/*
		Remove corrupt Buildings
	 */
	data.RemoveCorrupt()
	ph.P("Data cleansed - " + strconv.Itoa(data.Length()) + " records")
	data.WriteToFile("c://repos/GA_TEMP.csv")
	/*
		Create newEpcGA
	*/
	var epcGA 	*ga.EpcGA 			= ga.CreateEpcGA(&data, 100, 5, RETROFIT_TARGET_LABELS)
	ph.P("EPC-GA instantiated")
	epcGA.Run()
}


func P(str string){
	print("--- Message ---")
	print(str)
	print("===============")
}