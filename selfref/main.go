package main

import "fmt"

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
