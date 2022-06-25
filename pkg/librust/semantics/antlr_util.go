package semantics

import (
	"fmt"
	"strings"
)

type MessageType int

const (
	INFO MessageType = iota
	WARN
	ERROR
)

func (mt MessageType) String() string {
	switch mt {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	}
	return "???"
}

type Message struct {
	Type MessageType
	Desc string
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s", m.Type.String(), m.Desc)
}

type PatternReturn struct {
	Name     string
	Type     string
	Volatile bool
	Size     int
}

type TypeDef interface {
	String() string
}

type IDAttr struct {
	Name      string
	TypeParam TypeDef
}

func (id IDAttr) String() string {
	return fmt.Sprintf("{Name: %s; %s}", id.Name, id.TypeParam)
}

type ValueAttr struct {
	TypeDef
	Volatile bool
	BaseType TypeDef
}

func (val ValueAttr) String() string {
	return fmt.Sprintf("{VALUE; Volatile: %t; BaseType: %s}", val.Volatile, val.BaseType)
}

type RefAttr struct {
	TypeDef
	Volatile bool
	BaseType TypeDef
}

func (rf RefAttr) String() string {
	return fmt.Sprintf("{REF; Volatile: %t; BaseType: %s}", rf.Volatile, rf.BaseType)
}

type TypeAttr struct {
	TypeDef
	BaseType string
}

func (tp TypeAttr) String() string {
	return fmt.Sprintf("{TYPE; %s}", tp.BaseType)
}

type TypedefAttr struct {
	TypeDef
	Type TypeDef
}

func (tp TypedefAttr) String() string {
	return fmt.Sprintf("{TYPEDEF: %s}", tp.Type)
}

type PointerAttr struct {
	TypeDef
	Volatile bool
	Type     TypeDef
}

func (p PointerAttr) String() string {
	return fmt.Sprintf("{POINTER; Volatile: %t; %s}", p.Volatile, p.Type)
}

type ArrayAttr struct {
	TypeDef
	Type      TypeDef
	NumOfElem int
}

func (arr ArrayAttr) String() string {
	return fmt.Sprintf("{ARRAY; NumOfElem: %d; Type: %s}", arr.NumOfElem, arr.Type)
}

type FuncAttr struct {
	TypeDef
	ReturnType TypeDef
	CallParam  []TypeDef
}

func (fn FuncAttr) String() string {
	callParamStrings := make([]string, len(fn.CallParam))
	for i, param := range fn.CallParam {
		callParamStrings[i] = param.String()
	}
	return fmt.Sprintf("{FUNCTION; ReturnType: %s; CallParam: %s}", fn.ReturnType, strings.Join(callParamStrings, "; "))
}

type NamespaceAttr struct {
	TypeDef
	External string
}

func (nm NamespaceAttr) String() string {
	return fmt.Sprintf("{NAMESPACE; External: %s}", nm.External)
}

type ReturnAttr struct {
	TypeDef
	Type TypeDef
}

func (ra ReturnAttr) String() string {
	return fmt.Sprintf("{TO RETURN: %s}", ra.Type)
}

type Scope map[string]TypeDef

func GetToType(t TypeDef) TypeDef {
	switch attr := t.(type) {
	case ValueAttr:
		return GetToType(attr.BaseType)
	case TypedefAttr:
		return GetToType(attr.Type)
	case FuncAttr:
		return GetToType(attr.ReturnType)
	case ReturnAttr:
		return GetToType(attr.Type)
	case IDAttr:
		return GetToType(attr.TypeParam)
	default:
		return attr
	}
}
