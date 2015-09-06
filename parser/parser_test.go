// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "github.com/tallstreet/graphql/parser"

import (
	"log"
	"os"
	"testing"

	"sevki.org/lib/prettyprint"
	"github.com/tallstreet/graphql"
)

func TestKitchenSink(t *testing.T) {
	t.Parallel()
	var doc graphql.Document
	
	ks, _ := os.Open("../tests/relay-todo.graphql")
	if err := New("kitchenSink", ks).Decode(&doc); err != nil {
		t.Error(err.Error())

		if err != nil {
			t.Error(err)
		}

	} else {
		//		fs, _ := os.Open("../tests/complex-as-possible.graphql")
		//		e, _ := ioutil.ReadAll(fs)
		//		log.Printf(string(e))
		log.Printf(prettyprint.AsJSON(doc))
	}
	
}
