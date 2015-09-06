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
int process_visit_variable_definition_cgo(struct GraphQLAstVariableDefinition *node, void *parser);
void process_end_visit_variable_definition_cgo(struct GraphQLAstVariableDefinition *node, void *parser);
int process_visit_selection_set_cgo(struct GraphQLAstSelectionSet *node, void *parser);
void process_end_visit_selection_set_cgo(struct GraphQLAstSelectionSet *node, void *parser);
int process_visit_field_cgo(struct GraphQLAstField *node, void *parser);
void process_end_visit_field_cgo(struct GraphQLAstField *node, void *parser);
int process_visit_argument_cgo(struct GraphQLAstArgument *node, void *parser);
void process_end_visit_argument_cgo(struct GraphQLAstArgument *node, void *parser);
int process_visit_fragment_spread_cgo(struct GraphQLAstFragmentSpread *node, void *parser);
void process_end_visit_fragment_spread_cgo(struct GraphQLAstFragmentSpread *node, void *parser);
int process_visit_inline_fragment_cgo(struct GraphQLAstInlineFragment *node, void *parser);
void process_end_visit_inline_fragment_cgo(struct GraphQLAstInlineFragment *node, void *parser);
int process_visit_fragment_definition_cgo(struct GraphQLAstFragmentDefinition *node, void *parser);
void process_end_visit_fragment_definition_cgo(struct GraphQLAstFragmentDefinition *node, void *parser);
int process_visit_variable_cgo(struct GraphQLAstVariable *node, void *parser);
void process_end_visit_variable_cgo(struct GraphQLAstVariable *node, void *parser);
int process_visit_int_value_cgo(struct GraphQLAstIntValue *node, void *parser);
void process_end_visit_int_value_cgo(struct GraphQLAstIntValue *node, void *parser);
int process_visit_float_value_cgo(struct GraphQLAstFloatValue *node, void *parser);
void process_end_visit_float_value_cgo(struct GraphQLAstFloatValue *node, void *parser);
int process_visit_string_value_cgo(struct GraphQLAstStringValue *node, void *parser);
void process_end_visit_string_value_cgo(struct GraphQLAstStringValue *node, void *parser);
int process_visit_boolean_value_cgo(struct GraphQLAstBooleanValue *node, void *parser);
void process_end_visit_boolean_value_cgo(struct GraphQLAstBooleanValue *node, void *parser);
int process_visit_enum_value_cgo(struct GraphQLAstEnumValue *node, void *parser);
void process_end_visit_enum_value_cgo(struct GraphQLAstEnumValue *node, void *parser);
int process_visit_array_value_cgo(struct GraphQLAstArrayValue *node, void *parser);
void process_end_visit_array_value_cgo(struct GraphQLAstArrayValue *node, void *parser);
int process_visit_object_value_cgo(struct GraphQLAstObjectValue *node, void *parser);
void process_end_visit_object_value_cgo(struct GraphQLAstObjectValue *node, void *parser);
int process_visit_object_field_cgo(struct GraphQLAstObjectField *node, void *parser);
void process_end_visit_object_field_cgo(struct GraphQLAstObjectField *node, void *parser);
int process_visit_directive_cgo(struct GraphQLAstDirective *node, void *parser);
void process_end_visit_directive_cgo(struct GraphQLAstDirective *node, void *parser);
int process_visit_named_type_cgo(struct GraphQLAstNamedType *node, void *parser);
void process_end_visit_named_type_cgo(struct GraphQLAstNamedType *node, void *parser);
int process_visit_list_type_cgo(struct GraphQLAstListType *node, void *parser);
void process_end_visit_list_type_cgo(struct GraphQLAstListType *node, void *parser);
int process_visit_non_null_type_cgo(struct GraphQLAstNonNullType *node, void *parser);
void process_end_visit_non_null_type_cgo(struct GraphQLAstNonNullType *node, void *parser);
int process_visit_name_cgo(struct GraphQLAstName *node, void *parser);
void process_end_visit_name_cgo(struct GraphQLAstName *node, void *parser);
*/
import "C"

