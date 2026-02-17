const (
	nodes      = 5
	metrics    = 1000
	iterations = 200000
)
type Metric struct {
	bw   int
	loss int
	q    int
}
type Result struct {
	node int
	load int
}
func collect() Metric {
	return Metric{
		bw:   rand.Intn(1000),
		loss: rand.Intn(100),
		q:    rand.Intn(500),
	}
}
func analyze(m Metric) int {
	score := 0
	for i := 0; i < 50; i++ {
		score += (m.bw + m.q - m.loss)
	}
	return score
}
func nodeWorker(id int, out chan<- Result, wg *sync.WaitGroup) {
	for i := 0; i < iterations; i++ {
		m := collect()
		load := analyze(m)
		out <- Result{node: id, load: load}
	}
	wg.Done()
}
func collector(in <-chan Result, done chan<- bool) {
	total := 0
	for r := range in {
		total += r.load
	}
	fmt.Println("Central Aggregate:", total)
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
	ch := make(chan Result, 1000)
	done := make(chan bool)
	var wg sync.WaitGroup
	go collector(ch, done)
	for i := 0; i < nodes; i++ {
		wg.Add(1)
		go nodeWorker(i, ch, &wg)
	}
	wg.Wait()
	close(ch)
	cpuUsage(start)
}
