package themetxt

import (
	"fmt"
	"strings"
)

type Property struct {
	name  string
	value interface{}
}

type Length interface {
	GetConvertFunc() func(val float64) float64
}

// 50
type AbsNum int

func (v AbsNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		return float64(v)
	}
}

// 50%
type RelNum int

func (v RelNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		return float64(v) / 100.0 * val
	}
}

// 50%-10
// rel: 50
// abs: 10
type CombinedNum struct {
	Rel int
	Abs int
	Op  CombinedNumOp
}

type CombinedNumOp int

const (
	CombinedNumAdd CombinedNumOp = iota
	CombinedNumSub
)

func (v CombinedNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		rel := float64(v.Rel) / 100.0 * val
		return rel - float64(v.Abs)
	}
}

type Component struct {
	Id       string
	Props    []*Property
	Children []*Component
}

func (c *Component) GetPropString(name string) (string, bool) {
	return getPropString(c.Props, name)
}

func (c *Component) GetPropLength(name string) (Length, bool) {
	return getPropLength(c.Props, name)
}

func (c *Component) GetPropBool(name string) (bool, bool) {
	return getPropBool(c.Props, name)
}

func (c *Component) Dump(indent int) {
	indentStr := strings.Repeat(" ", indent*4)
	fmt.Printf("%s+ %s {\n", indentStr, c.Id)

	for _, prop := range c.Props {
		fmt.Printf("%s    %s = %T %#v\n", indentStr, prop.name, prop.value, prop.value)
	}

	for _, child := range c.Children {
		child.Dump(indent + 1)
	}

	fmt.Printf("%s}\n", indentStr)
}

func getPropString(props []*Property, name string) (string, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(string), true
	}
	return "", false
}

func getPropBool(props []*Property, name string) (bool, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(bool), true
	}
	return false, false
}

func getPropLength(props []*Property, name string) (Length, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(Length), true
	}
	return nil, false
}

func getProp(props []*Property, name string) (interface{}, bool) {
	for _, prop := range props {
		if prop.name == name {
			return prop.value, true
		}
	}
	return nil, false
}

type Theme struct {
	Props      []*Property
	Components []*Component
}

func (t *Theme) GetPropString(name string) (string, bool) {
	return getPropString(t.Props, name)
}

func (t *Theme) Dump() {
	for _, prop := range t.Props {
		fmt.Printf("%s : %T %#v\n", prop.name, prop.value, prop.value)
	}
	for _, comp := range t.Components {
		comp.Dump(0)
	}
}

func ParseThemeFile(filename string) (*Theme, error) {
	v, err := ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return v.(*Theme), nil
}
