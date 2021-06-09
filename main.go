package main

import (
	"./csv"
	"./ga"
	"./helpers"
	"./optimiser"
	bayesopt "go-bayesopt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)
const PATH string 		= "C:\\repos\\__shared_data__\\"
const CERTS string 		= "certificates.csv"
const OUTPUT string 	= "indices.csv"
const RETROFITS string	= "retrofits.csv"
const LOCATION string	= PATH + "Middlesbrough\\"
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
	"hotwater",
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
/*
	dhw, dhw-env, dhw-win, dhw-roof
	roof, roof-env, roof-win,
 */
const OPTIMISE bool = true
func main() {
	doWeatherNN()
	return
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
	ph.P("Data loaded - " + strconv.Itoa(data.Length()) + " records")

	/*
		Subtract original ratings from all values
	 */
	var targets csv.BuildingReader 	= csv.ParseBuildingCSV(GA_TARGETS,false)
	var targetLength int			= targets.Length()
	var targetsColumnIDX int		= targets.ColumnNameToIndex("energyEfficiency")
	var subtractColumns	[]int		= make([]int, len(RETROFIT_TARGET_LABELS))
	for i := 0; i < len(subtractColumns); i++{
		subtractColumns[i] = data.ColumnNameToIndex(RETROFIT_TARGET_LABELS[i] + "-Eff")
		/*
			!!! Patch boiler cost !!!
		 */
		if strings.Contains(RETROFIT_TARGET_LABELS[i], "hotwater"){
			for buildingID := 0; buildingID < data.Length(); buildingID++{
				data.Building(buildingID).Subtract( -1200,
					[]int{data.ColumnNameToIndex(RETROFIT_TARGET_LABELS[i] + "-Cost")})
			}
		}
	}
	for i := 0; i < targetLength; i++{
		if data.Building(i).Cell(0) != -9999 {
			data.Building(i).Subtract(targets.Row(i).Cell(targetsColumnIDX),subtractColumns)
		}
	}
	/*
		Remove corrupt Buildings
	 */
	data.RemoveCorrupt()
	ph.P("Data cleansed - " + strconv.Itoa(data.Length()) + " records")
	data = *data.Sample(0.5)
	if OPTIMISE{
		print("\nSHOE")
		var bayesOpt optimiser.GAOptimiser = optimiser.CreateBayesOptimiser()
		bayesOpt.AddParam("Hardness",0.05,0.15)
		bayesOpt.AddParam("CrossoverRate",0.05,0.2)
		bayesOpt.Run(func(params map[bayesopt.Param]float64)float64{
			/*
				Reset random number

				TODO: Is this serial? How does random state consistent over threads
			 */
			rand.Seed(1)
			var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 99,40, RETROFIT_TARGET_LABELS)
			for param, value := range params{
				switch param.GetName(){
				case "Hardness":
					epcGA.Hardness 		= float32(value)
				case "CrossoverRate":
					epcGA.CrossoverRate = float32(value)
				}
			}
			epcGA.Run()
			/*
				fmt.Println("V:" + strconv.FormatFloat(epcGA.ScoreGAstate(0),'f',2,64))
			 */
			return epcGA.ScoreGAstate(0)
		})

	}else {
		/*====
			From Bayes-Opt

			Hardness:		0.110466
			CrossoverRate:	0.191076
		 */
		/*
			Create newEpcGA
		*/
		var epcGA *ga.EpcGA = ga.CreateEpcGA(&data, 200, 40, RETROFIT_TARGET_LABELS)
		epcGA.Hardness		= 0.110466
		epcGA.CrossoverRate	= 0.191076
		ph.P("EPC-GA instantiated")
		epcGA.Run()
	}
}
func doWeatherNN(){
	/*
		data path
	 */
	var pathPattern string = "c:\\repos\\__shared_data__\\ml_anon\\*\\model_epc.inp"
	paths, _ := filepath.Glob(pathPattern)
	for i := 0; i < len(paths); i++{

	}

}
