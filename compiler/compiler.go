package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type Compiler struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func parse(input string) *ast.Program {
	parser := parser.NewParser(lexer.NewLexer(input))
	return parser.ParseProgram()
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}

	return nil
}
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.Instructions,
		Constants:    c.Constants,
	}
}

func (c *Compiler) addConstant(constant object.Object) int {
	c.Constants = append(c.Constants, constant)
	return len(c.Constants) - 1
}

func (c *Compiler) emit(opCode code.OpCode, opreands ...int) {
	c.Instructions = append(c.Instructions, code.Make(opCode, opreands...))
}
