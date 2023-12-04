package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct{}

type responseData struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func (r *Response) WriteJson(w http.ResponseWriter, req *http.Request, status int, data interface{}, err interface{}) {
	render.Status(req, status)
	render.JSON(w, req, responseData{
		Data:  data,
		Error: err,
	})
}
