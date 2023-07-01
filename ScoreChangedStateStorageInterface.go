package framework

type ScoreChangedStateStorageInterface interface {
	Set(studentId uint, lessonId uint, state string)
	Get(studentId uint, lessonId uint) string
}
