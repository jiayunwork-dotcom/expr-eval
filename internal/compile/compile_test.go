package compile

import (
	"testing"

	"expr-eval/internal/parser"
)

func TestCompileSimple(t *testing.T) {
	node, _ := parser.Parse("2 + 3")
	prog, err := Compile(node)
	if err != nil {
		t.Fatal(err)
	}
	if prog.InstructionCount() < 3 {
		t.Errorf("instructions = %d, want >=3", prog.InstructionCount())
	}
}

func TestCompileWithVariable(t *testing.T) {
	node, _ := parser.Parse("x * 2")
	prog, err := Compile(node)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Variables) != 1 || prog.Variables[0] != "x" {
		t.Errorf("variables = %v", prog.Variables)
	}
}

func TestCompileFunction(t *testing.T) {
	node, _ := parser.Parse("min(a, b)")
	prog, err := Compile(node)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, inst := range prog.Instructions {
		if inst.Op == OpCall && inst.Name == "min" && inst.ArgCount == 2 {
			found = true
		}
	}
	if !found {
		t.Error("expected CALL min/2 instruction")
	}
}

func TestOptimizeConstantFolding(t *testing.T) {
	node, _ := parser.Parse("2 + 3")
	prog, _ := Compile(node)
	optimized := Optimize(prog)
	pushCount := 0
	for _, inst := range optimized.Instructions {
		if inst.Op == OpPush {
			pushCount++
			if inst.Operand != 5 {
				t.Errorf("folded value = %g, want 5", inst.Operand)
			}
		}
	}
	if pushCount != 1 {
		t.Errorf("pushes after fold = %d, want 1", pushCount)
	}
}

func TestOptimizeDoubleNeg(t *testing.T) {
	node, _ := parser.Parse("--x")
	prog, _ := Compile(node)
	optimized := Optimize(prog)
	for _, inst := range optimized.Instructions {
		if inst.Op == OpNeg {
			t.Error("double negation should be eliminated")
		}
	}
}

func TestDeadCodeElimination(t *testing.T) {
	node, _ := parser.Parse("1")
	prog, _ := Compile(node)
	prog.Instructions = append(prog.Instructions, Instruction{Op: OpPush, Operand: 999})
	cleaned := DeadCodeElimination(prog)
	last := cleaned.Instructions[len(cleaned.Instructions)-1]
	if last.Op != OpHalt {
		t.Error("expected HALT as last instruction after DCE")
	}
}

func TestProgramString(t *testing.T) {
	node, _ := parser.Parse("x + 1")
	prog, _ := Compile(node)
	s := prog.String()
	if s == "" {
		t.Error("disassembly should not be empty")
	}
}
