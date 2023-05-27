package framework

import (
	"context"
	"fmt"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/redis/go-redis/v9"
	"io"
	"strconv"
	"time"
)

const IsAbsentScoreValue = "H"
const IsDeletedScoreValue = "D"

const ScoreFirstFieldName = "f"
const ScoreSecondFieldName = "s"

type ScoreChangeEventComposer struct {
	out           io.Writer
	redis         redis.UniversalClient
	storageExpire time.Duration
}

func (composer *ScoreChangeEventComposer) Compose(event *events.ScoreChangedEvent, currentScore *scoreApi.Score) scoreApi.Score {
	ctx := context.Background()
	storageScoreKey := composer.getStorageKey(event)

	/** Start section: Write new previous value to storage if there is no other value */
	var changedScoreFieldName string
	if event.LessonPart == 1 {
		changedScoreFieldName = ScoreFirstFieldName
	} else if event.LessonPart == 2 {
		changedScoreFieldName = ScoreSecondFieldName
	} else {
		_, _ = fmt.Fprintf(composer.out, "Wrong lesson part, storedKey %s, event:  %v", storageScoreKey, event)
		return scoreApi.Score{}
	}

	var scoreValueToSet string
	if event.Previous.IsDeleted {
		scoreValueToSet = IsDeletedScoreValue
	} else if event.Previous.IsAbsent {
		scoreValueToSet = IsAbsentScoreValue
	} else {
		scoreValueToSet = strconv.FormatFloat(float64(event.Previous.Value), 'f', -1, 32)
	}

	result, redisErr := composer.redis.HSetNX(ctx, storageScoreKey, changedScoreFieldName, scoreValueToSet).Result()
	if result {
		redisErr = composer.redis.Expire(ctx, storageScoreKey, composer.storageExpire).Err()
	}

	if redisErr != nil && redisErr != redis.Nil {
		_, _ = fmt.Fprintf(
			composer.out, "Redis error while composing changed scores: error %s, storageKey: %s",
			redisErr.Error(), storageScoreKey,
		)
	}
	/** End section: Write new previous value to storage if there is no other value */

	/** Read value from storage */
	previousScore := scoreApi.Score{
		Lesson: currentScore.Lesson,
		// fallback value to actual current values if there is not previous value
		FirstScore:  currentScore.FirstScore,
		SecondScore: currentScore.SecondScore,
		IsAbsent:    false,
	}
	for key, value := range composer.redis.HGetAll(ctx, storageScoreKey).Val() {
		if IsAbsentScoreValue == value {
			previousScore.IsAbsent = true
		}

		if key == ScoreFirstFieldName {
			previousScore.FirstScore = parseScore(value)
		} else if key == ScoreSecondFieldName {
			previousScore.SecondScore = parseScore(value)
		}
	}
	return previousScore
}

func (composer *ScoreChangeEventComposer) getStorageKey(event *events.ScoreChangedEvent) string {
	return fmt.Sprintf("SC:%d:%d", event.LessonId, event.StudentId)
}

func parseScore(input string) *float32 {
	if input == IsDeletedScoreValue || input == IsAbsentScoreValue || input == "" {
		return nil
	}

	f64, _ := strconv.ParseFloat(input, 10)
	f32 := float32(f64)
	return &f32
}
