// generated by stringer -type OperationType; DO NOT EDIT

package ast // import "github.com/tallstreet/graphql/ast"

import "fmt"

const _OperationType_name = "QueryMutation"

var _OperationType_index = [...]uint8{0, 5, 13}

func (i OperationType) String() string {
	if i < 0 || i+1 >= OperationType(len(_OperationType_index)) {
		return fmt.Sprintf("OperationType(%d)", i)
	}
	return _OperationType_name[_OperationType_index[i]:_OperationType_index[i+1]]
}
