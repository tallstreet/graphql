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
