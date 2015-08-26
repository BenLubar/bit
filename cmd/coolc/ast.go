package main

import (
	"go/token"

	"github.com/BenLubar/bit/bitgen"
)

type AST struct {
	FileSet *token.FileSet

	Classes []*ClassDecl
	main    *ClassDecl
}

type ID struct {
	Name   string
	Pos    token.Pos
	target interface{}
}

type TYPE struct {
	Name   string
	Pos    token.Pos
	target *ClassDecl
}

type ClassDecl struct {
	Name    TYPE
	Args    []*VarDecl
	Extends *ExtendsDecl
	Body    []Feature

	methods map[string]*MethodFeature
	size    uint
	depth   uint
}

type VarDecl struct {
	Name ID
	Type TYPE

	offset uint
	arg    uint
}

type ExtendsDecl struct {
	Type TYPE
	Args []Expr
}

type Feature interface{}

type MethodFeature struct {
	Override bool
	Name     ID
	Args     []*VarDecl
	Return   TYPE
	Body     Expr

	offset uint
}

type VarFeature struct {
	VarDecl
	Value Expr
}

type BlockFeature struct {
	Expr
}

type NativeFeature struct{}

type Expr interface{}

type AssignExpr struct {
	Left  ID
	Right Expr
}

type IfExpr struct {
	Pos       token.Pos
	Condition Expr
	Then      Expr
	Else      Expr
}

type WhileExpr struct {
	Pos       token.Pos
	Condition Expr
	Do        Expr
}

type MatchExpr struct {
	Left  Expr
	Cases []*Case
}

type CallExpr struct {
	Left Expr
	Name ID
	Args []Expr
}

type StaticCallExpr struct {
	Name ID
	Args []Expr
}

type NewExpr struct {
	Type TYPE
	Args []Expr
}

type VarExpr struct {
	VarFeature
	Expr Expr
}

type ChainExpr struct {
	Pre  Expr
	Expr Expr
}

type NullExpr struct {
}

type UnitExpr struct {
}

type NameExpr struct {
	Name ID
}

type IntegerExpr struct {
	N int32
}

type StringExpr struct {
	S string
}

type BooleanExpr struct {
	B bool
}

type ThisExpr struct {
}

type Case struct {
	Name ID
	Type TYPE
	Body Expr
}

type NativeExpr func(w *writer, start, end bitgen.Line)
