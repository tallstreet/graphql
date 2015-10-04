// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "github.com/tallstreet/graphql/parser"

import (
	"github.com/tallstreet/graphql"
)

// Decode decodes a graphql ast.
func (p *Parser) Decode(i interface{}) (err error) {
	p.Document = (i.(*graphql.Document))
	p.run()
	if p.Error != nil {
		return p.Error
	}

	return nil
}

func inlineFragmentsInSelection(doc *graphql.Document, j *graphql.SelectionSet) (err error) {
	for i := range *j {
		s := (*j)[i]
		if s.FragmentSpread != nil {
			frag := doc.LookupFragmentByName(s.FragmentSpread.Name)
			if frag != nil {
				s.InlineFragment = &graphql.InlineFragment{
					frag.TypeCondition,
					frag.Directives,
					frag.SelectionSet,
				}
				inlineFragmentsInSelection(doc, &frag.SelectionSet)
			}
		}
		if s.Field != nil {
			inlineFragmentsInSelection(doc, &s.Field.SelectionSet)
		}
	}
	return nil
}

// Goes through a graphql AST and replace fragment spreads with the fragment definitions
func InlineFragments(i interface{}) (err error) {
	doc := (i.(*graphql.Document))
	for o := range doc.Operations {
		inlineFragmentsInSelection(doc, &doc.Operations[o].SelectionSet)
	}
	return nil
}
