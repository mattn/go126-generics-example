package main

import (
	"errors"
	"fmt"
	"time"
)

type Promise[T any, P Promise[T, P]] interface {
	Then(func(T) (T, error)) P
	Catch(func(error) (T, error)) P
	Await() (T, error)
}

type promise[T any] struct {
	value T
	err   error
	done  chan struct{}
}

func (p *promise[T]) Then(fn func(T) (T, error)) *promise[T] {
	next := &promise[T]{done: make(chan struct{})}
	go func() {
		<-p.done
		if p.err != nil {
			next.err = p.err
		} else {
			next.value, next.err = fn(p.value)
		}
		close(next.done)
	}()
	return next
}

func (p *promise[T]) Catch(fn func(error) (T, error)) *promise[T] {
	next := &promise[T]{done: make(chan struct{})}
	go func() {
		<-p.done
		if p.err != nil {
			next.value, next.err = fn(p.err)
		} else {
			next.value = p.value
		}
		close(next.done)
	}()
	return next
}

func (p *promise[T]) Await() (T, error) {
	<-p.done
	return p.value, p.err
}

func NewPromise[T any](fn func() (T, error)) *promise[T] {
	p := &promise[T]{done: make(chan struct{})}
	go func() {
		p.value, p.err = fn()
		close(p.done)
	}()
	return p
}

func main() {
	p := NewPromise(func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 10, nil
	})

	result, err := p.
		Then(func(v int) (int, error) {
			fmt.Println("step1:", v)
			return v * 2, nil
		}).
		Then(func(v int) (int, error) {
			fmt.Println("step2:", v)
			return 0, errors.New("something went wrong")
		}).
		Then(func(v int) (int, error) {
			fmt.Println("step3:", v)
			return v + 100, nil
		}).
		Catch(func(err error) (int, error) {
			fmt.Println("caught error:", err)
			return -1, nil
		}).
		Await()

	fmt.Println("final result:", result, "error:", err)
}
