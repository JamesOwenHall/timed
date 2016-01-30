// Example query:
//     COMPUTE sum(f1), avg(f2) FROM table1 SINCE 2015-01-01T00:00:00Z UNTIL 2016-01-01T00:00:00Z
package query

import (
	"time"
)

type Query struct {
	Source string
	Since  time.Time
	Until  time.Time
	Calls  []FunctionCall
}

type FunctionCall struct {
	Function string
	Argument string
}
