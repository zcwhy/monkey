package code

import "fmt"

type Definition struct {
	Name          string
	OpreandsWidth []int
}

var instructionDefinitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
	OpSub:      {"OpSub", []int{}},
	OpMul:      {"OpSub", []int{}},
	OpDev:      {"OpSub", []int{}},
}

func LookUp(code OpCode) (*Definition, error) {
	def, ok := instructionDefinitions[code]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", code)
	}
	return def, nil
}
