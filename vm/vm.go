package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
	FrameSize   = 1024
)

type VM struct {
	constant []object.Object
	globals  []object.Object

	stack []object.Object
	sp    int

	frames     []*Frame
	frameIndex int
}

// call frame
type Frame struct {
	fn *object.CompiledFunction
	ip int // return address
}

func New(byteCode *compiler.Bytecode) *VM {
	mainFrame := Frame{fn: &object.CompiledFunction{Instructions: byteCode.Instructions}, ip: 0}

	vm := &VM{
		constant: byteCode.Constants,
		globals:  make([]object.Object, GlobalsSize),

		stack: make([]object.Object, 2048),
		sp:    -1,

		frames:     make([]*Frame, FrameSize),
		frameIndex: 0,
	}

	vm.frames[0] = &mainFrame
	return vm
}

func (v *VM) Run() error {
	for v.currentFrame().ip < len(v.currentFrame().fn.Instructions) {
		ip := v.currentFrame().ip
		instruction := v.currentFrame().fn.Instructions[ip]
		opCode := code.OpCode(instruction[0])

		v.currentFrame().ip++
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
		case code.OpTrue:
			v.push(&object.Boolean{Value: true})

		case code.OpSetGlobal:
			symbolIndex := code.ReadUint16(instruction[1:])
			v.globals[symbolIndex] = v.pop()

		case code.OpGetGlobal:
			symbolIndex := code.ReadUint16(instruction[1:])
			v.push(v.globals[symbolIndex])

		case code.OpCall:
			fn, ok := v.pop().(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}

			v.pushFrame(&Frame{fn: fn})

		case code.OpReturn:
			v.popFrame()
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

func (v *VM) currentFrame() *Frame {
	return v.frames[v.frameIndex]
}

func (v *VM) pushFrame(f *Frame) {
	v.frames[v.frameIndex+1] = f
	v.frameIndex += 1
}

func (v *VM) popFrame() *Frame {
	lastFrame := v.frames[v.frameIndex]
	v.frameIndex -= 1

	return lastFrame
}
