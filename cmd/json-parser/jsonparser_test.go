package main

import (
	"io"
	"strings"
	"testing"
)

func TestGetNextTokenEmptyObject(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{}`))
	expected := []token{
		{OBJECT_START, nil},
		{OBJECT_END, nil},
	}

	for _, exp := range expected {
		t.Run("Get "+exp.tokenType.String(), func(t *testing.T) {
			token, err := lexer.getNextToken()
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if token.tokenType != exp.tokenType {
				t.Errorf("Expected %s, got %s", exp.tokenType, token.tokenType)
			}
		})
	}
}

func TestGetNextTokenEmptyArray(t *testing.T) {
	lexer := newLexer(strings.NewReader(`[]`))
	expected := []token{
		{ARRAY_START, nil},
		{ARRAY_END, nil},
	}

	for _, exp := range expected {
		t.Run("Get "+exp.tokenType.String(), func(t *testing.T) {
			token, err := lexer.getNextToken()
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if token.tokenType != exp.tokenType {
				t.Errorf("Expected %s, got %s", exp.tokenType, token.tokenType)
			}
		})
	}
}

func TestGetNextTokenObjectStringValue(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{"key1":"value1"}`))
	expected := []token{
		{OBJECT_START, nil},
		{STRING, []byte("key1")},
		{COLON, nil},
		{STRING, []byte("value1")},
		{OBJECT_END, nil},
	}

	for _, exp := range expected {
		t.Run("Get "+exp.tokenType.String(), func(t *testing.T) {
			token, err := lexer.getNextToken()
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if token.tokenType != exp.tokenType {
				t.Errorf("Expected %s, got %s", exp.tokenType, token.tokenType)
			}
		})
	}
}

func TestGetNextTokenWithSpaceChars(t *testing.T) {
	lexer := newLexer(strings.NewReader("{\t} \r\n"))
	expected := []token{
		{OBJECT_START, nil},
		{OBJECT_END, nil},
	}

	for _, exp := range expected {
		t.Run("Get "+exp.tokenType.String(), func(t *testing.T) {
			token, err := lexer.getNextToken()
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if token.tokenType != exp.tokenType {
				t.Errorf("Expected %s, got %s", exp.tokenType, token.tokenType)
			}
		})
	}

	_, err := lexer.getNextToken()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %s", err)
	}
}

func TestGetNextTokenObjectMultiStringValue(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{"key1":"value1",
	"key2":"value2"}`))
	expected := []token{
		{OBJECT_START, nil},
		{STRING, []byte("key1")},
		{COLON, nil},
		{STRING, []byte("value1")},
		{COMMA, nil},
		{STRING, []byte("key2")},
		{COLON, nil},
		{STRING, []byte("value2")},
		{OBJECT_END, nil},
	}

	for _, exp := range expected {
		t.Run("Get "+exp.tokenType.String(), func(t *testing.T) {
			token, err := lexer.getNextToken()
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if token.tokenType != exp.tokenType {
				t.Errorf("Expected %s, got %s", exp.tokenType, token.tokenType)
			}
		})
	}
}

func TestParseEmptyObject(t *testing.T) {
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

func TestParseSimpleObject(t *testing.T) {
	tokens := []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: COLON},
		{tokenType: STRING, value: []byte("value1")},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	tokens = []token{
		{tokenType: OBJECT_START},
		{tokenType: COLON},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err == nil {
		t.Errorf("Expected error, got nil")
	}

	tokens = []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err == nil {
		t.Errorf("Expected error, got nil")
	}

	tokens = []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: COLON},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseObjectMultiStringValues(t *testing.T) {
	tokens := []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: COLON},
		{tokenType: STRING, value: []byte("value1")},
		{tokenType: COMMA},
		{tokenType: STRING, value: []byte("key2")},
		{tokenType: COLON},
		{tokenType: STRING, value: []byte("value2")},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	tokens = []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: COLON},
		{tokenType: STRING, value: []byte("value1")},
		{tokenType: COMMA},
		{tokenType: COMMA},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err == nil {
		t.Errorf("Expected error, got nil")
	}

	tokens = []token{
		{tokenType: OBJECT_START},
		{tokenType: STRING, value: []byte("key1")},
		{tokenType: COLON},
		{tokenType: STRING, value: []byte("value1")},
		{tokenType: COMMA},
		{tokenType: STRING, value: []byte("key2")},
		{tokenType: OBJECT_END},
	}

	if err := parse(tokens); err == nil {
		t.Errorf("Expected error, got nil")
	}
}
