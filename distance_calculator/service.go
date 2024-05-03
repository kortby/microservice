package main

import (
	"math"

	"github.com/godev/tolls/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoints []float64 
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}



func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoints) > 0 {
		// prevPoint := s.prevPoints[len(s.prevPoints) - 1]
		distance = calculateDistance(s.prevPoints[0], s.prevPoints[1], data.Lat, data.Lng)
	}
	s.prevPoints = []float64{data.Lat, data.Lng}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2 - x1, 2) + math.Pow(y2 - y1, 2))
}