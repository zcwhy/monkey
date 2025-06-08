package code

import "fmt"

type Definition struct {
	Name          string
	OpreandsWidth []int
}

var instructionDefinitions = map[OpCode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpSub", []int{}},
	OpDev:           {"OpSub", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}}, // 这里操作数是跳转到第几条指令
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpReturn:        {"OpReturn", []int{}},
	OpCall:          {"OpCall", []int{}},
}

func LookUp(code OpCode) (*Definition, error) {
	def, ok := instructionDefinitions[code]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", code)
	}
	return def, nil
}
