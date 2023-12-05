package csv

var RETROFIT_COST_ALIAS_SUFFIX = "-Cost"
var RETROFIT_EFF_ALIAS_SUFFIX = "-Eff"

type Retrofit struct {
	alias           string
	costAlias       string
	efficiencyAlias string
	cost            float32
	efficiency      float32
	reduction       float32
}

func CreateRetrofit(alias string, cost, efficiency, reduction float32) Retrofit {
	var retrofit Retrofit
	retrofit.alias = alias
	retrofit.costAlias = alias + RETROFIT_COST_ALIAS_SUFFIX
	retrofit.efficiencyAlias = alias + RETROFIT_EFF_ALIAS_SUFFIX
	retrofit.cost = cost
	retrofit.efficiency = efficiency
	retrofit.reduction = reduction
	return retrofit
}

/*
Getters
*/
func (retrofit *Retrofit) Cost() float32       { return retrofit.cost }
func (retrofit *Retrofit) Reduction() float32  { return retrofit.reduction }
func (retrofit *Retrofit) Efficiency() float32 { return retrofit.efficiency }
