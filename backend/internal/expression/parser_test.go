package expression

import "testing"

func TestParseValidExpressions(t *testing.T) {
	for _, input := range []string{"2+3*4", "(2+3)*4", "2^3^2", "sqrt(16)+2", "percent(200,10)", "-2^2", "2^-2", "percent(sqrt(16),-10)"} {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestPowerIsRightAssociative(t *testing.T) {
	node, err := Parse("2^3^2")
	if err != nil {
		t.Fatal(err)
	}
	if node.Operator != '^' || node.Right == nil || node.Right.Operator != '^' {
		t.Fatalf("unexpected tree: %#v", node)
	}
}
func TestUnaryHasLowerPrecedenceThanPower(t *testing.T) {
	node, err := Parse("-2^2")
	if err != nil {
		t.Fatal(err)
	}
	if node.Kind != Unary || node.Right.Operator != '^' {
		t.Fatalf("unexpected tree: %#v", node)
	}
}
func TestParseRejectsInvalidExpressions(t *testing.T) {
	for _, input := range []string{"", "   ", "2+", "(2+3", "2+3)", "eval(1)", "sqrt()", "percent(2)", "percent(2,3,4)", "2%3", "1..2", "2 3"} {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
func TestLengthLimit(t *testing.T) {
	input := make([]byte, 513)
	for i := range input {
		input[i] = '1'
	}
	if _, err := Parse(string(input)); err == nil {
		t.Fatal("expected error")
	}
}
