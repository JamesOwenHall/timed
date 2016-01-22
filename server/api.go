package server

import (
	"encoding/json"

	"github.com/JamesOwenHall/timed/query"
)

func (s *Server) HandleQuery(req *Request) {
	if req.r.Method != "POST" {
		req.BadRequest("invalid HTTP method, use POST")
		return
	}

	queryFormValue := req.r.PostFormValue("query")
	if queryFormValue == "" {
		req.BadRequest(`missing required "query" field`)
		return
	}

	var q query.Query
	if err := json.Unmarshal([]byte(queryFormValue), &q); err != nil {
		req.BadRequest(err.Error())
		return
	}

	result, err := s.context.ExecuteQuery(&q)
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
