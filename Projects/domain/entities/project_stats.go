//geova-back-1/Projects/domain/entities/project_stats.go
package entities

type DailyProjectCount struct {
    Date  string `json:"date"`
    Count int    `json:"count"`
}

type ProjectStats struct {
    UserId     int                 `json:"user_id"`
    TotalCount int                 `json:"total_count"`
    Daily      []DailyProjectCount `json:"daily"`
}