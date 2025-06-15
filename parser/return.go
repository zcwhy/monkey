package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	// link `return;`
	if p.curTokenIs(token.SEMICOLON) {
		return stmt
	}
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}
