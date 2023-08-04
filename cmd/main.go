package main

import (
	"context"
	"fmt"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	//example1()
	//example2()
	//example3()
	example4()

	//ex_with_rx()
}

func goroutine_example() {
	sum := func(s []int, c chan rxgo.Item) {
		_sum := 0
		for _, v := range s {
			_sum += v
		}
		c <- rxgo.Item{V: _sum}
	}
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan rxgo.Item)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y)
}

func ex_with_rx() {
	sum := func(s []int, c chan rxgo.Item) {
		_sum := 0
		for _, v := range s {
			_sum += v
		}
		c <- rxgo.Item{V: _sum}
	}
	s := []int{7, 2, 8, -9, 4, 0}
	ch := make(chan rxgo.Item)
	go sum(s[:len(s)/2], ch)
	go sum(s[len(s)/2:], ch)

	observable := rxgo.FromChannel(ch)
	fmt.Println(fmt.Sprintf("observable: %#v", observable))
	for x := range observable.Observe() {
		if x.Error() {
			fmt.Println(x.E)
			break
		}
		fmt.Println(fmt.Sprintf("%#v", x))
	}
}

func example1() {
	observable := rxgo.Just("Hello, World!")()
	ch := observable.Observe()
	fmt.Println(fmt.Sprintf("%#v", ch))
	item := <-ch
	fmt.Println(fmt.Sprintf("%#v", item))
}

func example2() {
	observable := rxgo.Just("Hello, World!")()

	// `ForEach` is non-blocking. Yet, it returns a notification channel that will be closed once the Observable complets.
	// Hence, to make the previous code blocking, we simply need to use `<-`:
	<-observable.ForEach(func(v interface{}) {
		fmt.Println(fmt.Sprintf("next: %#v", v))
	}, func(err error) {
		fmt.Println(fmt.Sprintf("error: %#v", err))
	}, func() {
		fmt.Println("observable is closed.")
	})
}

func example3() {
	producer := func(ch chan rxgo.Item) {
		for i := range []int{1, 2, 3, 4} {
			fmt.Printf("insert %v\n", i)
			ch <- rxgo.Item{V: i}
			fmt.Printf("insert %v <- done\n", i)
		}
		close(ch)
	}

	ch := make(chan rxgo.Item)
	go producer(ch)

	observable := rxgo.FromChannel(ch)
	// `Observe()` returns a `<- chan Item`
	for item := range observable.Observe() {
		if item.Error() {
			fmt.Println(item.E)
			break
		}
		fmt.Println(fmt.Sprintf("%#v", item))
	}
}

func example4() {
	observable := rxgo.Defer([]rxgo.Producer{
		func(_ context.Context, ch chan<- rxgo.Item) {
			for i := 0; i < 3; i++ {
				ch <- rxgo.Of(i)
			}
		},
	})

	for item := range observable.Observe() {
		fmt.Println("first", item.V)
	}
	for item := range observable.Observe() {
		fmt.Println("second", item.V)
	}
}
