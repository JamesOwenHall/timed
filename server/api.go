package server

import (
	"encoding/json"

	"github.com/JamesOwenHall/timed/query"
)

func (s *Server) HandleQuery(req *Request) {
	if req.r.Method != "GET" {
		req.BadRequest("invalid HTTP method, use GET")
		return
	}

	queryFormValue := req.r.FormValue("query")
	if queryFormValue == "" {
		req.BadRequest(`missing required "query" field`)
		return
	}

	q, err := query.Parse(queryFormValue)
	if err != nil {
		req.BadRequest(err.Error())
		return
	}

	result, err := s.context.ExecuteQuery(q)
	if err != nil {
		req.InternalServerError(err.Error())
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		req.InternalServerError(err.Error())
	}

	req.w.Write(response)
}
