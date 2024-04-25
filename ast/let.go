package ast

import (
	"bytes"
	"monkey/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) statementNode() {

}

func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.Token.Literal + " ")
	out.WriteString(l.Name.Value)
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")

	return out.String()
}
