package parser_test

import (
	"testing"

	"github.com/tmc/graphql/parser"
)

var shouldParse = []string{
	`foo()`,
	`node(42)`,
	`foo(){id}`,
	`foo(1,2){id}`,
	`foo(1,2){id,{nest,some,{fields}}}`,
}

func TestSuccessfulParses(t *testing.T) {
	for i, in := range shouldParse {
		_, err := parser.Parse("parser_test.go", []byte(in))
		if err != nil {
			t.Errorf("case %d: %v", i+1, err)
		}
	}
}