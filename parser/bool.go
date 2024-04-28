package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseBoolean() ast.Expression {
	stmt := &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return stmt
}
