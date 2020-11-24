package optimiser
import(
	bayesopt 	"go-bayesopt"
	"strconv"
)

type GAOptimiser struct{
	Parameters []bayesopt.Param
}
func CreateBayesOptimiser() GAOptimiser{
	var gaOptimiser GAOptimiser
	gaOptimiser.Parameters = make([]bayesopt.Param,0)
	return gaOptimiser
}
func(gaOptimiser *GAOptimiser)AddParam(name string, lower float64, upper float64){
	gaOptimiser.Parameters = append(gaOptimiser.Parameters,
		bayesopt.UniformParam{
			Name: name,
			Max:  upper,
			Min:  lower,
		})
}
func(gaOptimiser *GAOptimiser)Run(optimiseMethod func(map[bayesopt.Param]float64)float64){
	var optimiser *bayesopt.Optimizer = bayesopt.New(
		gaOptimiser.Parameters,
		bayesopt.WithRandomRounds(10),
		bayesopt.WithRounds(100),
		bayesopt.WithMinimize(true),
	)

	x, y, _ := optimiser.Optimize(optimiseMethod)
	for param, value := range x{
		print("\n" + param.GetName() + ":\t" +
			strconv.FormatFloat(y,'f',14,32) + "\t" +
			strconv.FormatFloat(value, 'f',14,32))
	}
}
