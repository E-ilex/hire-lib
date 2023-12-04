package question

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

var (
	//go:embed queries/insert_question.sql
	insertQuestion string

	//go:embed queries/list_questions.sql
	listQuestion string

	//go:embed queries/question_existence.sql
	checkQuestionExistence string

	//go:embed queries/update_question.sql
	updateQuestion string

	//go:embed queries/select_question.sql
	selectQuestion string

	//go:embed queries/delete_option.sql
	deleteOption string

	//go:embed queries/delete_question.sql
	deleteQuestion string
)

type QuestionService interface {
	ListQuestions(ctx context.Context, pagination *Pagination) ([]QuestionRaw, error)
	CreateQuestion(ctx context.Context, question *QuestionRaw) (*Question, error)
	UpdateQuestion(ctx context.Context, questionID string, question *QuestionRaw) (*Question, error)
	DeleteQuestion(ctx context.Context, questionID string) error
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

func (q *questionService) ListQuestions(ctx context.Context, pagination *Pagination) ([]QuestionRaw, error) {
	rows, err := q.database.QueryContext(ctx, listQuestion, pagination.LastID, pagination.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []QuestionRaw
	for rows.Next() {
		var question QuestionRaw
		var jsonData []byte
		err = rows.Scan(
			&question.Body,
			&jsonData,
		)

		// due to sqlite not supporting boolean variable
		var opts OptionsRaw
		if err := json.Unmarshal(jsonData, &opts); err != nil {
			return nil, err
		}

		for _, opt := range opts {
			question.Options = append(question.Options, OptionRaw{opt.Body, opt.Correct != 0})
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (q *questionService) CreateQuestion(ctx context.Context, question *QuestionRaw) (*Question, error) {
	tx, err := q.database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var qstn Question
	err = tx.QueryRowContext(ctx, insertQuestion, question.Body).Scan(
		&qstn.ID,
		&qstn.Body,
		&qstn.TsCreated)

	if errors.Is(err, sql.ErrNoRows) || err != nil {
		return nil, err
	}

	insertOption, values := getOptionQuery(qstn.ID, question.Options)
	_, err = tx.ExecContext(ctx, insertOption, values...)
	if err != nil {
		return nil, err
	}

	return &qstn, nil
}

func (q *questionService) UpdateQuestion(ctx context.Context, questionID string, question *QuestionRaw) (*Question, error) {
	tx, err := q.database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var exists bool
	if err := tx.QueryRowContext(ctx, checkQuestionExistence, questionID).Scan(&exists); err != nil {
		return nil, err
	}

	if !exists {
		return nil, errQuestionNotFound
	}

	_, err = tx.ExecContext(ctx, updateQuestion, question.Body, questionID)
	if err != nil {
		return nil, err
	}

	var qstn Question
	err = tx.QueryRowContext(ctx, selectQuestion, questionID).Scan(
		&qstn.ID,
		&qstn.Body,
		&qstn.TsCreated)

	if errors.Is(err, sql.ErrNoRows) || err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, deleteOption, questionID)
	if err != nil {
		return nil, err
	}

	insertOption, values := getOptionQuery(qstn.ID, question.Options)
	_, err = tx.ExecContext(ctx, insertOption, values...)
	if err != nil {
		return nil, err
	}

	return &qstn, nil
}

func (q *questionService) DeleteQuestion(ctx context.Context, questionID string) error {
	tx, err := q.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var exists bool
	if err := tx.QueryRowContext(ctx, checkQuestionExistence, questionID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return errQuestionNotFound
	}

	_, err = tx.ExecContext(ctx, deleteOption, questionID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, deleteQuestion, questionID)
	if err != nil {
		return err
	}

	return nil
}

func getOptionQuery(questionID int, options []OptionRaw) (string, []interface{}) {
	var values []interface{}
	var placeholders []string
	for _, option := range options {
		values = append(values, questionID, option.Body, option.Correct)
		placeholders = append(placeholders, "(?, ?, ?)")
	}

	insertOption := fmt.Sprintf("INSERT INTO option (question_id, body, correct) VALUES %s",
		strings.Join(placeholders, ","))

	return insertOption, values
}
