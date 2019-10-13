package model

type DailyCount struct {
	Date  string `json:"date"  meddler:"date"`
	Count int    `json:"count" meddler:"count"`
}

type StatsOverview struct {
	NewDiscussion int `json:"new_discussion" meddler:"discussions"`
	NewUser       int `json:"new_user" meddler:"users"`

	ActiveUser int `json:"active_user" meddler:"active"`
	TotalUser  int `json:"total_user"  meddler:"total"`
}
