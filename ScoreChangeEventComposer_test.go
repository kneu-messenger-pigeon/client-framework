package framework

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScoreChangeEventComposer_Compose(t *testing.T) {
	expectedExpire := time.Minute
	updatedAt := time.Date(2028, time.Month(11), 18, 14, 30, 40, 0, time.Local)

	newMiniRedis := func(t *testing.T) *redis.Client {
		return redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    miniredis.RunT(t).Addr(),
		})
	}

	t.Run("empty_storage_single_event", func(t *testing.T) {
		t.Run("create_and_delete_score", func(t *testing.T) {
			createdEvent := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     150,
					LessonPart:   1,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			createdScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(createdEvent.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			}

			deletedEvent := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     150,
					LessonPart:   1,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  false,
						IsDeleted: true,
					},
					UpdatedAt: updatedAt.Add(time.Second),
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			deletedScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(createdEvent.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualAfterCreated := composer.Compose(createdEvent, createdScore)
			fmt.Printf("%v\n", actualAfterCreated)
			actualAfterDeleted := composer.Compose(deletedEvent, deletedScore)
			fmt.Printf("%v\n", actualAfterDeleted)
			assert.True(t, actualAfterDeleted.IsDeleted())
			assert.True(t, deletedScore.IsEqual(actualAfterDeleted))
			redisS := composer.redis.HGetAll(context.Background(), "SC:150:123").Val()
			fmt.Printf("%v\n", redisS)

			assert.Empty(t, out.String())
		})

		t.Run("only_first_lesson_part__score_with_empty_storage", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     150,
					LessonPart:   1,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     1,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson:     currentScore.Lesson,
				FirstScore: &event.Previous.Value,
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualPreviousScore := composer.Compose(event, currentScore)
			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Empty(t, out.String())
		})

		t.Run("only_second_lesson_part_score_with_empty_storage", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   2,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     1,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				SecondScore: floatPointer(2.5),
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson:      currentScore.Lesson,
				SecondScore: &event.Previous.Value,
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Empty(t, out.String())
		})

		t.Run("both_lesson_parts_scores_with_empty_storage", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   2,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     1.5,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore:  floatPointer(0.5),
				SecondScore: floatPointer(1),
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson:      currentScore.Lesson,
				FirstScore:  floatPointer(0.5),
				SecondScore: floatPointer(1.5),
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Empty(t, out.String())
		})

		t.Run("add_absent_event_with_empty_storage", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   2,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						IsAbsent:  true,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				IsAbsent: true,
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson: currentScore.Lesson,
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Empty(t, out.String())
			assert.NotEqual(t, expectedPreviousScore.IsAbsent, currentScore.IsAbsent)
			assert.NotEqual(t, expectedPreviousScore.IsDeleted(), currentScore.IsDeleted())
		})

		t.Run("remove_absent_event_with_empty_storage", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   1,
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					IsAbsent:  true,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2),
				IsAbsent:   false,
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson:   currentScore.Lesson,
				IsAbsent: true,
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Empty(t, out.String())
		})
	})

	t.Run("error", func(t *testing.T) {

		t.Run("wrong_lesson_part", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   5,
					DisciplineId: 234,
					UpdatedAt:    updatedAt,
				},
				Previous: events.ScoreValue{
					Value:     2,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson:     scoreApi.Lesson{},
				FirstScore: floatPointer(2),
				IsAbsent:   false,
			}

			expectedPreviousScore := scoreApi.Score{}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out: out,
			}

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			assert.Contains(t, out.String(), "Wrong lesson part")
		})

		t.Run("redis_error", func(t *testing.T) {
			event := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    143,
					LessonId:     154,
					LessonPart:   1,
					DisciplineId: 234,
					UpdatedAt:    updatedAt,
				},
				Previous: events.ScoreValue{
					Value:     2,
					IsDeleted: false,
				},
			}

			currentScore := &scoreApi.Score{
				Lesson:     scoreApi.Lesson{},
				FirstScore: floatPointer(2),
				IsAbsent:   false,
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson:      currentScore.Lesson,
				FirstScore:  currentScore.FirstScore,
				SecondScore: currentScore.SecondScore,
				IsAbsent:    false,
			}

			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			composer.redis.Set(context.Background(), composer.getStorageKey(event), "string-not-hash-value", 0)

			actualPreviousScore := composer.Compose(event, currentScore)

			assert.Equal(t, expectedPreviousScore, *actualPreviousScore)
			expectedError := "Redis error while composing changed scores: error WRONGTYPE Operation against a key holding the wrong kind of value, storageKey: SC:154:143"
			assert.Equal(t, out.String(), expectedError)
		})
	})

	t.Run("multiple_events", func(t *testing.T) {
		t.Run("score1_then_score2_then_score1_remove_score1_remove_score2", func(t *testing.T) {
			composer := ScoreChangeEventComposer{
				out:           &bytes.Buffer{},
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			lesson := scoreApi.Lesson{
				Id:   150,
				Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
				Type: scoreApi.LessonType{
					Id:        5,
					ShortName: "МК",
					LongName:  "Модульний контроль.",
				},
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson: lesson,
			}

			/** Start step1 */
			step1CurrentScore := &scoreApi.Score{
				Lesson:     lesson,
				FirstScore: floatPointer(2.5),
			}

			step1Score1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step1ActualPreviousScore := composer.Compose(step1Score1, step1CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step1ActualPreviousScore)
			/** End step1 */

			/** Start step2 */
			step2CurrentScore := &scoreApi.Score{
				Lesson:      lesson,
				FirstScore:  floatPointer(2.5),
				SecondScore: floatPointer(1.5),
			}

			step2Score2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     1.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step2ActualPreviousScore := composer.Compose(step2Score2, step2CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step2ActualPreviousScore)
			/** End step2 */

			/** Start step3 */
			step3CurrentScore := &scoreApi.Score{
				Lesson: lesson,
			}

			step3RemoveScore1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  false,
						IsDeleted: true,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step1Score1.Value,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			step3ActualPreviousScore := composer.Compose(step3RemoveScore1, step3CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step3ActualPreviousScore)
			/** End step3 */

			/** Start step4 */
			step4CurrentScore := &scoreApi.Score{
				Lesson: lesson,
			}

			step4RemoveScore2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  false,
						IsDeleted: true,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step2Score2.Value,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			step4ActualPreviousScore := composer.Compose(step4RemoveScore2, step4CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step4ActualPreviousScore)
			/** End step4 */
		})

		t.Run("score1_then_score2_then_absent_score2", func(t *testing.T) {
			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			lesson := scoreApi.Lesson{
				Id:   150,
				Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
				Type: scoreApi.LessonType{
					Id:        5,
					ShortName: "МК",
					LongName:  "Модульний контроль.",
				},
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson: lesson,
			}

			/** Start step1 */
			step1CurrentScore := &scoreApi.Score{
				Lesson:     lesson,
				FirstScore: floatPointer(2.5),
			}

			step1Score1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step1ActualPreviousScore := composer.Compose(step1Score1, step1CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step1ActualPreviousScore)
			/** End step1 */

			/** Start step2 */
			step2CurrentScore := &scoreApi.Score{
				Lesson:      lesson,
				FirstScore:  floatPointer(2.5),
				SecondScore: floatPointer(1.5),
			}

			step2Score2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     1.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step2ActualPreviousScore := composer.Compose(step2Score2, step2CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step2ActualPreviousScore)
			/** End step2 */

			/** Start step3 */
			step3CurrentScore := &scoreApi.Score{
				Lesson:     lesson,
				FirstScore: floatPointer(2.5),
				IsAbsent:   true,
			}

			step3RemoveScore1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  true,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step1Score1.Value,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			step3ActualPreviousScore := composer.Compose(step3RemoveScore1, step3CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step3ActualPreviousScore)
			/** End step3 */
		})

		t.Run("absent1_then_absent2_remove_absent1_remove_absent2", func(t *testing.T) {
			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			lesson := scoreApi.Lesson{
				Id:   150,
				Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
				Type: scoreApi.LessonType{
					Id:        5,
					ShortName: "МК",
					LongName:  "Модульний контроль.",
				},
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson: lesson,
			}

			/** Start step1 */
			step1CurrentScore := &scoreApi.Score{
				Lesson:   lesson,
				IsAbsent: true,
			}

			step1Score1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  true,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step1ActualPreviousScore := composer.Compose(step1Score1, step1CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step1ActualPreviousScore)
			/** End step1 */

			/** Start step2 */
			step2CurrentScore := &scoreApi.Score{
				Lesson:   lesson,
				IsAbsent: true,
			}

			step2Score2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  true,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step2ActualPreviousScore := composer.Compose(step2Score2, step2CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step2ActualPreviousScore)
			/** End step2 */

			/** Start step3 */
			step3CurrentScore := &scoreApi.Score{
				Lesson: lesson,
			}

			step3RemoveScore1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  false,
						IsDeleted: true,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step1Score1.Value,
					IsAbsent:  true,
					IsDeleted: false,
				},
			}

			step3ActualPreviousScore := composer.Compose(step3RemoveScore1, step3CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step3ActualPreviousScore)
			/** End step3 */

			/** Start step4 */
			step4CurrentScore := &scoreApi.Score{
				Lesson: lesson,
			}

			step4RemoveScore2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  false,
						IsDeleted: true,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step2Score2.Value,
					IsAbsent:  true,
					IsDeleted: false,
				},
			}

			step4ActualPreviousScore := composer.Compose(step4RemoveScore2, step4CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step4ActualPreviousScore)
			/** End step4 */
		})

		t.Run("absent1_then_absent2_then_score1", func(t *testing.T) {

		})

		t.Run("score1_then_absent2_then_absent1", func(t *testing.T) {
			out := &bytes.Buffer{}
			composer := ScoreChangeEventComposer{
				out:           out,
				redis:         newMiniRedis(t),
				storageExpire: expectedExpire,
			}

			lesson := scoreApi.Lesson{
				Id:   150,
				Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
				Type: scoreApi.LessonType{
					Id:        5,
					ShortName: "МК",
					LongName:  "Модульний контроль.",
				},
			}

			expectedPreviousScore := scoreApi.Score{
				Lesson: lesson,
			}

			/** Start step1 */
			step1CurrentScore := &scoreApi.Score{
				Lesson:     lesson,
				FirstScore: floatPointer(2.5),
			}

			step1Score1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     2.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step1ActualPreviousScore := composer.Compose(step1Score1, step1CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step1ActualPreviousScore)
			/** End step1 */

			/** Start step2 */
			step2CurrentScore := &scoreApi.Score{
				Lesson:      lesson,
				FirstScore:  floatPointer(2.5),
				SecondScore: floatPointer(1.5),
			}

			step2Score2 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   2, // score2
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     1.5,
						IsAbsent:  false,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     0,
					IsAbsent:  false,
					IsDeleted: true,
				},
			}

			step2ActualPreviousScore := composer.Compose(step2Score2, step2CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step2ActualPreviousScore)
			/** End step2 */

			/** Start step3 */
			step3CurrentScore := &scoreApi.Score{
				Lesson: lesson,
			}

			step3RemoveScore1 := &events.ScoreChangedEvent{
				ScoreEvent: events.ScoreEvent{
					Id:           112233,
					StudentId:    123,
					LessonId:     uint(lesson.Id),
					LessonPart:   1, // score1
					DisciplineId: 234,
					Year:         2028,
					Semester:     1,
					ScoreValue: events.ScoreValue{
						Value:     0,
						IsAbsent:  true,
						IsDeleted: false,
					},
					UpdatedAt: updatedAt,
					SyncedAt:  updatedAt.Add(time.Second * 3),
				},
				Previous: events.ScoreValue{
					Value:     step1Score1.Value,
					IsAbsent:  false,
					IsDeleted: false,
				},
			}

			step3ActualPreviousScore := composer.Compose(step3RemoveScore1, step3CurrentScore)
			assert.Equal(t, expectedPreviousScore, *step3ActualPreviousScore)
			/** End step3 */
		})

	})
}
