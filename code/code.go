package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions [][]byte

type OpCode byte

const (
	OpConstant OpCode = iota
	OpAdd
	OpSub
	OpMul
	OpDev
	OpTrue
	OpJumpNotTruthy
	OpSetGlobal
	OpGetGlobal
	OpReturn
	OpCall
	OpGetLocal
	OpSetLocal
)

func Make(opCode OpCode, opreands ...int) []byte {
	def, err := LookUp(opCode)
	if err != nil {
		return nil
	}

	instructionLen := 1
	for _, width := range def.OpreandsWidth {
		instructionLen += width
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(opCode)

	pos := 1
	for i, opreand := range opreands {
		iWidth := def.OpreandsWidth[i]
		switch iWidth {
		case 1:
			instruction[pos] = byte(opreand)
		case 2:
			binary.BigEndian.PutUint16(instruction[pos:], uint16(opreand))
		}

		pos += iWidth
	}
	return instruction
}

func Disassemble(instructions Instructions) string {
	var out bytes.Buffer

	pos := 0
	fmt.Fprintf(&out, "\n")
	for _, instruction := range instructions {
		def, _ := LookUp(OpCode(instruction[0]))
		opreands, err := ReadOpreands(instruction)
		if err != nil {
			return ""
		}

		fmt.Fprintf(&out, "%04d %s %s\n", pos, def.Name, fomatOpreands(opreands))
		pos += len(instruction)
	}

	return out.String()
}

func ReadOpreands(instruction []byte) ([]int, error) {
	def, err := LookUp(OpCode(instruction[0]))
	if err != nil {
		return nil, err
	}

	offset := 1
	opreands := []int{}
	for _, width := range def.OpreandsWidth {
		switch width {
		case 2:
			opreands = append(opreands, int(binary.BigEndian.Uint16(instruction[offset:])))
		case 1:
			opreands = append(opreands, int(instruction[offset]))
		}
		offset += width
	}
	return opreands, nil
}

func ReadUint16(instruction []byte) uint16 {
	return binary.BigEndian.Uint16(instruction)
}

func fomatOpreands(opreands []int) string {
	var out bytes.Buffer
	for _, opreand := range opreands {
		fmt.Fprintf(&out, "%d ", opreand)
	}
	return out.String()
}
