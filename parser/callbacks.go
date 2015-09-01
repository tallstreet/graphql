/**
 *  Copyright (c) 2015, Facebook, Inc.
 *  All rights reserved.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree. An additional grant
 *  of patent rights can be found in the PATENTS file in the same directory.
 */
package parser

/*
struct GraphQLAstDocument;
int process_visit_document_cgo(struct GraphQLAstDocument *node, void *parser) {
  return processVisitDocument(node, parser);
}

void process_end_visit_document_cgo(struct GraphQLAstDocument *node, void *parser) {
  processEndVisitDocument(node, parser);
}

struct GraphQLAstOperationDefinition;
int process_visit_operation_definition_cgo(struct GraphQLAstOperationDefinition *node, void *parser) {
  return processVisitOperationDefinition(node, parser);
}

void process_end_visit_operation_definition_cgo(struct GraphQLAstOperationDefinition *node, void *parser) {
  processEndVisitOperationDefinition(node, parser);
}

struct GraphQLAstVariableDefinition;
int process_visit_variable_definition_cgo(struct GraphQLAstVariableDefinition *node, void *parser) {
  return processVisitVariableDefinition(node, parser);
}

void process_end_visit_variable_definition_cgo(struct GraphQLAstVariableDefinition *node, void *parser) {
  processEndVisitVariableDefinition(node, parser);
}

struct GraphQLAstSelectionSet;
int process_visit_selection_set_cgo(struct GraphQLAstSelectionSet *node, void *parser) {
  return processVisitSelectionSet(node, parser);
}

void process_end_visit_selection_set_cgo(struct GraphQLAstSelectionSet *node, void *parser) {
  processEndVisitSelectionSet(node, parser);
}

struct GraphQLAstField;
int process_visit_field_cgo(struct GraphQLAstField *node, void *parser) {
  return processVisitField(node, parser);
}

void process_end_visit_field_cgo(struct GraphQLAstField *node, void *parser) {
  processEndVisitField(node, parser);
}

struct GraphQLAstArgument;
int process_visit_argument_cgo(struct GraphQLAstArgument *node, void *parser) {
  return processVisitArgument(node, parser);
}

void process_end_visit_argument_cgo(struct GraphQLAstArgument *node, void *parser) {
  processEndVisitArgument(node, parser);
}

struct GraphQLAstFragmentSpread;
int process_visit_fragment_spread_cgo(struct GraphQLAstFragmentSpread *node, void *parser) {
  return processVisitFragmentSpread(node, parser);
}

void process_end_visit_fragment_spread_cgo(struct GraphQLAstFragmentSpread *node, void *parser) {
  processEndVisitFragmentSpread(node, parser);
}

struct GraphQLAstInlineFragment;
int process_visit_inline_fragment_cgo(struct GraphQLAstInlineFragment *node, void *parser) {
  return processVisitInlineFragment(node, parser);
}

void process_end_visit_inline_fragment_cgo(struct GraphQLAstInlineFragment *node, void *parser) {
  processEndVisitInlineFragment(node, parser);
}

struct GraphQLAstFragmentDefinition;
int process_visit_fragment_definition_cgo(struct GraphQLAstFragmentDefinition *node, void *parser) {
  return processVisitFragmentDefinition(node, parser);
}

void process_end_visit_fragment_definition_cgo(struct GraphQLAstFragmentDefinition *node, void *parser) {
  processEndVisitFragmentDefinition(node, parser);
}

struct GraphQLAstVariable;
int process_visit_variable_cgo(struct GraphQLAstVariable *node, void *parser){
  return processVisitVariable(node, parser);
}

void process_end_visit_variable_cgo(struct GraphQLAstVariable *node, void *parser){
  processEndVisitVariable(node, parser);
}


*/
import "C"
