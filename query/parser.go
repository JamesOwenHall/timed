package query

import (
	"fmt"
	"io"
	"time"
)

func Parse(input string) (*Query, error) {
	return (&Parser{}).Parse(input)
}

type Parser struct {
	tokenizer *Tokenizer
	current   *Token
	expected  TokenType
}

func (p *Parser) Parse(input string) (*Query, error) {
	p.tokenizer = NewTokenizer(input)
	q := new(Query)

	// Compute clause.
	if p.expect(Compute) == nil {
		return nil, p.currentError()
	}
	calls, err := p.parseFunctionCallList()
	if err != nil {
		return nil, err
	}
	q.Calls = calls

	// From clause.
	if p.expect(From) == nil {
		return nil, p.currentError()
	}
	source := p.expect(Identifier)
	if source == nil {
		return nil, p.currentError()
	}
	q.Source = source.Raw

	// Since clause.
	if p.expect(Since) == nil {
		return nil, p.currentError()
	}
	since := p.expect(Timestamp)
	if since == nil {
		return nil, p.currentError()
	}
	sinceTime, err := time.Parse(time.RFC3339, since.Raw)
	if err != nil {
		return nil, &ErrUnknownToken{Raw: since.Raw}
	}
	q.Since = sinceTime

	// Until clause.
	if p.expect(Until) == nil {
		return nil, p.currentError()
	}
	until := p.expect(Timestamp)
	if until == nil {
		return nil, p.currentError()
	}
	untilTime, err := time.Parse(time.RFC3339, until.Raw)
	if err != nil {
		return nil, &ErrUnknownToken{Raw: until.Raw}
	}
	q.Until = untilTime

	return q, nil
}

func (p *Parser) parseFunctionCallList() ([]FunctionCall, error) {
	list := []FunctionCall{}

	first, err := p.parseFunctionCall()
	if err != nil {
		return nil, err
	}
	list = append(list, *first)

	for {
		if p.expect(Comma) == nil {
			return list, nil
		}

		fc, err := p.parseFunctionCall()
		if err != nil {
			return nil, err
		}
		list = append(list, *fc)
	}
}

func (p *Parser) parseFunctionCall() (*FunctionCall, error) {
	function := p.expect(Identifier)
	if function == nil {
		return nil, p.currentError()
	}

	if p.expect(OpenParen) == nil {
		return nil, p.currentError()
	}

	arg := p.expect(Identifier)
	if arg == nil {
		return nil, p.currentError()
	}

	if p.expect(CloseParen) == nil {
		return nil, p.currentError()
	}

	return &FunctionCall{
		Function: function.Raw,
		Argument: arg.Raw,
	}, nil
}

// expect returns a token if the next token has the type t.  Otherwise, it
// stores the token in p.current for subsequent calls.
func (p *Parser) expect(t TokenType) *Token {
	p.expected = t

	if p.current == nil {
		p.current = p.tokenizer.Next()
	}

	// Check for EOF.
	if p.current == nil {
		return nil
	}

	if p.current.Type == t {
		tok := p.current
		p.current = nil
		return tok
	}

	return nil
}

func (p *Parser) currentError() error {
	if p.current == nil {
		return io.EOF
	}

	if p.current.Type == Unknown {
		return &ErrUnknownToken{Raw: p.current.Raw}
	}

	return &ErrUnexpectedToken{Exp: p.expected, Got: p.current.Type}
}

type ErrUnexpectedToken struct {
	Exp TokenType
	Got TokenType
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected %s; expected %s", e.Got.String(), e.Exp.String())
}

type ErrUnknownToken struct {
	Raw string
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("unknown token: %s", e.Raw)
}
