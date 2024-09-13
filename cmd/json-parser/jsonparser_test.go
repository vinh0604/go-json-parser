package main

import (
	"io"
	"slices"
	"strings"
	"testing"
)

func TestGetNextTokenEmptyObject(t *testing.T) {
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

func TestGetNextTokenEmptyArray(t *testing.T) {
	lexer := newLexer(strings.NewReader(`[]`))

	token, err := lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != ARRAY_START {
		t.Errorf("Expected ARRAY_START, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != ARRAY_END {
		t.Errorf("Expected ARRAY_END, got %s", token.tokenType)
	}
}

func TestGetNextTokenObjectStringValue(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{"key1":"value1"}`))

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
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("key1")) {
		t.Errorf("Expected key1, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != COLON {
		t.Errorf("Expected COLON, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("value1")) {
		t.Errorf("Expected value1, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != OBJECT_END {
		t.Errorf("Expected OBJECT_END, got %s", token.tokenType)
	}
}

func TestGetNextTokenWithSpaceChars(t *testing.T) {
	lexer := newLexer(strings.NewReader("{\t} \r\n"))

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

	token, err = lexer.getNextToken()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %s", err)
	}
}

func TestGetNextTokenObjectMultiStringValue(t *testing.T) {
	lexer := newLexer(strings.NewReader(`{"key1":"value1",
	"key2":"value2"}`))

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
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("key1")) {
		t.Errorf("Expected key1, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != COLON {
		t.Errorf("Expected COLON, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("value1")) {
		t.Errorf("Expected value1, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != COMMA {
		t.Errorf("Expected COMMA, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("key2")) {
		t.Errorf("Expected key2, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != COLON {
		t.Errorf("Expected COLON, got %s", token.tokenType)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != STRING {
		t.Errorf("Expected STRING, got %s", token.tokenType)
	}
	if !slices.Equal(token.value, []byte("value2")) {
		t.Errorf("Expected value2, got %s", token.value)
	}

	token, err = lexer.getNextToken()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if token.tokenType != OBJECT_END {
		t.Errorf("Expected OBJECT_END, got %s", token.tokenType)
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
