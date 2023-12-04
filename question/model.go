package question

import "time"

type Question struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	TsCreated time.Time `json:"ts_created"`
}

type Options struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	Body       string    `json:"body"`
	Correct    bool      `json:"correct"`
	TsCreated  time.Time `json:"ts_created"`
}

type QuestionRaw struct {
	Body    string      `json:"body" valid:"required"`
	Options []OptionRaw `json:"options" valid:"required"`
}

type OptionRaw struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}

type OptionsRaw []struct {
	Body    string `json:"body"`
	Correct int    `json:"correct"`
}

type Pagination struct {
	LastID *string
	Limit  *string
}
