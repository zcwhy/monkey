package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

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
		stack:       make([]object.Object, 2048),
		sp:          -1,
	}
}

func (v *VM) Run() error {
	for _, instruction := range v.intructions {
		opCode := code.OpCode(instruction[0])

		switch opCode {
		case code.OpConstant:
			opreandNo, err := code.ReadOpreands(instruction)
			if err != nil {

			}

			v.push(v.constant[opreandNo[0]])
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDev:
			err := v.executeBinaryOperation(opCode)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (v *VM) executeBinaryOperation(op code.OpCode) error {
	// 注意pop的顺序，先pop的是right
	right := v.pop()
	left := v.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return v.executeBinaryIntegerOperation(op, left.(*object.Integer), right.(*object.Integer))
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (v *VM) executeBinaryIntegerOperation(op code.OpCode, left, right *object.Integer) error {
	leftValue := left.Value
	rightValue := right.Value

	// fmt.Println(leftValue, " ", rightValue)

	switch op {
	case code.OpAdd:
		v.push(&object.Integer{Value: leftValue + rightValue})
	case code.OpSub:
		v.push(&object.Integer{Value: leftValue - rightValue})
	case code.OpMul:
		v.push(&object.Integer{Value: leftValue * rightValue})
	case code.OpDev:
		v.push(&object.Integer{Value: leftValue / rightValue})
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return nil
}

func (v *VM) push(elem ...object.Object) {
	for _, e := range elem {
		v.stack[v.sp+1] = e
		v.sp += 1
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
