package executor

import (
	"fmt"
	//"log"
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
}


func (e *Executor) callHandler(ctx context.Context, partial interface{}, field *graphql.Field, typeInfo schema.GraphQLTypeInfo, results *chan fieldResult, wg *sync.WaitGroup, child *graphql.Field) {
	defer wg.Done()
	
	fieldName := child.Name
	fieldHandler, _ := typeInfo.Fields[fieldName]
	partial, err := fieldHandler.Func(ctx, e, child)
	if err != nil {
		*results <- fieldResult{Err: err}
		return
	}
	resolved, err := e.Resolve(ctx, partial, child)
	if err != nil {
		*results <- fieldResult{Err: err}
		return
	}
	if child.Alias != "" {
		fieldName = child.Alias
	}
	*results <- fieldResult{
		FieldName: fieldName, Value: resolved, Err: err,
	}
}

func (e *Executor) resolveSelection(ctx context.Context, partial interface{}, field *graphql.Field, typeInfo schema.GraphQLTypeInfo, results *chan fieldResult, wg *sync.WaitGroup, child *graphql.Field)  (interface{}, error) {
	graphQLValue, _ := partial.(schema.GraphQLType)

	if child != nil {
		fieldName := child.Name
		_, ok := typeInfo.Fields[fieldName]
		if !ok {
			return nil, fmt.Errorf("No handler for field '%s' on type '%T'", fieldName, graphQLValue)
		}
		wg.Add(1)
		go e.callHandler(ctx, partial, field, typeInfo, results, wg, child)
	}
	return nil, nil
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
	
	fields := []*graphql.Field{}
	for _, selection := range field.SelectionSet {
		if selection.InlineFragment != nil {
			selectionSet := selection.InlineFragment.SelectionSet
			for _, selectionSetField := range selectionSet {
				fields = append(fields, selectionSetField.Field)
			}
		}
		
		if selection.Field != nil {
			fields = append(fields, selection.Field)
		}
	}
	
	

	result := map[string]interface{}{}
	typeInfo := schema.WithIntrospectionField(graphQLValue.GraphQLTypeInfo())
	results := make(chan fieldResult)
	wg := sync.WaitGroup{}

	for _, child := range fields {
		_, err := e.resolveSelection(ctx, partial, field, typeInfo, &results, &wg, child)
		if err != nil {
			return nil, err
		}
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		if r.Err != nil {
			return nil, r.Err
		}
		result[r.FieldName] = r.Value
	}
	return result, nil
}

func (e *Executor) resolveSlice(ctx context.Context, partials interface{}, field *graphql.Field) (interface{}, error) {
	v := reflect.ValueOf(partials)
	results := make([]interface{}, 0, v.Len())
	resChan := make(chan fieldResult)
	wg := sync.WaitGroup{}
	for i := 0; i < v.Len(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result, err := e.Resolve(ctx, v.Index(i).Interface(), field)
			resChan <- fieldResult{Value: result, Err: err}
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
		results = append(results, result.Value)
	}
	return results, nil
}
