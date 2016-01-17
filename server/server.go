package server

import (
	"fmt"
	"net/http"

	"github.com/JamesOwenHall/timed/cassandra"

	"github.com/gocql/gocql"
)

type Server struct {
	sources []cassandra.Source
	server  *http.Server
	mux     *http.ServeMux
}

func NewServer(config *Config) (*Server, error) {
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

	mux := http.NewServeMux()
	return &Server{
		sources: sources,
		server: &http.Server{
			Addr:    config.Listen,
			Handler: mux,
		},
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
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
