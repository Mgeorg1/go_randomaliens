// Package service uses for generate random frequencies for random mean and std
package service

import "math/rand"

type FrequenciesGenerator struct {
	Mean float64
	Std  float64
}

func (g *FrequenciesGenerator) init() {
	g.Mean = -10 + rand.Float64()*20
	g.Std = 0.3 + rand.Float64()*(1.5-0.3)
}

func (g FrequenciesGenerator) Generate() float64 {
	return rand.NormFloat64()*g.Std + g.Mean
}

func NewGenerator() *FrequenciesGenerator {
	newGen := &FrequenciesGenerator{}
	newGen.init()
	return newGen
}
