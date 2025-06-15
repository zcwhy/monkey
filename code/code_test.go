package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       OpCode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDev, []int{}, []byte{byte(OpDev)}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpJumpNotTruthy, []int{123}, []byte{byte(OpJumpNotTruthy), 0, 123}},
		{OpSetGlobal, []int{123}, []byte{byte(OpSetGlobal), 0, 123}},
		{OpGetGlobal, []int{123}, []byte{byte(OpGetGlobal), 0, 123}},
		{OpReturn, []int{}, []byte{byte(OpReturn)}},
		{OpCall, []int{1}, []byte{byte(OpCall), 1}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("%s, wrong byte at pos %d. want=%d, got=%d",
					instructionDefinitions[tt.op].Name, i, b, instruction[i])
			}
		}
	}
}

// func TestReadOperands(t *testing.T) {
// 	tests := []struct {
// 		op        OpCode
// 		operands  []int
// 		bytesRead int
// 	}{
// 		{OpGetLocal, []int{255}, 1},
// 	}

// 	for _, tt := range tests {
// 		ReadOpreands()
// 		instruction := Make(tt.op, tt.operands...)

// 		if len(instruction) != len(tt.expected) {
// 			t.Errorf("instruction has wrong length. want=%d, got=%d",
// 				len(tt.expected), len(instruction))
// 		}

// 		for i, b := range tt.expected {
// 			if instruction[i] != tt.expected[i] {
// 				t.Errorf("%s, wrong byte at pos %d. want=%d, got=%d",
// 					string(tt.op), i, b, instruction[i])
// 			}
// 		}
// 	}
// }

func TestFormatOreands(t *testing.T) {
	t.Log(fomatOpreands([]int{1, 2, 3}))
}

func TestDisaseemble(t *testing.T) {
	instructions := Instructions{
		// 0000
		Make(OpTrue),
		// 0001
		Make(OpJumpNotTruthy, 7),
		// 0004
		Make(OpConstant, 0),
	}

	t.Log(Disassemble(instructions))
}
