package main

import "fmt"

// 自己参照型制約を使った汎用 Sum 関数
type Adder[T Adder[T]] interface {
	Add(T) T
}

type Int int

func (i Int) Add(other Int) Int {
	return i + other
}

func Sum[T Adder[T]](values []T) T {
	var result T
	for _, v := range values {
		result = result.Add(v)
	}
	return result
}

func main() {
	nums := []Int{1, 2, 3, 4, 5}
	fmt.Println("Sum:", Sum(nums))
}
