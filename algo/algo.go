package algo

import "math"

var k float64
var d float64
var s float64

func SetParams(kFactor, derivation, score float64) {
	k = kFactor
	d = derivation
	s = score
}

func NewRatings(ratingW, ratingL float64) (w float64, l float64) {
	expW := 1 / (1 + math.Pow(10, (ratingL-ratingW)/d)) //expected score of winner
	expL := 1 / (1 + math.Pow(10, (ratingW-ratingL)/d)) //expected score of loser
	w = ratingW + (k * (s - float64(expW)))
	l = ratingL - (k * (s - float64(expL)))
	// allow ratings to go negative ?
	// if l < 0 {
	// 	l = 0
	// }
	return w, l
}

var n int = 0

func recurse(w, l float64) {
	if l < 0 {
		return
	}
	a, b := NewRatings(w, l)
	n++
	recurse(a, b)
}
