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
	reader   io.Reader
	buf      []byte
	readNext bool
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{reader, nil, true}
}

func (l *lexer) getNextToken() (token, error) {
	if l.readNext {
		buf := make([]byte, 16)
		n, err := l.reader.Read(buf)
		if err != nil {
			return token{}, err
		}

		l.buf = append(l.buf, buf[:n]...)
	}

	l.readNext = false
	firstChar := l.buf[0]
	switch firstChar {
	case '{':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: OBJECT_START}, nil
	case '}':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: OBJECT_END}, nil
	case '[':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: ARRAY_START}, nil
	case ']':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: ARRAY_END}, nil
	case '"':
		strVal := make([]byte, 0, 8)
		j := 1
		escape := false
		for ; j < len(l.buf); j++ {
			if l.buf[j] == '"' && !escape {
				break
			}
			strVal = append(strVal, l.buf[j])
			if l.buf[j] == '\\' {
				escape = !escape
			} else {
				escape = false
			}
		}

		if j < len(l.buf) {
			l.buf = l.buf[j+1:]
			if len(l.buf) == 0 {
				l.readNext = true
			}
			return token{tokenType: STRING, value: strVal}, nil
		} else {
			l.readNext = true
			return l.getNextToken()
		}
	case ':':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: COLON}, nil
	case ',':
		l.buf = l.buf[1:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return token{tokenType: COMMA}, nil
	case ' ', '\n', '\t', '\r':
		j := 1
		for ; j < len(l.buf); j++ {
			if l.buf[j] != ' ' && l.buf[j] != '\n' && l.buf[j] != '\t' && l.buf[j] != '\r' {
				break
			}
		}
		l.buf = l.buf[j:]
		if len(l.buf) == 0 {
			l.readNext = true
		}
		return l.getNextToken()
	default:
		return token{}, fmt.Errorf("unexpected character %c", firstChar)
	}
}

func main() {
	if len(os.Args) < 2 {
		if err := parseFromReader(os.Stdin); err != nil {
			panic(err)
		}
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := parseFromReader(file); err != nil {
		panic(err)
	}
}

func parseFromReader(reader io.Reader) error {
	tokens := make([]token, 0, 10)
	lexer := newLexer(reader)
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

	return nil
}

func parse(tokens []token) error {
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens to parse")
	}

	currentTokenType := tokens[0].tokenType
	possibleNextTokens := map[tokenType][]tokenType{
		OBJECT_START: {OBJECT_START, OBJECT_END, STRING},
		OBJECT_END:   {OBJECT_END, COMMA, ARRAY_END},
		STRING:       {COLON, COMMA, OBJECT_END},
		COLON:        {STRING, NUMBER, TRUE, FALSE, NULL, OBJECT_START, ARRAY_START},
		COMMA:        {STRING, NUMBER, TRUE, FALSE, NULL, OBJECT_START, ARRAY_START},
		ARRAY_START:  {ARRAY_START, ARRAY_END, OBJECT_START, STRING, NUMBER, TRUE, FALSE, NULL},
		ARRAY_END:    {ARRAY_END, COMMA, OBJECT_END},
	}

	var currentKey []byte
	for i := 1; i < len(tokens); i++ {
		nextTokenType := tokens[i].tokenType

		if !slices.Contains(possibleNextTokens[currentTokenType], nextTokenType) {
			return fmt.Errorf("unexpected token %s after %s", nextTokenType, currentTokenType)
		}

		if currentKey == nil {
			if currentTokenType == STRING {
				if nextTokenType != COLON {
					return fmt.Errorf("expected COLON after key %s, got %s", currentKey, nextTokenType)
				}
				currentKey = tokens[i-1].value
			}
		} else {
			if currentTokenType == STRING {
				if nextTokenType != COMMA && nextTokenType != OBJECT_END {
					return fmt.Errorf("expected COMMA or OBJECT_END after value %s, got %s", currentKey, nextTokenType)
				}
				currentKey = nil
			}
		}

		currentTokenType = nextTokenType
	}

	return nil
}
