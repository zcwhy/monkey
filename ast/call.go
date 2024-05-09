package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type CallExpression struct {
	Token     token.Token // '('词法单元
	Function  Expression  // 标识符或函数字⾯量
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
