package query

import (
	"time"
)

type Query struct {
	Source string         `json:"source"`
	Start  time.Time      `json:"start"`
	End    time.Time      `json:"end"`
	Calls  []FunctionCall `json:"compute"`
}

type FunctionCall struct {
	Function string `json:"function"`
	Argument string `json:"argument"`
}
