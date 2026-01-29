package main

import "fmt"

// Go 1.26: 自己参照型制約が可能に
// T が自分自身と同じ型を返す Clone メソッドを持つことを要求
type Cloner[T Cloner[T]] interface {
	Clone() T
}

type Person struct {
	Name string
}

func (p Person) Clone() Person {
	return Person{Name: p.Name}
}

func main() {
	p := Person{Name: "Alice"}
	c := p.Clone()
	fmt.Printf("original: %+v, cloned: %+v\n", p, c)
}
