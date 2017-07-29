package main

import (
	"fmt"
	"io"
	"log"

	"github.com/BenLubar/bit/bitnum"
)

const linux64Runtime = `
.data
.L_start_of_heap:
	.quad 0
.L_end_of_heap:
	.quad 0

.text
.globl var
var:
	mov .L_start_of_heap, %rbx
	test %rbx, %rbx
	jnz .L_var_ready
	push %rax
	mov $_end, %rdi
	mov $0xc, %rax
	syscall
	mov %rax, %rbx
	mov %rax, .L_start_of_heap
	mov %rax, .L_end_of_heap
	pop %rax
.L_var_ready:
	add %rbx, %rax
	call ensure
	ret

.text
.globl ensure
ensure:
	mov .L_end_of_heap, %rdi
	cmp %rdi, %rax
	jl .L_ensure_ret
	push %rax
	push %rdi
	mov %rax, %rdi
	add $0x1001, %rdi
	and $0xfffffffffffff000, %rdi
	push %rdi
	mov $0xc, %rax
	syscall
	mov %rax, .L_end_of_heap
	pop %rcx
	pop %rdi
	xor %rax, %rax
	sub %rdi, %rcx
	rep stosb
	pop %rax
.L_ensure_ret:
	ret

.data
.L_print_bits:
	.byte 0
.L_print_accum:
	.byte 0

.text
.globl print
print:
	xor %rbx, %rbx
	xor %rcx, %rcx
	mov .L_print_bits, %bl
	mov .L_print_accum, %cl
	shl $1, %rcx
	inc %rbx
	or %rax, %rcx
	mov %cl, .L_print_accum
	cmp $8, %rbx
	je .L_print_out
	mov %bl, .L_print_bits
	ret
.L_print_out:
	mov $1, %rax
	mov $1, %rdi
	mov $.L_print_accum, %rsi
	mov $1, %rdx
	syscall
	movb $0, .L_print_bits
	movb $0, .L_print_accum
	ret

.data
.L_read_bits:
	.byte 0
.L_read_accum:
	.byte 0

.text
.globl read
read:
	xor %rbx, %rbx
	mov .L_read_bits, %bl
	test %rbx, %rbx
	jz .L_read_in
	dec %rbx
	mov %bl, .L_read_bits
	mov .L_read_accum, %bl
	mov %rbx, %rax
	and $1, %rax
	shr $1, %rbx
	mov %bl, .L_read_accum
	ret
.L_read_in:
	push %rax
	mov $0, %rax
	mov $0, %rdi
	mov $.L_read_accum, %rsi
	mov $1, %rdx
	syscall
	test %rax, %rax
	jz .L_read_fail
	pop %rax
	movb $8, .L_read_bits
	jmp read
.L_read_fail:
	pop %rax
	test %rax, %rax
	jz exit
	mov %rax, (%rsp)
	ret

.text
.globl exit
exit:
	xor %rax, %rax
	mov .L_print_bits, %al
	test %rax, %rax
	jz .L_exit_exit
	mov $7, %rcx
	sub %rax, %rcx
	mov .L_print_accum, %al
	shl %cl, %rax
	mov %al, .L_print_accum
	movb $7, .L_print_bits
	xor %rax, %rax
	call print
.L_exit_exit:
	mov $60, %rax
	mov $0, %rdi
	syscall

.data
.globl jump_register
jump_register:
	.byte 0
`

type Linux64AssemblyWriter struct {
	W io.Writer
}

func (w Linux64AssemblyWriter) Runtime() error {
	_, err := io.WriteString(w.W, linux64Runtime)
	return err
}

func (w Linux64AssemblyWriter) DataSegment() error {
	_, err := io.WriteString(w.W, "\n.data\n")
	return err
}

func (w Linux64AssemblyWriter) TextSegment() error {
	_, err := io.WriteString(w.W, "\n.text\n")
	return err
}

func (w Linux64AssemblyWriter) DeclarePointer(n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "ptr_%s:\n\t.quad 0\n", n.ShortString())
	return err
}

func (w Linux64AssemblyWriter) Start() error {
	_, err := io.WriteString(w.W, ".globl _start\n_start:\n")
	return err
}

func (w Linux64AssemblyWriter) DeclareLine(n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, ".L_l%s:\n", n.ShortString())
	return err
}

