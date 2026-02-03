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

func DoubleThen[T any, P Promise[T, P]](p P, fn func(T) (T, error)) P {
	return p.Then(fn).Then(fn)
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

func (p *promise[T]) Finally(fn func()) *promise[T] {
	next := &promise[T]{done: make(chan struct{})}
	go func() {
		<-p.done
		fn()
		next.value = p.value
		next.err = p.err
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

func Parallel[T any](promises ...*promise[T]) *promise[[]T] {
	p := &promise[[]T]{done: make(chan struct{})}
	go func() {
		results := make([]T, len(promises))
		for i, pr := range promises {
			v, err := pr.Await()
			if err != nil {
				p.err = err
				close(p.done)
				return
			}
			results[i] = v
		}
		p.value = results
		close(p.done)
	}()
	return p
}

func main() {
	p := NewPromise(func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 10, nil
	})

	doubled := DoubleThen(p, func(v int) (int, error) {
		fmt.Println("double:", v)
		return v * 2, nil
	})

	result, err := doubled.
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
			return -1, err
		}).
		Finally(func() {
			fmt.Println("finally called!")
		}).
		Await()

	fmt.Println("final result:", result, "error:", err)

	// Parallel example
	p1 := NewPromise(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("p1 done")
		return 100, nil
	})
	p2 := NewPromise(func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("p2 done")
		return 200, nil
	})
	p3 := NewPromise(func() (int, error) {
		time.Sleep(150 * time.Millisecond)
		fmt.Println("p3 done")
		return 300, nil
	})

	results, err := Parallel(p1, p2, p3).Await()
	fmt.Println("parallel results:", results, "error:", err)
}
