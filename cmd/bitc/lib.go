package main

const assemblyLib = `
.data
start_of_heap:
	.quad 0
end_of_heap:
	.quad 0

.text
.globl var
var:
	mov start_of_heap, %rbx
	test %rbx, %rbx
	jnz .L_var_ready
	push %rax
	mov $_end, %rdi
	mov $0xc, %rax
	syscall
	mov %rax, %rbx
	mov %rax, start_of_heap
	mov %rax, end_of_heap
	pop %rax
.L_var_ready:
	add %rbx, %rax
	call ensure
	ret

.text
.globl ensure
ensure:
	mov end_of_heap, %rdi
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
	mov %rax, end_of_heap
	pop %rcx
	pop %rdi
	xor %rax, %rax
	sub %rdi, %rcx
	rep stosb
	pop %rax
.L_ensure_ret:
	ret

.data
print_bits:
	.byte 0
print_accum:
	.byte 0

.text
.globl print
print:
	xor %rbx, %rbx
	xor %rcx, %rcx
	mov print_bits, %bl
	mov print_accum, %cl
	shl $1, %rcx
	inc %rbx
	or %rax, %rcx
	mov %cl, print_accum
	cmp $8, %rbx
	je .L_print_out
	mov %bl, print_bits
	ret
.L_print_out:
	mov $1, %rax
	mov $1, %rdi
	mov $print_accum, %rsi
	mov $1, %rdx
	syscall
	movb $0, print_bits
	movb $0, print_accum
	ret

.data
read_bits:
	.byte 0
read_accum:
	.byte 0

.text
.globl read
read:
	xor %rbx, %rbx
	mov read_bits, %bl
	test %rbx, %rbx
	jz .L_read_in
	dec %rbx
	mov %bl, read_bits
	mov read_accum, %bl
	mov %rbx, %rax
	and $1, %rax
	shr $1, %rbx
	mov %bl, read_accum
	ret
.L_read_in:
	push %rax
	mov $0, %rax
	mov $0, %rdi
	mov $read_accum, %rsi
	mov $1, %rdx
	syscall
	test %rax, %rax
	jz .L_read_fail
	pop %rax
	movb $8, read_bits
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
	call print
	xor %rax, %rax
	call print
	xor %rax, %rax
	call print
	xor %rax, %rax
	call print
	xor %rax, %rax
	call print
	xor %rax, %rax
	call print
	xor %rax, %rax
	call print
	mov $60, %rax
	mov $0, %rdi
	syscall

.data
.globl jump_register
jump_register:
	.byte 0
`
