// Copyright 2015 Gary Roberts <contact@tallstreet.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "github.com/tallstreet/graphql/parser"

/*
#cgo CFLAGS: -I ../libgraphqlparser/c/
#cgo LDFLAGS: -L ../libgraphqlparser -lgraphqlparser
#include "GraphQLAst.h"
#include "GraphQLAstNode.h"
#include "GraphQLAstVisitor.h"
#include "GraphQLParser.h"
#include <stdlib.h>

int process_visit_document_cgo(struct GraphQLAstDocument *node, void *parser);
void process_end_visit_document_cgo(struct GraphQLAstDocument *node, void *parser);
int process_visit_operation_definition_cgo(struct GraphQLAstOperationDefinition *node, void *parser);
void process_end_visit_operation_definition_cgo(struct GraphQLAstOperationDefinition *node, void *parser);

*/
import "C"

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"unsafe"
	"errors"
	
	"github.com/tallstreet/graphql"
	"github.com/oleiade/lane"
)

type Parser struct {
	name     string
	query    string
	Error    error
	Document *graphql.Document
	nodes    *lane.Stack
}

type stateFn func(*Parser) stateFn

//export processVisitDocument
func processVisitDocument(node *C.struct_GraphQLAstDocument, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	p.visitNode(p.Document)
	return 1
}

//export processEndVisitDocument
func processEndVisitDocument(node *C.struct_GraphQLAstDocument, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	doc := p.nodes.Head().(*graphql.Document)
	doc.DefinitionSize = int(C.GraphQLAstDocument_get_definitions_size(node))
	p.endVisitNode()
}

//export processVisitOperationDefinition
func processVisitOperationDefinition(node *C.struct_GraphQLAstOperationDefinition, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	doc := p.nodes.Head().(*graphql.Document)
	operation := &graphql.Operation {}
	doc.Operations = append(doc.Operations, operation)
	p.visitNode(operation)
	return 1
}

//export processEndVisitOperationDefinition
func processEndVisitOperationDefinition(node *C.struct_GraphQLAstOperationDefinition, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	operation := p.nodes.Head().(*graphql.Operation)
	name := C.GraphQLAstOperationDefinition_get_name(node)
	if name != nil {
		operation.Name = C.GoString(C.GraphQLAstName_get_value(name))
	}
	operation.Type = graphql.OperationType(C.GoString(C.GraphQLAstOperationDefinition_get_operation(node)))
	p.endVisitNode()
}


func (p *Parser) visitNode(node interface{}) {
	p.nodes.Push(node)
}

func (p *Parser) endVisitNode() {
	p.nodes.Pop()
}

func parse(query string) (*C.struct_GraphQLAstNode, error) {
	graphql := C.CString(query)
	cError := (*C.char)(nil)
	ast := C.graphql_parse_string(graphql, &cError)
	C.free(unsafe.Pointer(graphql))

	if ast == nil {
		err := errors.New(C.GoString(cError))
		C.graphql_error_free(cError)
		return nil, err
	}
	return ast, nil
}

func New(name string, r io.Reader) *Parser {
	var doc graphql.Document

	query, _ := ioutil.ReadAll(r)
	log.Printf(string(query));
	p := &Parser{
		name:     name,
		query:    string(query),
		Document: &doc,
		nodes:    lane.NewStack(),
	}
	return p
}

