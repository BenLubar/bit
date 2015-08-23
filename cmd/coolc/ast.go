package main

import (
	"go/token"
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
}

type VarDecl struct {
	Name ID
	Type TYPE
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
}

type VarFeature struct {
	VarDecl
	Value Expr
}

type BlockFeature struct {
	Expr
}

type Expr interface{}

type AssignExpr struct {
	Left  ID
	Right Expr
}

type NotExpr struct {
	Right Expr
}

type NegativeExpr struct {
	Right Expr
}

type IfExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

type WhileExpr struct {
	Condition Expr
	Do        Expr
}

type LessThanOrEqualExpr struct {
	Left  Expr
	Right Expr
}

type LessThanExpr struct {
	Left  Expr
	Right Expr
}

type EqualEqualExpr struct {
	Left  Expr
	Right Expr
}

type MultiplyExpr struct {
	Left  Expr
	Right Expr
}

type DivideExpr struct {
	Left  Expr
	Right Expr
}

type AddExpr struct {
	Left  Expr
	Right Expr
}

type SubtractExpr struct {
	Left  Expr
	Right Expr
}

type MatchExpr struct {
	Left  Expr
	Cases *Cases
}

type CallExpr struct {
	Left Expr
	Name ID
	Args []Expr
}

type SelfCallExpr struct {
	Super bool
	Name  ID
	Args  []Expr
}

type NewExpr struct {
	Type TYPE
	Args []Expr
}

type VarExpr struct {
	Var  *VarFeature
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

type Cases struct {
	Cases []*Case
	Null  Expr
}

type Case struct {
	Name ID
	Type TYPE
	Body Expr
}
