package framework

import (
	"encoding/binary"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"math"
	"strconv"
)

func CalculateState(newScore *scoreApi.Score, previousScore *scoreApi.Score) string {
	return strconv.FormatUint(
		uint64(
			binary.LittleEndian.Uint32(
				[]byte{
					scoreToUnt8(newScore.FirstScore, newScore.IsAbsent),
					scoreToUnt8(newScore.SecondScore, newScore.IsAbsent),
					scoreToUnt8(previousScore.FirstScore, previousScore.IsAbsent),
					scoreToUnt8(previousScore.SecondScore, previousScore.IsAbsent),
				},
			),
		),
		36,
	)
}

func scoreToUnt8(score *float32, isAbsent bool) uint8 {
	if score == nil {
		if isAbsent {
			return 191
		}
		return 190
	}

	var scoreInteger uint8
	_scoreInteger, scoreFraction := math.Modf(float64(*score))
	if _scoreInteger < 0 {
		scoreInteger = uint8(256 + _scoreInteger)
	} else {
		scoreInteger = uint8(_scoreInteger)
	}

	return scoreInteger + uint8(scoreFraction*40)
}
