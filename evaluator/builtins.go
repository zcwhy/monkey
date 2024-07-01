package evaluator

import (
	"monkey/object"
)

func InstallBuiltinFn(env *object.Environment) {
	env.Set("len", LenBuiltin())

	// array builtin
	env.Set("first", FirstBuiltin())
	env.Set("last", LastBuiltin())
	env.Set("rest", RestBuiltin())
	env.Set("push", PushBuiltin())
}

func LenBuiltin() *object.Builtin {
	lenBuiltin := func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].Type())
		}
	}
	return &object.Builtin{Fn: lenBuiltin}
}

func FirstBuiltin() *object.Builtin {
	firstBuiltin := func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return NULL
	}
	return &object.Builtin{Fn: firstBuiltin}
}

func LastBuiltin() *object.Builtin {
	lastBuiltin := func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		l := len(arr.Elements)
		if l > 0 {
			return arr.Elements[l-1]
		}
		return NULL
	}
	return &object.Builtin{Fn: lastBuiltin}
}

func RestBuiltin() *object.Builtin {
	restBuiltin := func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		l := len(arr.Elements)
		if l > 0 {
			newElements := make([]object.Object, l-1, l-1)
			copy(newElements, arr.Elements[1:l])
			return &object.Array{Elements: newElements}
		}
		return NULL
	}
	return &object.Builtin{Fn: restBuiltin}
}

func PushBuiltin() *object.Builtin {
	pushBuiltin := func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		l := len(arr.Elements)
		if l > 0 {
			newElements := make([]object.Object, l+1, l+1)
			copy(newElements, arr.Elements)
			newElements[l] = args[1]
			return &object.Array{Elements: newElements}
		}
		return NULL
	}
	return &object.Builtin{Fn: pushBuiltin}
}
