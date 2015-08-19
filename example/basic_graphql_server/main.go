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
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/tallstreet/graphql"
	"github.com/tallstreet/graphql/executor"
	"github.com/tallstreet/graphql/executor/resolver"
	"github.com/tallstreet/graphql/handler"
	"github.com/tallstreet/graphql/schema"
	"golang.org/x/net/context"
)

var listenAddr = flag.String("l", ":8080", "listen addr")

type nowProvider struct {
	start time.Time
}

func (n *nowProvider) now(ctx context.Context, r resolver.Resolver, f *graphql.Field) (interface{}, error) {
	return time.Now(), nil
}

func (n *nowProvider) uptime(ctx context.Context, r resolver.Resolver, f *graphql.Field) (interface{}, error) {
	return time.Now().Sub(n.start).Seconds(), nil
}

func (n *nowProvider) GraphQLTypeInfo() schema.GraphQLTypeInfo {
	return schema.GraphQLTypeInfo{
		Name:        "nowProvider",
		Description: "example root field provider",
		Fields: map[string]*schema.GraphQLFieldSpec{
			"now":    {"now", "Provides the current server time", n.now, []graphql.Argument{}, true},
			"uptime": {"uptime", "Provides the current server uptime", n.uptime, []graphql.Argument{}, true},
		},
	}
}

func main() {
	flag.Parse()
	// create a new schema (which self-registers)
	now := &nowProvider{time.Now()}

	schema := schema.New()
	schema.Register(now)

	executor := executor.New(schema)
	handler := handler.New(executor)
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	log.Fatalln(http.ListenAndServe(*listenAddr, mux))
}
