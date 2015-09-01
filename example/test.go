// Program basic_graphql_server shows a simple HTTP server that exposes a bare schema.
//
// Example:
//  $ go get github.com/tallstreet/graphql/example/basic_graphql_server
//  $ basic_graphql_server &
//  $ curl -g 'http://localhost:8080/?q={__schema{root_fields{name,description}}}'
//  {"data":[{"root_fields":[{"description": "Schema entry root field","name":"__schema"}]}}]
//
// Here we see the server showing the available root fields ("schema").
package main

import (
	"log"
	"os"
	"sevki.org/lib/prettyprint"

	"github.com/tallstreet/graphql/parser"
	"github.com/tallstreet/graphql"
)

func main() {
	var doc graphql.Document
	fs, _ := os.Open("./tests/kitchen-sink.graphql")
	if err := parser.New("graphql", fs).Decode(&doc); err != nil {
		
		log.Printf(err.Error())
		
	} else {
		
		log.Printf(prettyprint.AsJSON(doc))
		log.Print(doc.DefinitionSize)
	}
}
