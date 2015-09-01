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
  return processVisitOperationDefinition(node, parser);
}

void process_end_visit_variable_definition_cgo(struct GraphQLAstVariableDefinition *node, void *parser) {
  processEndVisitOperationDefinition(node, parser);
}
*/
import "C"
