package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	nodes      = 5
	iterations = 200000
)

type Metric struct {
	bw   int
	loss int
	q    int
}

type Packet struct {
	node   int
	metric Metric
}

func collect() Metric {
	return Metric{
		bw:   rand.Intn(1000),
		loss: rand.Intn(100),
		q:    rand.Intn(500),
	}
}

func nodeWorker(id int, out chan<- Packet, wg *sync.WaitGroup) {
	for i := 0; i < iterations; i++ {
		m := collect()
		out <- Packet{node: id, metric: m}
	}
	wg.Done()
}

func correlate(p Packet) int {
	score := 0
	for i := 0; i < 50; i++ {
		score += (p.metric.bw + p.metric.q - p.metric.loss)
	}
	return score
}

func correlationEngine(in <-chan Packet, done chan<- bool) {
	total := 0
	for p := range in {
		total += correlate(p)
	}
	fmt.Println("Correlation Aggregate:", total)
	done <- true
}

func cpuUsage(start time.Time) {
	elapsed := time.Since(start)
	fmt.Println("Elapsed:", elapsed)
	fmt.Println("Goroutines:", runtime.NumGoroutine())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	ch := make(chan Packet, 1000)
	done := make(chan bool)

	var wg sync.WaitGroup

	go correlationEngine(ch, done)

	for i := 0; i < nodes; i++ {
		wg.Add(1)
		go nodeWorker(i, ch, &wg)
	}

	wg.Wait()
	close(ch)

	<-done

	cpuUsage(start)
}
