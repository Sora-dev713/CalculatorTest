package expression

import (
	"fmt"
	"strconv"
	"unicode"
)

type Kind int

const (
	Number Kind = iota
	Binary
	Unary
	Function
)

type Node struct {
	Kind        Kind
	Value       float64
	Operator    rune
	Name        string
	Left, Right *Node
	Args        []*Node
}

type tokenKind int

const (
	tEOF tokenKind = iota
	tNumber
	tIdentifier
	tPlus
	tMinus
	tMultiply
	tDivide
	tPower
	tLeftParen
	tRightParen
	tComma
)

type token struct {
	kind     tokenKind
	text     string
	position int
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	for i := 0; i < len(input); {
		r := rune(input[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		start := i
		if (r >= '0' && r <= '9') || r == '.' {
			dots, digits := 0, 0
			for i < len(input) {
				c := input[i]
				if c == '.' {
					dots++
					i++
					continue
				}
				if c < '0' || c > '9' {
					break
				}
				digits++
				i++
			}
			text := input[start:i]
			if dots > 1 || digits == 0 {
				return nil, fmt.Errorf("invalid number at position %d", start+1)
			}
			tokens = append(tokens, token{tNumber, text, start})
			continue
		}
		if unicode.IsLetter(r) {
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			tokens = append(tokens, token{tIdentifier, input[start:i], start})
			continue
		}
		kinds := map[byte]tokenKind{'+': tPlus, '-': tMinus, '*': tMultiply, '/': tDivide, '^': tPower, '(': tLeftParen, ')': tRightParen, ',': tComma}
		kind, ok := kinds[input[i]]
		if !ok {
			return nil, fmt.Errorf("unsupported character %q at position %d", input[i], i+1)
		}
		tokens = append(tokens, token{kind, input[i : i+1], i})
		i++
	}
	tokens = append(tokens, token{kind: tEOF, position: len(input)})
	return tokens, nil
}

type parser struct {
	tokens  []token
	current int
}

func Parse(input string) (*Node, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("expression cannot be empty")
	}
	if len(input) > 512 {
		return nil, fmt.Errorf("expression exceeds the 512 character limit")
	}
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens}
	node, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tEOF {
		return nil, fmt.Errorf("unexpected %q at position %d", p.peek().text, p.peek().position+1)
	}
	return node, nil
}

func (p *parser) expression() (*Node, error) { return p.addition() }
func (p *parser) addition() (*Node, error) {
	left, err := p.multiplication()
	if err != nil {
		return nil, err
	}
	for p.match(tPlus, tMinus) {
		op := rune(p.previous().text[0])
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		left = &Node{Kind: Binary, Operator: op, Left: left, Right: right}
	}
	return left, nil
}
func (p *parser) multiplication() (*Node, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(tMultiply, tDivide) {
		op := rune(p.previous().text[0])
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		left = &Node{Kind: Binary, Operator: op, Left: left, Right: right}
	}
	return left, nil
}
func (p *parser) unary() (*Node, error) {
	if p.match(tPlus, tMinus) {
		op := rune(p.previous().text[0])
		operand, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Node{Kind: Unary, Operator: op, Right: operand}, nil
	}
	return p.power()
}
func (p *parser) power() (*Node, error) {
	left, err := p.primary()
	if err != nil {
		return nil, err
	}
	if p.match(tPower) {
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Node{Kind: Binary, Operator: '^', Left: left, Right: right}, nil
	}
	return left, nil
}
func (p *parser) primary() (*Node, error) {
	if p.match(tNumber) {
		value, _ := strconv.ParseFloat(p.previous().text, 64)
		return &Node{Kind: Number, Value: value}, nil
	}
	if p.match(tIdentifier) {
		name := p.previous().text
		if name != "sqrt" && name != "percent" {
			return nil, fmt.Errorf("unknown function %q", name)
		}
		if !p.match(tLeftParen) {
			return nil, fmt.Errorf("expected '(' after %s", name)
		}
		first, err := p.expression()
		if err != nil {
			return nil, err
		}
		args := []*Node{first}
		if name == "percent" {
			if !p.match(tComma) {
				return nil, fmt.Errorf("percent requires two arguments separated by a comma")
			}
			second, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, second)
		}
		if !p.match(tRightParen) {
			return nil, fmt.Errorf("expected ')' to close %s", name)
		}
		return &Node{Kind: Function, Name: name, Args: args}, nil
	}
	if p.match(tLeftParen) {
		node, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(tRightParen) {
			return nil, fmt.Errorf("expected ')' at position %d", p.peek().position+1)
		}
		return node, nil
	}
	t := p.peek()
	if t.kind == tEOF {
		return nil, fmt.Errorf("expected a number, function, or '(' at end of expression")
	}
	return nil, fmt.Errorf("expected a number, function, or '(' at position %d", t.position+1)
}
func (p *parser) match(kinds ...tokenKind) bool {
	for _, k := range kinds {
		if p.peek().kind == k {
			p.current++
			return true
		}
	}
	return false
}
func (p *parser) peek() token     { return p.tokens[p.current] }
func (p *parser) previous() token { return p.tokens[p.current-1] }
