package main

const runtimeASM = `.data
.Lzero:
.ascii "ZERO\n"
.LzeroLen = . - .Lzero
.Lone:
.ascii "ONE\n"
.LoneLen = . - .Lone

.Lreadidx:
.byte 0
.Lreadcur:
.byte 0
.Lreadbuf:
.byte 0

.LpageSize = 0x1000
.LpageSizeMask = .LpageSize - 1
.LpageSizeInvMask = 0 - .LpageSize
.LpageSizeMinus4 = .LpageSize - 4
.LpageSizeMinus5 = .LpageSize - 5

.text
read:
.cfi_startproc
	push %eax
	push %ebx
	push %ecx
.LreadLoop:
	movl $3, %eax
	movl $0, %ebx
	xorl %ecx, %ecx
	leal .Lreadbuf, %ecx
	movl $1, %edx
	int $0x80
	test %eax, %eax
	jz exit
	movb .Lreadbuf, %al
	movb .Lreadcur, %bh
	movb .Lreadidx, %bl
	cmp $79, %eax
	je .LreadO
	cmp $90, %eax
	je .LreadZ
	cmp $9, %eax
	ja .LreadReset
	cmp $13, %eax
	ja .LreadLoop
	cmp $32, %eax
	je .LreadLoop
	cmp $69, %eax
	jne .LreadNotE
	test $0x0001, %ebx
	je .LreadNext
	test $0x0102, %ebx
	je .LreadDone
	jmp .LreadReset
.LreadO:
	test $0x0003, %ebx
	je .LreadDone
	movb $1, .Lreadcur
	movb $1, .Lreadidx
	jmp .LreadLoop
.LreadZ:
	movb $0, .Lreadcur
	movb $1, .Lreadidx
	jmp .LreadLoop
.LreadNotE:
	cmp $78, %eax
	jne .LreadMaybeR
	test $0x0101, %ebx
	je .LreadNext
	jmp .LreadReset
.LreadMaybeR:
	cmp $82, %eax
	jne .LreadReset
	test $0x0002, %ebx
	jne .LreadReset
	jmp .LreadNext
.LreadNext:
	addb $1, .Lreadidx
	jmp .LreadLoop
.LreadReset:
	movb $0, .Lreadidx
	jmp .LreadLoop
.LreadDone:
	xorl %edx, %edx
	movb .Lreadcur, %dl
	movb $0, .Lreadidx
	pop %ecx
	pop %ebx
	pop %eax
	ret
.cfi_endproc
print0:
.cfi_startproc
	push %ecx
	push %edx
	mov $.Lzero, %ecx
	mov $.LzeroLen, %edx
	jmp .Lprint
print1:
	push %ecx
	push %edx
	mov $.Lone, %ecx
	mov $.LoneLen, %edx
.Lprint:
	push %eax
	push %ebx
	mov $0x04, %eax
	mov $1, %ebx
	int $0x80
	pop %ebx
	pop %eax
	pop %edx
	pop %ecx
	ret
.cfi_endproc
exit:
.cfi_startproc
	movl $1, %eax
	movl $0, %ebx
	int $0x80
.cfi_endproc
incr:
.cfi_startproc
	push %ecx
	movl %eax, %ecx
	andl $.LpageSizeMask, %ecx
	addl %ebx, %ecx
	cmpl $.LpageSizeMinus5, %ecx
	jae .LincrDeref
	addl %ebx, %eax
	pop %ecx
	ret
.LincrDeref:
	andl $.LpageSizeInvMask, %eax
	addl $.LpageSizeMinus4, %eax
	subl $.LpageSizeMinus5, %ecx
	call alloc
	movl (%eax), %eax
	cmpl $.LpageSizeMinus5, %ecx
	jae .LincrDeref
	addl %ecx, %eax
	pop %ecx
	ret
.cfi_endproc
alloc:
.cfi_startproc
	cmpl $0, (%eax)
	jne .LallocDone
	push %eax
	push %ebx
	push %ecx
	movl %eax, %ecx
	movl $45, %eax
	movl $0, %ebx
	int $0x80
	addl $.LpageSizeMask, %eax
	andl $.LpageSizeInvMask, %eax
	movl %eax, (%ecx)
	movl %eax, %edi
	leal .LpageSize(%eax), %ebx
	movl $45, %eax
	int $0x80
	leal -6(%eax), %ecx
	subl %edi, %ecx
	movl $0, -4(%eax)
	movb $255, -5(%eax)
	movb $0, %al
	rep stosb
	pop %ecx
	pop %ebx
	pop %eax
.LallocDone:
	ret
.cfi_endproc

.data
zero:
.byte 0
one:
.byte 1
.align 4
`
