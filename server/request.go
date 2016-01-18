package server

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Request struct {
	w   http.ResponseWriter
	r   *http.Request
	log *logrus.Entry
	err string
}

func (s *Server) NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	r.ParseForm()
	return &Request{
		w: w,
		r: r,
		log: s.log.WithFields(logrus.Fields{
			"url":            r.URL.String(),
			"path":           r.URL.Path,
			"method":         r.Method,
			"client_address": r.RemoteAddr,
			"form":           r.Form,
		}),
	}
}

func (r *Request) Error(status int, message string) {
	r.log = r.log.WithField("status", status)
	r.err = message
	r.w.WriteHeader(status)
	r.w.Write([]byte(message))
}

func (r *Request) NotFound(message string) {
	r.Error(http.StatusNotFound, message)
}

func (r *Request) BadRequest(message string) {
	r.Error(http.StatusBadRequest, message)
}

func (r *Request) InternalServerError(message string) {
	r.Error(http.StatusInternalServerError, message)
}
