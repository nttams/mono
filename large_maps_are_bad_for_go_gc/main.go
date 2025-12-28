package main

import (
	"fmt"
	"runtime"
	"time"
)

/*
go run ./...
n=1000000 | avg GC pause=10.400ms
n=10000000 | avg GC pause=101.300ms
n=20000000 | avg GC pause=314.600ms
*/

func main() {
	run(1_000_000)
	run(10_000_000)
	run(20_000_000)
}

func run(n int) {
	routes := make(map[string]string, n)

	for i := range n {
		routes[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value-%d", i)
	}

	const runs = 10
	var totalPause time.Duration
	for range runs {
		start := time.Now()
		runtime.GC()
		pause := time.Since(start)
		totalPause += pause
	}

	avgMs := float64(totalPause.Milliseconds()) / float64(runs)
	fmt.Printf("n=%d | avg GC pause=%.3fms\n", n, avgMs)

	_ = routes["key-0"] // prevent the map from being garbage collected
}
