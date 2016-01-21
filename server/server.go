package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JamesOwenHall/timed/cassandra"
	"github.com/JamesOwenHall/timed/query"

	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
)

type Server struct {
	context query.Context
	server  *http.Server
	log     *logrus.Logger
}

func NewServer(log *logrus.Logger, config *Config) (*Server, error) {
	cluster := gocql.NewCluster(config.Cassandra.Addresses...)
	cluster.Keyspace = config.Cassandra.Keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	sources := []cassandra.Source{}
	for _, src := range config.Sources {
		c, err := parseConsistency(src.Consistency)
		if err != nil {
			return nil, err
		}

		sources = append(sources, cassandra.Source{
			Session:     session,
			Consistency: c,
			Name:        src.Name,
			TimeKey:     src.Timekey,
		})
	}

	server := &Server{
		context: query.Context{Sources: sources},
		server: &http.Server{
			Addr: config.Listen,
		},
		log: log,
	}
	server.server.Handler = server

	return server, nil
}

func (s *Server) ListenAndServe() error {
	s.log.WithField("address", s.server.Addr).Info("server listening")
	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := s.NewRequest(w, r)

	if path := r.URL.Path; path == "/query" {
		s.Query(req)
	} else {
		req.NotFound("unknown endpoint " + path)
	}

	if req.err != "" {
		req.log.Error(req.err)
	} else {
		req.log.Info("success")
	}
}

func (s *Server) Query(req *Request) {
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

func parseConsistency(s string) (c gocql.Consistency, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ErrInvalidConfig{
				Message: "unknown consistency",
				Value:   s,
			}
		}
	}()

	c = gocql.ParseConsistency(s)
	return c, err
}

type ErrInvalidConfig struct {
	Message string
	Value   string
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf(
		"invalid configuration => %s (%s)",
		e.Message, e.Value,
	)
}