import (
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
	"errors"
	"strconv"
	
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


//export processVisitVariableDefinition
func processVisitVariableDefinition(node *C.struct_GraphQLAstVariableDefinition, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	variable := &graphql.VariableDefinition {
	}
	operation := p.nodes.Head().(*graphql.Operation)
	operation.VariableDefinitions = append(operation.VariableDefinitions, variable)
	p.visitNode(variable)
	return 1
}

//export processEndVisitVariableDefinition
func processEndVisitVariableDefinition(node *C.struct_GraphQLAstVariableDefinition, parser unsafe.Pointer) {
	var variable *graphql.VariableDefinition
	var ok bool
	p := (*Parser)(parser)
	last1 := p.endVisitNode()
	last2 := p.endVisitNode()
	last3 := p.endVisitNode()
	value, ok := last1.(*graphql.Value)
	if ok {
		variable, ok = last3.(*graphql.VariableDefinition)
		variable.Variable = last2.(*graphql.Variable)
		variable.DefaultValue = value
	} else {
		p.visitNode(last3)
		variable, ok = last2.(*graphql.VariableDefinition)
		variable.Variable = last1.(*graphql.Variable)
	}	
	typeT := (*C.struct_GraphQLAstNamedType)(C.GraphQLAstVariableDefinition_get_type(node))
	typeName := C.GraphQLAstNamedType_get_name(typeT)
	variable.Type.Name = C.GoString(C.GraphQLAstName_get_value(typeName))
	
}


//export processVisitSelectionSet
func processVisitSelectionSet(node *C.struct_GraphQLAstSelectionSet, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitSelectionSet
func processEndVisitSelectionSet(node *C.struct_GraphQLAstSelectionSet, parser unsafe.Pointer) {
}

//export processVisitField
func processVisitField(node *C.struct_GraphQLAstField, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	selection := &graphql.Selection{
		Field: &graphql.Field{},
	}
	operation, ok := p.nodes.Head().(*graphql.Operation)
	if ok {
	  operation.SelectionSet = append(operation.SelectionSet, selection)
	}
	p.visitNode(selection)
	return 1
}

//export processEndVisitField
func processEndVisitField(node *C.struct_GraphQLAstField, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	field := p.nodes.Head().(*graphql.Selection)
	if field.Field != nil {
		alias := C.GraphQLAstField_get_alias(node)
		if alias != nil {
			field.Field.Alias = C.GoString(C.GraphQLAstName_get_value(alias))
		}
		name := C.GraphQLAstField_get_name(node)
		if name != nil {
			field.Field.Name = C.GoString(C.GraphQLAstName_get_value(name))
		}
		p.endVisitNode()
	}
}

//export processVisitArgument
func processVisitArgument(node *C.struct_GraphQLAstArgument, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitArgument
func processEndVisitArgument(node *C.struct_GraphQLAstArgument, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	name := C.GraphQLAstArgument_get_name(node)
	value := p.endVisitNode()
	argument := graphql.Argument {
		Name: C.GoString(C.GraphQLAstName_get_value(name)), 
		Value: value,
	}
	selection, ok := p.nodes.Head().(*graphql.Selection)
	if ok {
		selection.Field.Arguments = append(selection.Field.Arguments, argument)
	}
	directive, ok := p.nodes.Head().(*graphql.Directive)
	if ok {
		directive.Arguments = append(directive.Arguments, argument)
	}
}

//export processVisitFragmentSpread
func processVisitFragmentSpread(node *C.struct_GraphQLAstFragmentSpread, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	name := C.GraphQLAstFragmentSpread_get_name(node)
	selection := &graphql.Selection{
		FragmentSpread: &graphql.FragmentSpread{
			Name: C.GoString(C.GraphQLAstName_get_value(name)),
		},
	}
	operation, ok := p.nodes.Head().(*graphql.Operation)
	if ok {
	  operation.SelectionSet = append(operation.SelectionSet, selection)
	}
	return 0
}

//export processEndVisitFragmentSpread
func processEndVisitFragmentSpread(node *C.struct_GraphQLAstFragmentSpread, parser unsafe.Pointer) {
}

//export processVisitInlineFragment
func processVisitInlineFragment(node *C.struct_GraphQLAstInlineFragment, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	selection := &graphql.Selection{
		InlineFragment: &graphql.InlineFragment{},
	}
	operation, ok := p.endVisitNode().(*graphql.Operation)
	if ok {
	  operation.SelectionSet = append(operation.SelectionSet, selection)
	}
	p.visitNode(selection)
	return 1
}

//export processEndVisitInlineFragment
func processEndVisitInlineFragment(node *C.struct_GraphQLAstInlineFragment, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	
	fragment := p.nodes.Head().(*graphql.Selection).InlineFragment
	condition := C.GraphQLAstInlineFragment_get_type_condition(node)
	condition_name := C.GraphQLAstNamedType_get_name(condition)
	fragment.TypeCondition = C.GoString(C.GraphQLAstName_get_value(condition_name))
	p.endVisitNode()
}

//export processVisitFragmentDefinition
func processVisitFragmentDefinition(node *C.struct_GraphQLAstFragmentDefinition, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	doc := p.nodes.Head().(*graphql.Document)
	name := C.GraphQLAstFragmentDefinition_get_name(node)
	fragment := &graphql.FragmentDefinition {
		Name: C.GoString(C.GraphQLAstName_get_value(name)),
	}
	doc.FragmentDefinitions = append(doc.FragmentDefinitions, fragment)
	p.visitNode(fragment)
	return 1
}

//export processEndVisitFragmentDefinition
func processEndVisitFragmentDefinition(node *C.struct_GraphQLAstFragmentDefinition, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	fragment := p.nodes.Head().(*graphql.FragmentDefinition)
	condition := C.GraphQLAstFragmentDefinition_get_type_condition(node)
	condition_name := C.GraphQLAstNamedType_get_name(condition)
	fragment.TypeCondition = C.GoString(C.GraphQLAstName_get_value(condition_name))
	p.endVisitNode()
}

//export processVisitVariable
func processVisitVariable(node *C.struct_GraphQLAstVariable, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitVariable
func processEndVisitVariable(node *C.struct_GraphQLAstVariable, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	variable := &graphql.Variable {}
	name := C.GraphQLAstVariable_get_name(node)
	if name != nil {
		variable.Name = C.GoString(C.GraphQLAstName_get_value(name))
	}
	p.visitNode(variable)
}

//export processVisitIntValue
func processVisitIntValue(node *C.struct_GraphQLAstIntValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitIntValue
func processEndVisitIntValue(node *C.struct_GraphQLAstIntValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	value := C.GoString(C.GraphQLAstIntValue_get_value(node))
	i, _ := strconv.ParseInt(value, 10, 64)
	v := &graphql.Value{ Value: i }
	p.visitNode(v)
}

//export processVisitFloatValue
func processVisitFloatValue(node *C.struct_GraphQLAstFloatValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitFloatValue
func processEndVisitFloatValue(node *C.struct_GraphQLAstFloatValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	value := C.GoString(C.GraphQLAstFloatValue_get_value(node))
	f, _ := strconv.ParseFloat(value, 64)
	v := &graphql.Value{ Value: f }
	p.visitNode(v)
}

//export processVisitStringValue
func processVisitStringValue(node *C.struct_GraphQLAstStringValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitStringValue
func processEndVisitStringValue(node *C.struct_GraphQLAstStringValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	value := C.GraphQLAstStringValue_get_value(node)
	v := &graphql.Value{ Value: C.GoString(value)}
	p.visitNode(v)
}

//export processVisitBooleanValue
func processVisitBooleanValue(node *C.struct_GraphQLAstBooleanValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitBooleanValue
func processEndVisitBooleanValue(node *C.struct_GraphQLAstBooleanValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	value := C.GraphQLAstBooleanValue_get_value(node)
	v := &graphql.Value{ Value: value == 1 }
	p.visitNode(v)
}

//export processVisitEnumValue
func processVisitEnumValue(node *C.struct_GraphQLAstEnumValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitEnumValue
func processEndVisitEnumValue(node *C.struct_GraphQLAstEnumValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	value := C.GraphQLAstEnumValue_get_value(node)
	v := &graphql.Value{ Value: C.GoString(value)}
	p.visitNode(v)
}


//export processVisitArrayValue
func processVisitArrayValue(node *C.struct_GraphQLAstArrayValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitArrayValue
func processEndVisitArrayValue(node *C.struct_GraphQLAstArrayValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	size := int(C.GraphQLAstArrayValue_get_values_size(node))
	array := make([]interface{}, size, size)
	for i := size - 1; i > 0; i-- {
		array[i] = p.endVisitNode()
	}
	v := &graphql.Value{ Value: array }
	p.visitNode(v)
}

//export processVisitObjectValue
func processVisitObjectValue(node *C.struct_GraphQLAstObjectValue, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitObjectValue
func processEndVisitObjectValue(node *C.struct_GraphQLAstObjectValue, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	size := int(C.GraphQLAstObjectValue_get_fields_size(node))
	object := make(map[string]interface{}, size)
	for i := 0; i < size; i++ {
		val := p.endVisitNode().(map[string]interface{})
		for k, v := range val {
		    object[k] = v
		}
	}
	v := &graphql.Value{ Value: object }
	p.visitNode(v)
}

//export processVisitObjectField
func processVisitObjectField(node *C.struct_GraphQLAstObjectField, parser unsafe.Pointer) int {
	return 1
}

//export processEndVisitObjectField
func processEndVisitObjectField(node *C.struct_GraphQLAstObjectField, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	field := make(map[string]interface{}, 1)
	name := C.GraphQLAstObjectField_get_name(node)
	value := p.endVisitNode()
	field[C.GoString(C.GraphQLAstName_get_value(name))] = value
	p.visitNode(field)
}

//export processVisitDirective
func processVisitDirective(node *C.struct_GraphQLAstDirective, parser unsafe.Pointer) int {
	p := (*Parser)(parser)
	directive := &graphql.Directive{}
	p.visitNode(directive)
	return 1
}

//export processEndVisitDirective
func processEndVisitDirective(node *C.struct_GraphQLAstDirective, parser unsafe.Pointer) {
	p := (*Parser)(parser)
	directive := p.nodes.Head().(*graphql.Directive)
	name := C.GraphQLAstDirective_get_name(node)
	if name != nil {
		directive.Name = C.GoString(C.GraphQLAstName_get_value(name))
	}
	p.endVisitNode()
}

//export processVisitNamedType
func processVisitNamedType(node *C.struct_GraphQLAstNamedType, parser unsafe.Pointer) int {
	return 0
}

//export processEndVisitNamedType
func processEndVisitNamedType(node *C.struct_GraphQLAstNamedType, parser unsafe.Pointer) {
}

//export processVisitListType
func processVisitListType(node *C.struct_GraphQLAstListType, parser unsafe.Pointer) int {
	return 0
}

//export processEndVisitListType
func processEndVisitListType(node *C.struct_GraphQLAstListType, parser unsafe.Pointer) {
}

//export processVisitNonNullType
func processVisitNonNullType(node *C.struct_GraphQLAstNonNullType, parser unsafe.Pointer) int {
	return 0
}

//export processEndVisitNonNullType
func processEndVisitNonNullType(node *C.struct_GraphQLAstNonNullType, parser unsafe.Pointer) {
}

//export processVisitName
func processVisitName(node *C.struct_GraphQLAstName, parser unsafe.Pointer) int {
	return 0
}

//export processEndVisitName
func processEndVisitName(node *C.struct_GraphQLAstName, parser unsafe.Pointer) {
}

func (p *Parser) visitNode(node interface{}) {
	p.nodes.Push(node)
}

func (p *Parser) endVisitNode() interface{} {
	node := p.nodes.Pop()
	return node
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
	}
	C.graphql_node_visit(ast, &visitor_callbacks, unsafe.Pointer(p))

	C.graphql_node_free(ast)
}
