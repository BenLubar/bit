package main

import (
	"fmt"
	"io"
	"log"
)

func (p *Program) Compile(w io.Writer) error {
	if _, err := io.WriteString(w, assemblyLib); err != nil {
		return err
	}

	if len(p.Pointers) != 0 {
		if _, err := io.WriteString(w, "\n.data\n"); err != nil {
			return err
		}
		for i := range p.Pointers {
			if _, err := fmt.Fprintf(w, "ptr_%d:\n\t.quad 0\n", i); err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(w, "\n.text\n.globl _start\n_start:\n"); err != nil {
		return err
	}
	for _, l := range p.Lines {
		if _, err := fmt.Fprintf(w, ".L_l%s:\n", l.Num.shortString()); err != nil {
			return err
		}
		if err := p.compileStmt(w, l.Stmt); err != nil {
			return err
		}
		if l.Zero != nil && l.One != nil {
			if l.Zero.Equal(l.One) {
				if _, err := fmt.Fprintf(w, "\tjmp .L_l%s\n", l.Zero.shortString()); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjz .L_l%s\n\tjmp .L_l%s\n", l.Zero.shortString(), l.One.shortString()); err != nil {
					return err
				}
			}
		} else if l.Zero != nil {
			if _, err := fmt.Fprintf(w, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjz .L_l%s\n\tcall exit\n", l.Zero.shortString()); err != nil {
				return err
			}
		} else if l.One != nil {
			if _, err := fmt.Fprintf(w, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjnz .L_l%s\n\tcall exit\n", l.One.shortString()); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(w, "\tcall exit\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Program) compileStmt(w io.Writer, stmt Stmt) error {
	switch s := stmt.(type) {
	case *ReadStmt:
		if s.EOFLine != nil {
			if _, err := fmt.Fprintf(w, "\tmov $.L_l%s, %%rax\n", s.EOFLine.shortString()); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(w, "\txor %rax, %rax\n"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\tcall read\n\tmov %al, jump_register\n"); err != nil {
			return err
		}
	case *PrintStmt:
		if s.Bit {
			if _, err := io.WriteString(w, "\tmov $1, %rax\n"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(w, "\txor %rax, %rax\n"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\tcall print\n"); err != nil {
			return err
		}
	case *EqualsStmt:
		if s.Left.Pointer() && s.Right.Pointer() && s.Left.Value() && s.Right.Value() {
			log.Panicln("pointer or value???:", s)
		} else if s.Left.Pointer() && s.Right.Pointer() {
			if err := p.compileLValue(w, s.Left); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\tpush %rax\n"); err != nil {
				return err
			}
			if err := p.compileRValue(w, s.Right); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\tpop %rbx\n\tmov %rax, (%rbx)\n"); err != nil {
				return err
			}
		} else if s.Left.Value() && s.Right.Value() {
			if err := p.compileLValue(w, s.Left); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\tpush %rax\n"); err != nil {
				return err
			}
			if err := p.compileRValue(w, s.Right); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\tpop %rbx\n\tmov %al, (%rbx)\n"); err != nil {
				return err
			}
		} else {
			log.Panicln("assign???:", s)
		}
	default:
		panic("unreachable")
	}
	return nil
}

func (p *Program) compileLValue(w io.Writer, expr Expr) error {
	switch e := expr.(type) {
	case *PointerVariable:
		if p.Pointers.Contains(e.Num) {
			if _, err := fmt.Fprintf(w, "\tmov $ptr_%d, %%rax\n", p.Pointers.Search(e.Num)); err != nil {
				return err
			}
		} else {
			log.Panicln("pointer???:", expr)
		}
	case *BitVariable:
		if _, err := fmt.Fprintf(w, "\tmov $%d, %%rax\n\tcall var\n", e.Num.Uint64()); err != nil {
			return err
		}
	case *JumpRegister:
		if _, err := io.WriteString(w, "\tmov $jump_register, %rax\n"); err != nil {
			return err
		}
	case *ValueAt:
		if err := p.compileRValue(w, e.Target); err != nil {
			return err
		}
		if e.Offset != 0 {
			if _, err := fmt.Fprintf(w, "\tadd $%d, %%rax\n\tcall ensure\n", e.Offset); err != nil {
				return err
			}
		}
	case *AddressOf:
		log.Panicln("invalid lvalue:", e)
	case *BitValue:
		log.Panicln("invalid lvalue:", e)
	case *Nand:
		log.Panicln("invalid lvalue:", e)
	case *Parenthesis:
		return p.compileLValue(w, e.Inner)
	default:
		panic("unreachable")
	}
	return nil
}

func (p *Program) compileRValue(w io.Writer, expr Expr) error {
	switch e := expr.(type) {
	case *PointerVariable:
		if p.Pointers.Contains(e.Num) {
			if _, err := fmt.Fprintf(w, "\tmov ptr_%d, %%rax\n", p.Pointers.Search(e.Num)); err != nil {
				return err
			}
		} else {
			log.Panicln("pointer???:", expr)
		}
	case *BitVariable:
		if _, err := fmt.Fprintf(w, "\tmov $%d, %%rax\n\tcall var\n\tmov (%%rax), %%bl\n\txor %%rax, %%rax\n\tmov %%bl, %%al\n", e.Num.Uint64()); err != nil {
			return err
		}
	case *JumpRegister:
		if _, err := io.WriteString(w, "\txor %rax, %rax\n\tmov jump_register, %al\n"); err != nil {
			return err
		}
	case *ValueAt:
		if err := p.compileRValue(w, e.Target); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "\tmov %d(%%rax), %%bl\n\txor %%rax, %%rax\n\tmov %%bl, %%al\n", e.Offset); err != nil {
			return err
		}
	case *AddressOf:
		if err := p.compileLValue(w, e.Variable); err != nil {
			return err
		}
	case *BitValue:
		if e.Bit {
			if _, err := io.WriteString(w, "\tmov $1, %rax\n"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(w, "\txor %rax, %rax\n"); err != nil {
				return err
			}
		}
	case *Nand:
		if err := p.compileRValue(w, e.Left); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\tpush %rax\n"); err != nil {
			return err
		}
		if err := p.compileRValue(w, e.Right); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\tpop %rbx\n\tand %rbx, %rax\n\txor $1, %rax\n"); err != nil {
			return err
		}
	case *Parenthesis:
		return p.compileRValue(w, e.Inner)
	default:
		panic("unreachable")
	}
	return nil
}
