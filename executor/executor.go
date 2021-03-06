package executor

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/tallstreet/graphql"
	"github.com/tallstreet/graphql/schema"
	"golang.org/x/net/context"
)

type Executor struct {
	schema *schema.Schema
}

func New(schema *schema.Schema) *Executor {
	return &Executor{
		schema: schema,
	}
}

func (e *Executor) HandleOperation(ctx context.Context, o *graphql.Operation) (interface{}, error) {
	rootSelections := o.SelectionSet
	rootFields := e.schema.RootFields()
	result := make(map[string]interface{})

	for _, selection := range rootSelections {
		rootFieldHandler, ok := rootFields[selection.Field.Name]
		if !ok {
			return nil, fmt.Errorf("Root field '%s' is not registered", selection.Field.Name)
		}
		partial, err := rootFieldHandler.Func(ctx, e, selection.Field)
		if err != nil {
			return nil, err
		}
		resolved, err := e.Resolve(ctx, partial, selection.Field)
		if err != nil {
			return nil, err
		}
		result[selection.Field.Alias] = resolved
	}
	return result, nil
}

func isSlice(value interface{}) bool {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

type fieldResult struct {
	FieldName string
	Value     interface{}
	Err       error
	Index     int
}

func (e *Executor) findFields(selectionSet graphql.SelectionSet) []*graphql.Field {

	fields := []*graphql.Field{}
	for _, selection := range selectionSet {
		if selection.InlineFragment != nil {
			childFields := e.findFields(selection.InlineFragment.SelectionSet)
			fields = append(fields, childFields...)
		}

		if selection.Field != nil {
			fields = append(fields, selection.Field)
		}
	}
	return fields
}

func (e *Executor) mergeValues(oldValue interface{}, newValue interface{}) interface{} {
	if isSlice(oldValue) {
		oldValueSlice := oldValue.([]interface{})
		newValueSlice := newValue.([]interface{})
		newSlice := make([]interface{}, 0, len(oldValueSlice))
		for k := range oldValueSlice {
			newSlice = append(newSlice, e.mergeValues(oldValueSlice[k], newValueSlice[k]))
		}
		return newSlice
	}
	if reflect.TypeOf(oldValue).Kind() == reflect.Map {
		oldValueMap := oldValue.(map[string]interface{})
		newValueMap := newValue.(map[string]interface{})
		newMap := newValueMap
		for k := range oldValueMap {
			if oldValueMap[k] != nil && newValueMap[k] != nil {
				newMap[k] = e.mergeValues(oldValueMap[k], newValueMap[k])
			} else {
				newMap[k] = oldValueMap[k]
			}
		}
		return newMap

	}
	if oldValue != nil {
		return oldValue
	}
	if newValue != nil {
		return newValue
	}
	return nil
}

func (e *Executor) Resolve(ctx context.Context, partial interface{}, field *graphql.Field) (interface{}, error) {
	if partial != nil && isSlice(partial) {
		return e.resolveSlice(ctx, partial, field)
	}
	graphQLValue, ok := partial.(schema.GraphQLType)
	// if we have a scalar we're done
	if !ok {
		//log.Printf("returning scalar %T: %v\n", partial, partial)
		return partial, nil
	}
	// check against returning object as non-leaf
	if len(field.SelectionSet) == 0 {
		return nil, fmt.Errorf("Cannot return a '%T' as a leaf", graphQLValue)
	}

	fields := e.findFields(field.SelectionSet)

	result := map[string]interface{}{}
	typeInfo := schema.WithIntrospectionField(graphQLValue.GraphQLTypeInfo())
	results := make(chan fieldResult)
	wg := sync.WaitGroup{}

	for _, selection := range fields {
		fieldName := selection.Name
		fieldHandler, ok := typeInfo.Fields[fieldName]
		if !ok {
			return nil, fmt.Errorf("No handler for field '%s' on type '%T'", fieldName, graphQLValue)
		}
		wg.Add(1)
		go func(selection *graphql.Field) {
			defer wg.Done()
			partial, err := fieldHandler.Func(ctx, e, selection)
			if err != nil {
				results <- fieldResult{Err: err}
				return
			}
			resolved, err := e.Resolve(ctx, partial, selection)
			if err != nil {
				results <- fieldResult{Err: err}
				return
			}
			if selection.Alias != "" {
				fieldName = selection.Alias
			}
			results <- fieldResult{
				FieldName: fieldName, Value: resolved, Err: err,
			}
		}(selection)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		if r.Err != nil {
			return nil, r.Err
		}
		if result[r.FieldName] != nil {
			result[r.FieldName] = e.mergeValues(result[r.FieldName], r.Value)
		} else {
			result[r.FieldName] = r.Value
		}
	}
	return result, nil
}

func (e *Executor) resolveSlice(ctx context.Context, partials interface{}, field *graphql.Field) (interface{}, error) {
	v := reflect.ValueOf(partials)
	results := make([]interface{}, v.Len(), v.Len())
	resChan := make(chan fieldResult)
	wg := sync.WaitGroup{}
	for i := 0; i < v.Len(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result, err := e.Resolve(ctx, v.Index(i).Interface(), field)
			resChan <- fieldResult{Value: result, Err: err, Index: i}
		}(i)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()
	for result := range resChan {
		if result.Err != nil {
			return nil, result.Err
		}
		results[result.Index] = result.Value
	}
	return results, nil
}
