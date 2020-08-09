package retry

import (
	"fmt"
	"time"
)

func ExampleRetryer() {
	noSleepingInTests := func(_ time.Duration) time.Duration { return 0 }

	fmt.Println("> Happy path")
	i := 0
	var err error
	for r := New(3, noSleepingInTests); r.Next(err); {
		i++
		fmt.Println("Attempt", i)
		err = nil
		if i < 2 {
			err = fmt.Errorf("connection error")
		}
	}
	fmt.Println("Err:", err)

	fmt.Println("> Unhappy path")
	i = 0
	for r := New(3, noSleepingInTests); r.Next(err); {
		i++
		fmt.Println("Attempt", i)
		err = fmt.Errorf("connection error")
	}
	fmt.Println("Err:", err)
	// Output: > Happy path
	// Attempt 1
	// Attempt 2
	// Err: <nil>
	// > Unhappy path
	// Attempt 1
	// Attempt 2
	// Attempt 3
	// Err: connection error
}

func ExampleExpDuration() {
	df := ExpDuration(time.Millisecond)
	for i := time.Duration(1); i <= 5; i++ {
		fmt.Println(df(i).Round(time.Millisecond))
	}
	// Output: 1ms
	// 3ms
	// 7ms
	// 20ms
	// 55ms
}
