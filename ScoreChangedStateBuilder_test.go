package framework

import (
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateState(t *testing.T) {
	makePointerFloat := func(value float32) *float32 {
		return &value
	}

	t.Run("idempotence", func(t *testing.T) {
		newScore := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(14.75),
		}

		previousScore := &scoreApi.Score{}

		state1 := CalculateState(newScore, previousScore)
		state2 := CalculateState(newScore, previousScore)
		state3 := CalculateState(newScore, previousScore)
		assert.NotEmpty(t, state1)
		assert.Equal(t, state1, state2)
		assert.Equal(t, state1, state3)
	})

	t.Run("idempotence", func(t *testing.T) {
		newScore := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(14.75),
		}

		previousScore := &scoreApi.Score{}

		state1 := CalculateState(newScore, previousScore)
		state2 := CalculateState(newScore, previousScore)
		state3 := CalculateState(newScore, previousScore)
		assert.NotEmpty(t, state1)
		assert.Equal(t, state1, state2)
		assert.Equal(t, state1, state3)
	})

	t.Run("idempotence_is_absent", func(t *testing.T) {
		newScore := &scoreApi.Score{
			IsAbsent: true,
		}

		previousScore := &scoreApi.Score{}

		state1 := CalculateState(newScore, previousScore)
		state2 := CalculateState(newScore, previousScore)
		state3 := CalculateState(newScore, previousScore)
		assert.NotEmpty(t, state1)
		assert.Equal(t, state1, state2)
		assert.Equal(t, state1, state3)
	})

	t.Run("elasticity", func(t *testing.T) {
		newScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(14.75),
		}

		previousScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(21.35),
			SecondScore: makePointerFloat(7),
		}

		newScore2 := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(15.75),
		}

		previousScore2 := &scoreApi.Score{
			FirstScore:  makePointerFloat(22.35),
			SecondScore: makePointerFloat(7),
		}

		state1 := CalculateState(newScore1, previousScore1)
		state2 := CalculateState(newScore2, previousScore2)
		assert.NotEqual(t, state1, state2)
	})

	t.Run("elasticity_is_absent", func(t *testing.T) {
		newScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(14.75),
		}

		previousScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(21.35),
			SecondScore: makePointerFloat(7),
		}

		newScore2 := &scoreApi.Score{
			IsAbsent: true,
		}

		previousScore2 := &scoreApi.Score{}

		state1 := CalculateState(newScore1, previousScore1)
		state2 := CalculateState(newScore2, previousScore2)
		assert.NotEqual(t, state1, state2)
	})

	t.Run("elasticity_negative_", func(t *testing.T) {
		newScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(-32),
			SecondScore: makePointerFloat(-5),
		}

		previousScore1 := &scoreApi.Score{
			FirstScore:  makePointerFloat(21.35),
			SecondScore: makePointerFloat(7),
		}

		newScore2 := &scoreApi.Score{
			FirstScore:  makePointerFloat(32),
			SecondScore: makePointerFloat(5),
		}

		previousScore2 := &scoreApi.Score{}

		state1 := CalculateState(newScore1, previousScore1)
		state2 := CalculateState(newScore2, previousScore2)
		assert.NotEqual(t, state1, state2)
	})
}
