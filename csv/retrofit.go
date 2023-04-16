package csv

type Retrofit struct {
	Code         string
	EPCIndexDiff float32
	Cost         float32
	Savings      float32
}

func CreateRetrofit(code string, epcIndexDiff, cost, savings float32) Retrofit {
	return Retrofit{Code: code, Cost: cost, Savings: savings, EPCIndexDiff: epcIndexDiff}
}
