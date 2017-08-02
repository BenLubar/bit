package main

import (
	"fmt"
	"io"
	"log"

	"github.com/BenLubar/bit/bitnum"
)

const linux32Runtime = `
.data
.L_start_of_heap:
	.int 0
.L_end_of_heap:
	.int 0

.text
.globl var
var:
	mov .L_start_of_heap, %ebx
	test %ebx, %ebx
	jnz .L_var_ready
	push %eax
	xor %ebx, %ebx
	mov $0x2d, %eax
	int $0x80
	mov %eax, %ebx
	mov %eax, .L_start_of_heap
	mov %eax, .L_end_of_heap
	pop %eax
.L_var_ready:
	add %ebx, %eax
	call ensure
	ret

.text
.globl ensure
ensure:
	mov .L_end_of_heap, %edi
	cmp %edi, %eax
	jl .L_ensure_ret
	push %eax
	push %edi
	mov %eax, %ebx
	add $0x1001, %ebx
	and $0xfffff000, %ebx
	push %ebx
	mov $0x2d, %eax
	int $0x80
	mov %eax, .L_end_of_heap
	pop %ecx
	pop %edi
	xor %eax, %eax
	sub %edi, %ecx
	rep stosb
	pop %eax
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
	xor %ebx, %ebx
	xor %ecx, %ecx
	mov .L_print_bits, %bl
	mov .L_print_accum, %cl
	shl $1, %ecx
	inc %ebx
	or %eax, %ecx
	mov %cl, .L_print_accum
	cmp $8, %ebx
	je .L_print_out
	mov %bl, .L_print_bits
	ret
.L_print_out:
	mov $4, %eax
	mov $1, %ebx
	mov $.L_print_accum, %ecx
	mov $1, %edx
	int $0x80
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
	xor %ebx, %ebx
	mov .L_read_bits, %bl
	test %ebx, %ebx
	jz .L_read_in
	dec %ebx
	mov %bl, .L_read_bits
	mov .L_read_accum, %bl
	mov %ebx, %eax
	and $1, %eax
	shr $1, %ebx
	mov %bl, .L_read_accum
	ret
.L_read_in:
	push %eax
	mov $3, %eax
	mov $0, %ebx
	mov $.L_read_accum, %ecx
	mov $1, %edx
	int $0x80
	test %eax, %eax
	jz .L_read_fail
	pop %eax
	movb $8, .L_read_bits
	jmp read
.L_read_fail:
	pop %eax
	test %eax, %eax
	jz exit
	mov %eax, (%esp)
	ret

.text
.globl exit
exit:
	xor %eax, %eax
	mov .L_print_bits, %al
	test %eax, %eax
	jz .L_exit_exit
	mov $7, %ecx
	sub %eax, %ecx
	mov .L_print_accum, %al
	shl %cl, %eax
	mov %al, .L_print_accum
	movb $7, .L_print_bits
	xor %eax, %eax
	call print
.L_exit_exit:
	mov $1, %eax
	mov $0, %ebx
	int $0x80

.data
.globl jump_register
jump_register:
	.byte 0
`

type Linux32AssemblyWriter struct {
	W io.Writer
}

func (w Linux32AssemblyWriter) Runtime() error {
	_, err := io.WriteString(w.W, linux32Runtime)
	return err
}

func (w Linux32AssemblyWriter) DataSegment() error {
	_, err := io.WriteString(w.W, "\n.data\n")
	return err
}

func (w Linux32AssemblyWriter) TextSegment() error {
	_, err := io.WriteString(w.W, "\n.text\n")
	return err
}

func (w Linux32AssemblyWriter) DeclarePointer(n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "ptr_%s:\n\t.int 0\n", n.ShortString())
	return err
}

func (w Linux32AssemblyWriter) Start() error {
	_, err := io.WriteString(w.W, ".globl _start\n_start:\n")
	return err
}

func (w Linux32AssemblyWriter) DeclareLine(n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, ".L_l%s:\n", n.ShortString())
	return err
}

func (w Linux32AssemblyWriter) Goto(zero, one *bitnum.Number) error {
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

func (w Linux32AssemblyWriter) Read(eof *bitnum.Number) error {
	if eof != nil {
		if _, err := fmt.Fprintf(w.W, "\tmov $.L_l%s, %%eax\n", eof.ShortString()); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w.W, "\txor %eax, %eax\n"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w.W, "\tcall read\n\tmov %al, jump_register\n")
	return err
}

func (w Linux32AssemblyWriter) BitValue(register int, bit bool) error {
	if bit {
		_, err := fmt.Fprintf(w.W, "\tmov $1, %%e%cx\n", register+'a')
		return err

	}
	_, err := fmt.Fprintf(w.W, "\txor %%e%cx, %%e%cx\n", register+'a', register+'a')
	return err
}

func (w Linux32AssemblyWriter) Print() error {
	_, err := io.WriteString(w.W, "\tcall print\n")
	return err
}

func (w Linux32AssemblyWriter) SaveRegister(register int) error {
	_, err := fmt.Fprintf(w.W, "\tpush %%e%cx\n", register+'a')
	return err
}

func (w Linux32AssemblyWriter) LoadRegister(register int) error {
	_, err := fmt.Fprintf(w.W, "\tpop %%e%cx\n", register+'a')
	return err
}

func (w Linux32AssemblyWriter) WritePointer(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov %%e%cx, (%%e%cx)\n", src+'a', dest+'a')
	return err
}

func (w Linux32AssemblyWriter) WriteBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov %%%cl, (%%e%cx)\n", src+'a', dest+'a')
	return err
}

func (w Linux32AssemblyWriter) PointerAddress(reg int, n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "\tmov $ptr_%s, %%e%cx\n", n.ShortString(), reg+'a')
	return err
}

func (w Linux32AssemblyWriter) BitAddress(reg int, n *bitnum.Number) error {
	if reg != 0 {
		log.Panicln("internal compiler error: unsupported BitAddress register for Linux32AssemblyWriter:", reg)
	}
	_, err := fmt.Fprintf(w.W, "\tmov $%d, %%eax\n\tcall var\n", n.Uint64())
	return err
}

func (w Linux32AssemblyWriter) PointerValue(reg int, n *bitnum.Number) error {
	_, err := fmt.Fprintf(w.W, "\tmov ptr_%s, %%e%cx\n", n.ShortString(), reg+'a')
	return err
}

func (w Linux32AssemblyWriter) ReadBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tmov (%%e%cx), %%%cl\n\tand $1, %%e%cx\n", src+'a', dest+'a', dest+'a')
	return err
}

func (w Linux32AssemblyWriter) JumpAddress(reg int) error {
	_, err := fmt.Fprintf(w.W, "\tmov $jump_register, %%e%cx\n", reg+'a')
	return err
}

func (w Linux32AssemblyWriter) JumpValue(reg int) error {
	_, err := fmt.Fprintf(w.W, "\txor %%e%cx, %%e%cx\n\tmov jump_register, %%%cl\n", reg+'a', reg+'a', reg+'a')
	return err
}

func (w Linux32AssemblyWriter) Advance(reg, offset int) error {
	if reg != 0 {
		log.Panicln("internal compiler error: unsupported Advance register for Linux32AssemblyWriter:", reg)
	}
	_, err := fmt.Fprintf(w.W, "\tadd $%d, %%eax\n\tcall ensure\n", offset)
	return err
}

func (w Linux32AssemblyWriter) NandBit(dest, src int) error {
	_, err := fmt.Fprintf(w.W, "\tand %%e%cx, %%e%cx\n\txor $1, %%e%cx\n", src+'a', dest+'a', dest+'a')
	return err
}
