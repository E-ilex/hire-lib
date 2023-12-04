package question

import (
	"hire-test-lib/utils"

	"net/http"

	"github.com/rs/zerolog"
)

type QuestionHandler interface {
	GetQuestions(w http.ResponseWriter, r *http.Request)
	CreateQuestion(w http.ResponseWriter, r *http.Request)
	UpdateQuestion(w http.ResponseWriter, r *http.Request)
	DeleteQuestion(w http.ResponseWriter, r *http.Request)
}

type questionHandler struct {
	questionService QuestionService
	resp            *utils.Response
	logging         *zerolog.Logger
}

func NewQuestionHandler(logging *zerolog.Logger, resp *utils.Response, serv QuestionService) QuestionHandler {
	return &questionHandler{
		questionService: serv,
		resp:            resp,
		logging:         logging,
	}
}

func (q *questionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	// @TODO
	q.resp.WriteJson(w, r, 200, "placeholder", nil)
}

func (q *questionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// @TODO
	q.resp.WriteJson(w, r, 200, "placeholder", nil)
}

func (q *questionHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	// @TODO
	q.resp.WriteJson(w, r, 200, "placeholder", nil)
}

func (q *questionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	// @TODO
	q.resp.WriteJson(w, r, 200, "placeholder", nil)
}