// run runs the state machine for the Scanner.
func (p *Parser) run() {
	ast, err := parse(p.query)
	if err != nil {
		fmt.Printf("BUG: unexpected parse error: %s", err)
		return
	}
	visitor_callbacks := C.struct_GraphQLAstVisitorCallbacks{
		visit_document: (C.visit_document_func)(C.process_visit_document_cgo),
		end_visit_document: (C.end_visit_document_func)(C.process_end_visit_document_cgo),
		visit_operation_definition: (C.visit_operation_definition_func)(C.process_visit_operation_definition_cgo),
		end_visit_operation_definition: (C.end_visit_operation_definition_func)(C.process_end_visit_operation_definition_cgo),
		/*
		visit_variable_definition: (C.visit_variable_definition_func)(C.process_visit_variable_definition_cgo),
		end_visit_variable_definition: (C.end_visit_variable_definition_func)(C.process_end_visit_variable_definition_cgo),
		visit_selection_set: (C.visit_selection_set_func)(C.process_visit_selection_set_cgo),
		end_visit_selection_set: (C.end_visit_selection_set_func)(C.process_end_visit_selection_set_cgo),
		visit_field: (C.visit_field_func)(C.process_visit_field_cgo),
		end_visit_field: (C.end_visit_field_func)(C.process_end_visit_field_cgo),
		visit_argument: (C.visit_argument_func)(C.process_visit_argument_cgo),
		end_visit_argument: (C.end_visit_argument_func)(C.process_end_visit_argument_cgo),
		visit_fragment_spread: (C.visit_fragment_spread_func)(C.process_visit_fragment_spread_cgo),
		end_visit_fragment_spread: (C.end_visit_fragment_spread_func)(C.process_end_visit_fragment_spread_cgo),
		visit_inline_fragment: (C.visit_inline_fragment_func)(C.process_visit_inline_fragment_cgo),
		end_visit_inline_fragment: (C.end_visit_inline_fragment_func)(C.process_end_visit_inline_fragment_cgo),
		visit_fragment_definition: (C.visit_fragment_definition_func)(C.process_visit_fragment_definition_cgo),
		end_visit_fragment_definition: (C.end_visit_fragment_definition_func)(C.process_end_visit_fragment_definition_cgo),
		visit_variable: (C.visit_variable_func)(C.process_visit_variable_cgo),
		end_visit_variable: (C.end_visit_variable_func)(C.process_end_visit_variable_cgo),
		visit_int_value: (C.visit_int_value_func)(C.process_visit_int_value_cgo),
		end_visit_int_value: (C.end_visit_int_value_func)(C.process_end_visit_int_value_cgo),
		visit_float_value: (C.visit_float_value_func)(C.process_visit_float_value_cgo),
		end_visit_float_value: (C.end_visit_float_value_func)(C.process_end_visit_float_value_cgo),
		visit_string_value: (C.visit_string_value_func)(C.process_visit_string_value_cgo),
		end_visit_string_value: (C.end_visit_string_value_func)(C.process_end_visit_string_value_cgo),
		visit_boolean_value: (C.visit_boolean_value_func)(C.process_visit_boolean_value_cgo),
		end_visit_boolean_value: (C.end_visit_boolean_value_func)(C.process_end_visit_boolean_value_cgo),
		visit_enum_value: (C.visit_enum_value_func)(C.process_visit_enum_value_cgo),
		end_visit_enum_value: (C.end_visit_enum_value_func)(C.process_end_visit_enum_value_cgo),
		visit_array_value: (C.visit_array_value_func)(C.process_visit_array_value_cgo),
		end_visit_array_value: (C.end_visit_array_value_func)(C.process_end_visit_array_value_cgo),
		visit_object_value: (C.visit_object_value_func)(C.process_visit_object_value_cgo),
		end_visit_object_value: (C.end_visit_object_value_func)(C.process_end_visit_object_value_cgo),
		visit_object_field: (C.visit_object_field_func)(C.process_visit_object_field_cgo),
		end_visit_object_field: (C.end_visit_object_field_func)(C.process_end_visit_object_field_cgo),
		visit_directive: (C.visit_directive_func)(C.process_visit_directive_cgo),
		end_visit_directive: (C.end_visit_directive_func)(C.process_end_visit_directive_cgo),
		visit_named_type: (C.visit_named_type_func)(C.process_visit_named_type_cgo),
		end_visit_named_type: (C.end_visit_named_type_func)(C.process_end_visit_named_type_cgo),
		visit_list_type: (C.visit_list_type_func)(C.process_visit_list_type_cgo),
		end_visit_list_type: (C.end_visit_list_type_func)(C.process_end_visit_list_type_cgo),
		visit_non_null_type: (C.visit_non_null_type_func)(C.process_visit_non_null_type_cgo),
		end_visit_non_null_type: (C.end_visit_non_null_type_func)(C.process_end_visit_non_null_type_cgo),
		visit_name: (C.visit_name_func)(C.process_visit_name_cgo),
		end_visit_name: (C.end_visit_name_func)(C.process_end_visit_name_cgo),
		*/
	}
	C.graphql_node_visit(ast, &visitor_callbacks, unsafe.Pointer(p))

	C.graphql_node_free(ast)
}
