package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct {
	Instructions code.Instructions
	Constants    []object.Object
	SymbolTable  *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type CompilationScope struct {
	instructions code.Instructions
}

func New() *Compiler {
	return &Compiler{
		SymbolTable: NewSymbolTable(),
		scopes: []CompilationScope{
			{
				instructions: code.Instructions{},
			},
		},
		scopeIndex: 0,
	}
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

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDev)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		pos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		c.changeOpreand(pos, len(c.Instructions))

		// err = c.Compile(node.Alternative)
		// if err != nil {
		// 	return err
		// }

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		entry := c.SymbolTable.Set(node.Name.Value, GlobalScope)
		c.emit(code.OpSetGlobal, entry.Index)

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		}
	case *ast.Identifier:
		entry, ok := c.SymbolTable.Get(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.emit(code.OpGetGlobal, entry.Index)

	case *ast.FunctionLiteral:
		c.enterScope()

		c.Compile(node.Body)

		compiledFunc := &object.CompiledFunction{
			Instructions: c.scopes[c.scopeIndex].instructions,
		}

		c.leaveScope()

		c.emit(code.OpConstant, c.addConstant(compiledFunc))
	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}
		c.emit(code.OpReturn)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		c.emit(code.OpCall)
	}

	return nil
}
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.Constants,
	}
}

func (c *Compiler) addConstant(constant object.Object) int {
	c.Constants = append(c.Constants, constant)
	return len(c.Constants) - 1
}

func (c *Compiler) emit(opCode code.OpCode, opreands ...int) int {
	instruciton := code.Make(opCode, opreands...)
	c.scopes[c.scopeIndex].instructions = append(c.scopes[c.scopeIndex].instructions, instruciton)
	return len(c.scopes[c.scopeIndex].instructions) - 1
}

func (c *Compiler) changeOpreand(pos int, opreand ...int) {
	opCode := code.OpCode(c.Instructions[pos][0])

	newInstruction := code.Make(opCode, opreand...)
	c.replaceInstrucion(pos, newInstruction)
}

func (c *Compiler) replaceInstrucion(pos int, newInstruction []byte) {
	c.Instructions[pos] = newInstruction
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	newScope := CompilationScope{
		instructions: code.Instructions{},
	}
	c.scopes = append(c.scopes, newScope)
	c.scopeIndex += 1
}

func (c *Compiler) leaveScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex -= 1
}
