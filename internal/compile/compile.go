package compile

import (
	"fmt"

	"expr-eval/internal/parser"
)

type OpCode byte

const (
	OpPush OpCode = iota
	OpLoad
	OpNeg
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
	OpCall
	OpHalt
)

type Instruction struct {
	Op       OpCode
	Operand  float64
	Name     string
	ArgCount int
}

type Program struct {
	Instructions []Instruction
	Constants    []float64
	Variables    []string
}

func Compile(node parser.Node) (*Program, error) {
	c := &compiler{
		prog: &Program{},
	}
	if err := c.compile(node); err != nil {
		return nil, err
	}
	c.emit(Instruction{Op: OpHalt})
	return c.prog, nil
}

type compiler struct {
	prog *Program
}

func (c *compiler) emit(inst Instruction) {
	c.prog.Instructions = append(c.prog.Instructions, inst)
}

func (c *compiler) compile(node parser.Node) error {
	switch n := node.(type) {
	case *parser.NumberNode:
		c.emit(Instruction{Op: OpPush, Operand: n.Val})
		c.prog.Constants = append(c.prog.Constants, n.Val)
	case *parser.IdentNode:
		c.emit(Instruction{Op: OpLoad, Name: n.Name})
		c.addVar(n.Name)
	case *parser.UnaryNode:
		if err := c.compile(n.Expr); err != nil {
			return err
		}
		if n.Op == '-' {
			c.emit(Instruction{Op: OpNeg})
		}
	case *parser.BinaryNode:
		if err := c.compile(n.Left); err != nil {
			return err
		}
		if err := c.compile(n.Right); err != nil {
			return err
		}
		op, err := binOp(n.Op)
		if err != nil {
			return err
		}
		c.emit(Instruction{Op: op})
	case *parser.CallNode:
		for _, arg := range n.Args {
			if err := c.compile(arg); err != nil {
				return err
			}
		}
		c.emit(Instruction{Op: OpCall, Name: n.Name, ArgCount: len(n.Args)})
	default:
		return fmt.Errorf("compile: unsupported node type %T", node)
	}
	return nil
}

func (c *compiler) addVar(name string) {
	for _, v := range c.prog.Variables {
		if v == name {
			return
		}
	}
	c.prog.Variables = append(c.prog.Variables, name)
}

func binOp(op string) (OpCode, error) {
	switch op {
	case "+":
		return OpAdd, nil
	case "-":
		return OpSub, nil
	case "*":
		return OpMul, nil
	case "/":
		return OpDiv, nil
	case "%":
		return OpMod, nil
	case "^":
		return OpPow, nil
	default:
		return 0, fmt.Errorf("compile: unknown operator %q", op)
	}
}

func (p *Program) InstructionCount() int {
	return len(p.Instructions)
}

func (p *Program) String() string {
	var s string
	for i, inst := range p.Instructions {
		switch inst.Op {
		case OpPush:
			s += fmt.Sprintf("%04d PUSH %g\n", i, inst.Operand)
		case OpLoad:
			s += fmt.Sprintf("%04d LOAD %s\n", i, inst.Name)
		case OpNeg:
			s += fmt.Sprintf("%04d NEG\n", i)
		case OpAdd:
			s += fmt.Sprintf("%04d ADD\n", i)
		case OpSub:
			s += fmt.Sprintf("%04d SUB\n", i)
		case OpMul:
			s += fmt.Sprintf("%04d MUL\n", i)
		case OpDiv:
			s += fmt.Sprintf("%04d DIV\n", i)
		case OpMod:
			s += fmt.Sprintf("%04d MOD\n", i)
		case OpPow:
			s += fmt.Sprintf("%04d POW\n", i)
		case OpCall:
			s += fmt.Sprintf("%04d CALL %s/%d\n", i, inst.Name, inst.ArgCount)
		case OpHalt:
			s += fmt.Sprintf("%04d HALT\n", i)
		}
	}
	return s
}
