package ga

type IScorer interface{
	Score(gaState GAState) float32
}
