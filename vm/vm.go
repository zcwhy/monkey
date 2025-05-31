package vm

import (
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

type VM struct {
	intructions code.Instructions
	constant    []object.Object

	stack []object.Object
	sp    int
}

func New(byteCode *compiler.Bytecode) *VM {
	return &VM{
		intructions: byteCode.Instructions,
		constant:    byteCode.Constants,

		sp: -1,
	}
}

func (v *VM) Run() error {
	for _, instruction := range v.intructions {
		opCode := instruction[0]

		switch code.OpCode(opCode) {
		case code.OpConstant:
			opreandNo, err := code.ReadOpreands(instruction)
			if err != nil {

			}

			v.push(v.constant[opreandNo[0]])
		case code.OpAdd:
			op1 := v.pop()
			op2 := v.pop()

			if op1.Type() == object.INTEGER_OBJ && op2.Type() == object.INTEGER_OBJ {
				leftVal := op1.(*object.Integer).Value
				rightVal := op2.(*object.Integer).Value
				result := &object.Integer{Value: leftVal + rightVal}
				v.push(result)
			}
		}
	}
	return nil
}

func (v *VM) push(elem ...object.Object) {
	for _, e := range elem {
		v.stack = append(v.stack, e)
		v.sp = len(v.stack) - 1
	}
}

func (v *VM) pop() object.Object {
	if v.sp < 0 {
		return nil
	}

	o := v.stack[v.sp]
	v.sp -= 1

	return o
}

func (v *VM) StackTop() object.Object {
	return v.stack[v.sp]
}
