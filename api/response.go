package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type statusResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (s *Server) respondNotOK(w http.ResponseWriter, statusCode int, err error) {
	resp := statusResponse{
		Code: statusCode,
		Text: err.Error(),
	}
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Error(fmt.Errorf("marshal error response: %w", err))
		return
	}
}

func (s *Server) respondOK(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if resp == nil {
		return
	}

	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	if err := e.Encode(resp); err != nil {
		logrus.Error(fmt.Errorf("marshal response: %w", err))
	}
}
