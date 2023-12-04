package question

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

type QuestionService interface {
	GetQuestions(ctx context.Context) ([]Question, error)
	CreateQuestion(ctx context.Context, question *Question) (*Question, error)
	UpdateQuestion(ctx context.Context, id string, question *Question) (*Question, error)
	DeleteQuestion(ctx context.Context, id string) error
}

type questionService struct {
	database *sql.DB
	logging  *zerolog.Logger
}

func NewQuestionService(logging *zerolog.Logger, database *sql.DB) QuestionService {
	return &questionService{
		database: database,
		logging:  logging,
	}
}

// @TODO add pagination
func (q *questionService) GetQuestions(ctx context.Context) ([]Question, error) {
	return nil, nil
}

func (q *questionService) CreateQuestion(ctx context.Context, question *Question) (*Question, error) {
	return nil, nil
}

func (q *questionService) UpdateQuestion(ctx context.Context, id string, question *Question) (*Question, error) {
	return nil, nil
}

func (q *questionService) DeleteQuestion(ctx context.Context, id string) error {
	return nil
}
