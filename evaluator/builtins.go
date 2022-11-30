package evaluator

import "github.com/henningstorck/monkey-interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Function: func(args ...object.Object) object.Object {
			if err := expectArguments(args, 1); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("invalid argument. got %s", args[0].Type())
			}
		},
	},
	"first": {
		Function: func(args ...object.Object) object.Object {
			if err := expectArguments(args, 1); err != nil {
				return err
			}

			if err := expectArgumentType(args[0], object.ArrayObj); err != nil {
				return err
			}

			arr := args[0].(*object.Array)

			if (len(arr.Elements)) > 0 {
				return arr.Elements[0]
			}

			return NullObj
		},
	},
	"last": {
		Function: func(args ...object.Object) object.Object {
			if err := expectArguments(args, 1); err != nil {
				return err
			}

			if err := expectArgumentType(args[0], object.ArrayObj); err != nil {
				return err
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			if (length) > 0 {
				return arr.Elements[length-1]
			}

			return NullObj
		},
	},
	"rest": {
		Function: func(args ...object.Object) object.Object {
			if err := expectArguments(args, 1); err != nil {
				return err
			}

			if err := expectArgumentType(args[0], object.ArrayObj); err != nil {
				return err
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			if (length) > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NullObj
		},
	},
	"push": {
		Function: func(args ...object.Object) object.Object {
			if err := expectArguments(args, 2); err != nil {
				return err
			}

			if err := expectArgumentType(args[0], object.ArrayObj); err != nil {
				return err
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
}

func expectArguments(args []object.Object, expected int) *object.Error {
	if len(args) != expected {
		return newError("wrong number of arguments. got %d, but expected %d", len(args), expected)
	}

	return nil
}

func expectArgumentType(arg object.Object, objType object.ObjectType) *object.Error {
	if arg.Type() != object.ArrayObj {
		return newError("invalid argument. got %s, but expected %s", arg.Type(), objType)
	}

	return nil
}
