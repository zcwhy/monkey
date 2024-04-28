package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	// 这里右括号的优先级是LOWEST的（所有未显式声明优先级的token均会分会LOWEST），所以在解析遇到LPARTEN时会返回当前表达式
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
