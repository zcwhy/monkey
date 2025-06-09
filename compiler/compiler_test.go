// compiler/compiler_test.go

package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions code.Instructions
}

func parse(input string) *ast.Program {
	parser := parser.NewParser(lexer.NewLexer(input))
	return parser.ParseProgram()
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				if (true) { 10 };
				`,
			expectedConstants: []interface{}{10},
			expectedInstructions: code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 3),
				// 0003
				code.Make(code.OpConstant, 0),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
	let one = 1;
	let two = 2;
	`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
	let one = 1;
	one;
	`,
			expectedConstants: []interface{}{1}, expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
			}},
		{
			input: `
	let one = 1;
	let two = one;
	two;
	`,
			expectedConstants: []interface{}{1},
			expectedInstructions: code.Instructions{

				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
			}},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10; }`,
			expectedConstants: []interface{}{
				5,
				10,
				code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturn),
				}},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 2),
			}},
	}
	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)
	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	// last := compiler.scopes[compiler.scopeIndex].lastInstruction
	// if last.Opcode != code.OpSub {
	// 	t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
	// 		last.Opcode, code.OpSub)
	// }

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d",
			compiler.scopeIndex, 0)
	}

	compiler.emit(code.OpAdd)
	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	// last = compiler.scopes[compiler.scopeIndex].lastInstruction
	// if last.Opcode != code.OpAdd {
	// 	t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
	// 		last.Opcode, code.OpAdd)
	// }

	// previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	// if previous.Opcode != code.OpMul {
	// 	t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
	// 		previous.Opcode, code.OpMul)
	// }
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 24; }();`,
			expectedConstants: []interface{}{
				24,
				code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 1), // The compiled function
				code.Make(code.OpCall),
			},
		},
		{
			input: `
	let noArg = fn() { return 24; };
	noArg();
	`,
			expectedConstants: []interface{}{
				24,
				code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 1), // The compiled function
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall),
			}},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func testInstructions(
	expected code.Instructions,
	actual code.Instructions,
) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			code.Disassemble(expected), code.Disassemble(actual))
	}

	for i, ins := range expected {
		if string(actual[i]) != string(ins) {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, code.Disassemble(expected), code.Disassemble(actual))
		}
	}

	return nil
}

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		case code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T", i, actual[i])
			}

			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func TestCompile(t *testing.T) {
	input := `
					   let a = fn() { return 1 };
					   let b = fn() { return a() + 1 };
					   b();
	`

	program := parse(input)

	compiler := New()
	err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compiler error: %s", err)
	}

	bytecode := compiler.Bytecode()

	fmt.Println(code.Disassemble(bytecode.Instructions))
	fmt.Println(bytecode.Constants[0])
}
