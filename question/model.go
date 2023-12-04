package question

import "time"

type Question struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	TsCreated time.Time `json:"ts_created"`
	TsUpdated time.Time `json:"ts_updated"`
}

type Options struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	Body       string    `json:"body"`
	Rank       int       `json:"rank"`
	Correct    bool      `json:"correct"`
	TsCreated  time.Time `json:"ts_created"`
	TsUpdated  time.Time `json:"ts_updated"`
}
