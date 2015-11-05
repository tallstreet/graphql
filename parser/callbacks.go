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

struct GraphQLAstIntValue;
int process_visit_int_value_cgo(struct GraphQLAstIntValue *node, void *parser){
  return processVisitIntValue(node, parser);
}

void process_end_visit_int_value_cgo(struct GraphQLAstIntValue *node, void *parser){
  processEndVisitIntValue(node, parser);
}

struct GraphQLAstFloatValue;
int process_visit_float_value_cgo(struct GraphQLAstFloatValue *node, void *parser){
  return processVisitFloatValue(node, parser);
}

void process_end_visit_float_value_cgo(struct GraphQLAstFloatValue *node, void *parser){
  processEndVisitFloatValue(node, parser);
}

struct GraphQLAstStringValue;
int process_visit_string_value_cgo(struct GraphQLAstStringValue *node, void *parser){
  return processVisitStringValue(node, parser);
}

void process_end_visit_string_value_cgo(struct GraphQLAstStringValue *node, void *parser){
  processEndVisitStringValue(node, parser);
}

struct GraphQLAstBooleanValue;
int process_visit_boolean_value_cgo(struct GraphQLAstBooleanValue *node, void *parser){
  return processVisitBooleanValue(node, parser);
}

void process_end_visit_boolean_value_cgo(struct GraphQLAstBooleanValue *node, void *parser){
  processEndVisitBooleanValue(node, parser);
}

struct GraphQLAstEnumValue;
int process_visit_enum_value_cgo(struct GraphQLAstEnumValue *node, void *parser){
  return processVisitEnumValue(node, parser);
}

void process_end_visit_enum_value_cgo(struct GraphQLAstEnumValue *node, void *parser){
  processEndVisitEnumValue(node, parser);
}

struct GraphQLAstArrayValue;
int process_visit_array_value_cgo(struct GraphQLAstArrayValue *node, void *parser){
  return processVisitArrayValue(node, parser);
}

void process_end_visit_array_value_cgo(struct GraphQLAstArrayValue *node, void *parser){
  processEndVisitArrayValue(node, parser);
}

struct GraphQLAstObjectValue;
int process_visit_object_value_cgo(struct GraphQLAstObjectValue *node, void *parser){
  return processVisitObjectValue(node, parser);
}

void process_end_visit_object_value_cgo(struct GraphQLAstObjectValue *node, void *parser){
  processEndVisitObjectValue(node, parser);
}

struct GraphQLAstObjectField;
int process_visit_object_field_cgo(struct GraphQLAstObjectField *node, void *parser){
  return processVisitObjectField(node, parser);
}

void process_end_visit_object_field_cgo(struct GraphQLAstObjectField *node, void *parser){
  processEndVisitObjectField(node, parser);
}

struct GraphQLAstDirective;
int process_visit_directive_cgo(struct GraphQLAstDirective *node, void *parser){
  return processVisitDirective(node, parser);
}

void process_end_visit_directive_cgo(struct GraphQLAstDirective *node, void *parser){
  processEndVisitDirective(node, parser);
}

struct GraphQLAstNamedType;
int process_visit_named_type_cgo(struct GraphQLAstNamedType *node, void *parser){
  return processVisitNamedType(node, parser);
}

void process_end_visit_named_type_cgo(struct GraphQLAstNamedType *node, void *parser){
  processEndVisitNamedType(node, parser);
}

struct GraphQLAstListType;
int process_visit_list_type_cgo(struct GraphQLAstListType *node, void *parser){
  return processVisitListType(node, parser);
}

void process_end_visit_list_type_cgo(struct GraphQLAstListType *node, void *parser){
  processEndVisitListType(node, parser);
}

struct GraphQLAstNonNullType;
int process_visit_non_null_type_cgo(struct GraphQLAstNonNullType *node, void *parser){
  return processVisitNonNullType(node, parser);
}

void process_end_visit_non_null_type_cgo(struct GraphQLAstNonNullType *node, void *parser){
  processEndVisitNonNullType(node, parser);
}

struct GraphQLAstName;
int process_visit_name_cgo(struct GraphQLAstName *node, void *parser){
  return processVisitName(node, parser);
}

void process_end_visit_name_cgo(struct GraphQLAstName *node, void *parser){
  processEndVisitName(node, parser);
}

*/
import "C"