func (w Linux64AssemblyWriter) Goto(zero, one *bitnum.Number) error {
	if zero != nil && one != nil {
		if zero.Equal(one) {
			_, err := fmt.Fprintf(w.W, "\tjmp .L_l%s\n", zero.ShortString())
			return err
		}
		_, err := fmt.Fprintf(w.W, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjz .L_l%s\n\tjmp .L_l%s\n", zero.ShortString(), one.ShortString())
		return err
	}
	if zero != nil {
		_, err := fmt.Fprintf(w.W, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjz .L_l%s\n\tcall exit\n", zero.ShortString())
		return err
	}
	if one != nil {
		_, err := fmt.Fprintf(w.W, "\tmov jump_register, %%al\n\ttest %%al, %%al\n\tjnz .L_l%s\n\tcall exit\n", one.ShortString())
		return err
	}
	_, err := io.WriteString(w.W, "\tcall exit\n")
	return err
}

func (w Linux64AssemblyWriter) Read(eof *bitnum.Number) error {
	if eof != nil {
		if _, err := fmt.Fprintf(w.W, "\tmov $.L_l%s, %%rax\n", eof.ShortString()); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w.W, "\txor %rax, %rax\n"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w.W, "\tcall read\n\tmov %al, jump_register\n")
	return err
}

func (w Linux64AssemblyWriter) BitValue(register int, bit bool) error {
	if bit {
		_, err := fmt.Fprintf(w.W, "\tmov $1, %%r%cx\n", register+'a')
		return err

	}
	_, err := fmt.Fprintf(w.W, "\txor %%r%cx, %%r%cx\n", register+'a', register+'a')
	return err
}

func (w Linux64AssemblyWriter) Print() error {
	_, err := io.WriteString(w.W, "\tcall print\n")
	return err
}

func (w Linux64AssemblyWriter) SaveRegister(register int) error {
	_, err := fmt.Fprintf(w.W, "\tpush %%r%cx\n", register+'a')
	return err
}

func (w Linux64AssemblyWriter) LoadRegister(register int) error {
	_, err := fmt.Fprintf(w.W, "\tpop %%r%cx\n", register+'a')
	return err
}

func (w Linux64AssemblyWriter) WritePointer(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov %%r%cx, (%%r%cx)\n", src+'a', dest+'a')
	return err
}

func (w Linux64AssemblyWriter) WriteBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov %%%cl, (%%r%cx)\n", src+'a', dest+'a')
	return err
}

func (w Linux64AssemblyWriter) PointerAddress(reg int, n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "\tmov $ptr_%s, %%r%cx\n", n.ShortString(), reg+'a')
	return err
}

func (w Linux64AssemblyWriter) BitAddress(reg int, n *bitnum.Number) error {
	if reg != 0 {
		log.Panicln("internal compiler error: unsupported BitAddress register for Linux64AssemblyWriter:", reg)
	}
	_, err := fmt.Fprintf(w.W, "\tmov $%d, %%rax\n\tcall var\n", n.Uint64())
	return err
}

func (w Linux64AssemblyWriter) PointerValue(reg int, n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "\tmov ptr_%s, %%r%cx\n", n.ShortString(), reg+'a')
	return err
}

func (w Linux64AssemblyWriter) ReadBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov (%%r%cx), %%%cl\n\tand $1, %%r%cx\n", src+'a', dest+'a', dest+'a')
	return err
}

func (w Linux64AssemblyWriter) JumpAddress(reg int) error {
	_, err := fmt.Fprintf(w.W, "\tmov $jump_register, %%r%cx\n", reg+'a')
	return err
}

func (w Linux64AssemblyWriter) JumpValue(reg int) error {
	_, err := fmt.Fprintf(w.W, "\txor %%r%cx, %%r%cx\n\tmov jump_register, %%%cl\n", reg+'a', reg+'a', reg+'a')
	return err
}

func (w Linux64AssemblyWriter) Advance(reg, offset int) error {
	if reg != 0 {
		log.Panicln("internal compiler error: unsupported Advance register for Linux64AssemblyWriter:", reg)
	}
	_, err := fmt.Fprintf(w.W, "\tadd $%d, %%rax\n\tcall ensure\n", offset)
	return err
}

func (w Linux64AssemblyWriter) NandBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tand %%r%cx, %%r%cx\n\txor $1, %%r%cx\n", src+'a', dest+'a', dest+'a')
	return err
}
