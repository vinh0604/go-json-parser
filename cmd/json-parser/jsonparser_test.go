package main

import (
	"strings"
	"testing"
)

func TestGetNextToken(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{}`))

	token, err := lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != OBJECT_START {
		t.Errorf("Expected OBJECT_START, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != OBJECT_END {
		t.Errorf("Expected OBJECT_END, got %s", token.tokenType)
	}
}

func TestParse(t *testing.T) {
	tokens := []token{
		{tokenType: OBJECT_START},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if err := parse([]token{}); err == nil {
		t.Errorf("Expected error, got nil")
	}
}
