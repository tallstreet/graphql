package schema

import (
	"fmt"

	"github.com/tmc/graphql"
)

type Schema struct {
	registeredTypes []GraphQLType
	rootCalls       map[string]CallHandler
}

type CallHandler func(graphql.Call) (interface{}, error)

type GraphQLType interface {
	RootCalls() map[string]CallHandler
}

func New() *Schema {
	s := &Schema{
		registeredTypes: []GraphQLType{},
		rootCalls:       map[string]CallHandler{},
	}
	// self-register
	s.Register(s)
	return s
}

func (s *Schema) HandleCall(c graphql.Call) (interface{}, error) {
	handler, ok := s.rootCalls[c.Name]
	if !ok {
		return nil, fmt.Errorf("schema: no registered types handle the root call '%s'", c.Name)
	}
	return handler(c)
}

func (s *Schema) Register(t GraphQLType) {
	s.registeredTypes = append(s.registeredTypes, t)
	for call, handler := range t.RootCalls() {
		// TODO(tmc): collision handling
		s.rootCalls[call] = handler
	}
}

func (s *Schema) RootCalls() map[string]CallHandler {
	return map[string]CallHandler{
		"schema": s.handleSchemaCall,
	}
}

func (s *Schema) handleSchemaCall(c graphql.Call) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
