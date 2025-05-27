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
		case 2:
			binary.BigEndian.PutUint16(instruction[pos:], uint16(opreand))
		}

		pos += iWidth
	}
	return instruction
}

func Disassemble(instructions [][]byte) string {
	var out bytes.Buffer

	pos := 0
	for _, instruction := range instructions {
		def, _ := LookUp(OpCode(instruction[0]))
		opreands, err := readOpreands(instruction)
		if err != nil {
			return ""
		}

		fmt.Fprintf(&out, "%04d %s %s\n", pos, def.Name, fomatOpreands(opreands))
		pos += len(instruction)
	}

	return out.String()
}

func readOpreands(instruction []byte) ([]int, error) {
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
		}
		offset += width
	}
	return opreands, nil
}

func fomatOpreands(opreands []int) string {
	var out bytes.Buffer
	for _, opreand := range opreands {
		fmt.Fprintf(&out, "%d ", opreand)
	}
	return out.String()
}
