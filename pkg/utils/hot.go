package utils

import (
	"math"
	"time"

	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

func ComputeHotScore(views, likes int64, createdAt time.Time) float64 {
	var DF float64
	DF = constants.DecayFactor
	activity := float64(views + likes + 1)
	age := float64(time.Now().Unix() - createdAt.Unix())
	return math.Log10(activity) - age/DF
}

func DefaultComputeHotScore(views, likes int64, createdAt time.Time) float64 {
	var DF float64
	DF = constants.DecayFactor
	activity := float64(views + likes + 1)
	age := float64(time.Now().Unix() - createdAt.Unix())
	result := math.Log10(activity) - age/DF
	return result + constants.DefaultHotScoreDelta
}
