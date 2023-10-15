package framework

import "github.com/VictoriaMetrics/metrics"

var (
	loginCount  = metrics.NewCounter(`login_count`)
	logoutCount = metrics.NewCounter(`logout_count`)

	userCount = metrics.NewCounter(`user_count`)

	scoreChangesSendCount = metrics.NewCounter(`score_changes_send_count`)

	welcomeAuthorizedActionErrorCount = metrics.NewCounter(`error_count{type="WelcomeAuthorizedAction"}`)
	scoreChangeActionErrorCount       = metrics.NewCounter(`error_count{type="ScoreChangeAction"}`)
)
