*{{.Discipline.Discipline.Name}}*: {{.Discipline.ScoreRating.Total}}
рейтинг #{{.Discipline.ScoreRating.Rating}}/{{.Discipline.ScoreRating.StudentsCount}}

Загалом по групі: max {{.Discipline.ScoreRating.MaxTotal}}, min {{.Discipline.ScoreRating.MinTotal}}

{{range .Discipline.Scores}}{{date .Lesson.Date}} *{{renderScore .}}* {{if .Lesson.Type.LongName}}_{{.Lesson.Type.LongName}}_{{end}}
{{end}}