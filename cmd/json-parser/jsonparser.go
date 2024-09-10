package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

type tokenType int

const (
	// Token types
	OBJECT_START tokenType = iota
	OBJECT_END
	ARRAY_START
	ARRAY_END
	COMMA
	COLON
	STRING
	NUMBER
	TRUE
	FALSE
	NULL
)

func (t tokenType) String() string {
	switch t {
	case OBJECT_START:
		return "OBJECT_START"
	case OBJECT_END:
		return "OBJECT_END"
	case ARRAY_START:
		return "ARRAY_START"
	case ARRAY_END:
		return "ARRAY_END"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	default:
		return "UNKNOWN"
	}
}

type token struct {
	tokenType tokenType
	value     []byte
}

type lexer struct {
	reader io.Reader
	buf    []byte
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{reader, nil}
}

func (l *lexer) getNextToken() (token, error) {
	if len(l.buf) == 0 {
		buf := make([]byte, 32)
		n, err := l.reader.Read(buf)
		if err != nil {
			return token{}, err
		}

		l.buf = buf[:n]
	}

	for i := 0; i < len(l.buf); i++ {
		switch l.buf[i] {
		case '{':
			l.buf = l.buf[i+1:]
			return token{tokenType: OBJECT_START}, nil
		case '}':
			l.buf = l.buf[i+1:]
			return token{tokenType: OBJECT_END}, nil
		}
	}

	return token{}, io.EOF
}

func main() {
	if len(os.Args) < 2 {
		panic("not yet implemented. TODO: Read from stdin")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tokens := make([]token, 0, 10)
	lexer := newLexer(file)
	for {
		token, err := lexer.getNextToken()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		tokens = append(tokens, token)
	}

	if err := parse(tokens); err != nil {
		panic(err)
	}
}

func parse(tokens []token) error {
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens to parse")
	}

	currentTokenType := tokens[0].tokenType
	possibleNextTokens := map[tokenType][]tokenType{
		OBJECT_START: {OBJECT_START, OBJECT_END, STRING},
		OBJECT_END:   {OBJECT_END, COMMA},
		STRING:       {COLON},
	}

	for i := 1; i < len(tokens); i++ {
		nextTokenType := tokens[i].tokenType
		if !slices.Contains(possibleNextTokens[currentTokenType], nextTokenType) {
			return fmt.Errorf("unexpected token %s after %s", nextTokenType, currentTokenType)
		}

		currentTokenType = nextTokenType
	}

	return nil
}
