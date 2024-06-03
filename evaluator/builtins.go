package evaluator

import "monkey/object"

func InstallBuiltinFn(env *object.Environment) {
	env.Set("len", LenBuiltin())
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
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].Type())
		}
	}
	return &object.Builtin{Fn: lenBuiltin}
}
