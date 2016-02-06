package query

import (
	"fmt"
	"strings"
)

type TokenType int8

const (
	Unknown TokenType = iota
	Compute
	From
	Since
	Until
	Identifier
	OpenParen
	CloseParen
	Comma
	Timestamp
)

func (t TokenType) String() string {
	switch t {
	default:
		return "unknown"
	case Compute:
		return "compute"
	case From:
		return "from"
	case Since:
		return "since"
	case Until:
		return "until"
	case Identifier:
		return "identifier"
	case OpenParen:
		return "open-paren"
	case CloseParen:
		return "close-paren"
	case Comma:
		return "comma"
	case Timestamp:
		return "timestamp"
	}
}

type Token struct {
	Type TokenType
	Raw  string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Raw: \"%s\"},\n", t.Type.String(), t.Raw)
}

type Tokenizer struct {
	input []rune
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input: []rune(input),
	}
}

func (t *Tokenizer) Collect() []Token {
	tokens := []Token{}
	for tok := t.Next(); tok != nil; tok = t.Next() {
		tokens = append(tokens, *tok)
	}

	return tokens
}

func (t *Tokenizer) Next() *Token {
	// Skip spaces.
	for {
		if !isSpace(t.current()) {
			break
		}

		t.input = t.input[1:]
	}

	current := t.current()
	if current == -1 {
		return nil
	}

	// Punctuation.
	switch current {
	case '(':
		t.input = t.input[1:]
		return &Token{Type: OpenParen, Raw: "("}
	case ')':
		t.input = t.input[1:]
		return &Token{Type: CloseParen, Raw: ")"}
	case ',':
		t.input = t.input[1:]
		return &Token{Type: Comma, Raw: ","}
	}

	// Keywords and identifiers.
	if isLetter(current) || current == '_' {
		word := t.readWord()
		switch strings.ToUpper(word) {
		case "COMPUTE":
			return &Token{Type: Compute, Raw: word}
		case "FROM":
			return &Token{Type: From, Raw: word}
		case "SINCE":
			return &Token{Type: Since, Raw: word}
		case "UNTIL":
			return &Token{Type: Until, Raw: word}
		default:
			return &Token{Type: Identifier, Raw: word}
		}
	}

	// Timestamps.
	if current == '\'' {
		return t.readTimestamp()
	}

	return &Token{Type: Unknown, Raw: string([]rune{current})}
}

func (t *Tokenizer) current() rune {
	if len(t.input) == 0 {
		return -1
	}

	return t.input[0]
}

func (t *Tokenizer) readWord() string {
	word := []rune{}

	for {
		r := t.current()
		if !isLetter(r) && r != '_' && !isDigit(r) {
			return string(word)
		}

		word = append(word, r)
		t.input = t.input[1:]
	}
}

func (t *Tokenizer) readTimestamp() *Token {
	ts := []rune{t.current()}
	t.input = t.input[1:]

	for {
		r := t.current()
		if r == -1 {
			return &Token{Type: Unknown, Raw: string(ts)}
		}

		ts = append(ts, r)
		t.input = t.input[1:]

		if r == '\'' {
			return &Token{Type: Timestamp, Raw: string(ts)}
		}
	}
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
