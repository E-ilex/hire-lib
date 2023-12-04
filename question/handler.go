package question

import (
	"context"
	"encoding/json"
	"errors"
	"hire-test-lib/utils"
	"time"

	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const (
	contextTimeout             = time.Second * 10
	component                  = "component"
	requestID                  = "requestID"
	respInvalidPayload         = "Invalid payload to create question"
	respFailedToCreateQuestion = "Could not create question"
	respFailedToListQuestion   = "Could not list question"
	respFailedToUpdateQuestion = "Could not update question"
	respFailedToDeleteQuestion = "Could not delete question"
	respFailedToFindQuestion   = "Could not find question"
)

var (
	errFailedToDecodePayload   = errors.New("failed to decode payload")
	errFailedToValidatePayload = errors.New("failed to validate payload")
	errFailedToCreateQuestion  = errors.New("failed to create question")
	errFailedToListQuestion    = errors.New("failed to list question")
	errFailedToUpdateQuestion  = errors.New("failed to update question")
	errFailedToDeleteQuestion  = errors.New("failed to delete question")
	errQuestionNotFound        = errors.New("question not found")
)

type QuestionHandler interface {
	ListQuestions(w http.ResponseWriter, r *http.Request)
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

func (q *questionHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	currentComponent := "question/list-questions"
	reqCtx := r.Context()
	reqID := middleware.GetReqID(reqCtx)

	ctx, cancel := context.WithTimeout(reqCtx, contextTimeout)
	defer cancel()

	pagination := getPagination(r)
	questions, err := q.questionService.ListQuestions(ctx, pagination)
	if err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToListQuestion.Error())
		q.resp.WriteJson(w, r, http.StatusInternalServerError, nil, respFailedToListQuestion)
		return
	}

	q.resp.WriteJson(w, r, http.StatusOK, questions, nil)
}

func (q *questionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	currentComponent := "question/create-question"
	reqCtx := r.Context()
	reqID := middleware.GetReqID(reqCtx)

	var payload QuestionRaw
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToDecodePayload.Error())
		q.resp.WriteJson(w, r, http.StatusBadRequest, nil, respInvalidPayload)
		return
	}

	if _, err := govalidator.ValidateStruct(payload); err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToValidatePayload.Error())
		q.resp.WriteJson(w, r, http.StatusBadRequest, nil, respInvalidPayload)
		return
	}

	ctx, cancel := context.WithTimeout(reqCtx, contextTimeout)
	defer cancel()

	question, err := q.questionService.CreateQuestion(ctx, &payload)
	if err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToCreateQuestion.Error())
		q.resp.WriteJson(w, r, http.StatusInternalServerError, nil, respFailedToCreateQuestion)
		return
	}

	q.resp.WriteJson(w, r, http.StatusCreated, question, nil)
}

func (q *questionHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	currentComponent := "question/update-question"
	reqCtx := r.Context()
	reqID := middleware.GetReqID(reqCtx)

	questionID := chi.URLParam(r, "questionID")

	var payload QuestionRaw
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToDecodePayload.Error())
		q.resp.WriteJson(w, r, http.StatusBadRequest, nil, respInvalidPayload)
		return
	}

	if _, err := govalidator.ValidateStruct(payload); err != nil {
		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToValidatePayload.Error())
		q.resp.WriteJson(w, r, http.StatusBadRequest, nil, respInvalidPayload)
		return
	}

	ctx, cancel := context.WithTimeout(reqCtx, contextTimeout)
	defer cancel()

	question, err := q.questionService.UpdateQuestion(ctx, questionID, &payload)
	if err != nil {
		if errors.Is(err, errQuestionNotFound) {
			q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errQuestionNotFound.Error())
			q.resp.WriteJson(w, r, http.StatusNotFound, nil, respFailedToFindQuestion)
			return
		}

		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToUpdateQuestion.Error())
		q.resp.WriteJson(w, r, http.StatusInternalServerError, nil, respFailedToUpdateQuestion)
		return
	}

	q.resp.WriteJson(w, r, http.StatusOK, question, nil)
}

func (q *questionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	currentComponent := "question/delete-question"
	reqCtx := r.Context()
	reqID := middleware.GetReqID(reqCtx)

	questionID := chi.URLParam(r, "questionID")

	ctx, cancel := context.WithTimeout(reqCtx, contextTimeout)
	defer cancel()

	err := q.questionService.DeleteQuestion(ctx, questionID)
	if err != nil {
		if errors.Is(err, errQuestionNotFound) {
			q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errQuestionNotFound.Error())
			q.resp.WriteJson(w, r, http.StatusNotFound, nil, respFailedToFindQuestion)
			return
		}

		q.logging.Error().Err(err).Str(requestID, reqID).Str(component, currentComponent).Msg(errFailedToDeleteQuestion.Error())
		q.resp.WriteJson(w, r, http.StatusInternalServerError, nil, respFailedToDeleteQuestion)
		return
	}

	q.resp.WriteJson(w, r, http.StatusNoContent, nil, nil)
}

func getPagination(r *http.Request) *Pagination {
	pagination := Pagination{}
	lastID := r.FormValue("last_id")
	if lastID != "" {
		pagination.LastID = &lastID
	}

	limit := r.FormValue("limit")
	if limit != "" {
		pagination.Limit = &limit
	}
	return &pagination
}
