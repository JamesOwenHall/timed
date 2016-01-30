// Example query:
//     COMPUTE sum(f1), avg(f2) FROM table1 SINCE 2015-01-01T00:00:00Z UNTIL 2016-01-01T00:00:00Z
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
