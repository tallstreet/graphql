// Copyright 2015 Gary Roberts <contact@tallstreet.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "github.com/tallstreet/graphql/parser"

import (
	"io"
	"io/ioutil"
	"fmt"

	"github.com/oleiade/lane"
	"github.com/tallstreet/graphql"
	
	"github.com/chris-ramon/graphql-go/language/parser"
	"github.com/chris-ramon/graphql-go/language/ast"
	
)

type Parser struct {
	name     string
	query    string
	Error    error
	Document *graphql.Document
	parsedDoc *ast.Document
	nodes    *lane.Stack
}

func New(name string, r io.Reader) *Parser {
	
	
	var doc graphql.Document
	query, _ := ioutil.ReadAll(r)
	
	opts := parser.ParseOptions{
		NoSource: true,
	}
	params := parser.ParseParams{
		Source:  string(query),
		Options: opts,
	}
	parsedDoc, _ := parser.Parse(params)
	
	p := &Parser{
		name:     name,
		query:    string(query),
		Document: &doc,
		parsedDoc: parsedDoc,
		nodes:    lane.NewStack(),
	}
	return p
}

// run runs the state machine for the Scanner.
func (p *Parser) run() {
	
	fmt.Printf("%s\n", p.parsedDoc.Kind)
	for i := range p.parsedDoc.Definitions {
		a := p.parsedDoc.Definitions[i]
		fmt.Printf("%s", a.GetKind())
	}
}
