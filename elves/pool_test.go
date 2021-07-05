/*
	go test -bench=. -memprofile mem.prof
	go tool pprof -http :9402 mem.prof
	go test -bench=. -test.benchmem
*/

package elves

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPool(b *testing.B) {
	run := func() {
		pool, _ := NewPool(10)
		var wg sync.WaitGroup
		cnt := 20
		wg.Add(cnt)
		for i := 1; i <= cnt; i++ {
			go func(i int) {
				pool.Execute(func() {
					defer wg.Done()

					log.Printf("id: %d, time: %v", i, time.Now())
				})
			}(i)
		}

		wg.Wait()
		pool.Destroy()
	}

	for n := 0; n < b.N; n++ {
		run()
	}
}

func TestNewPool(t *testing.T) {
	tests := []struct {
		given int
		want  error
	}{
		{given: 10, want: nil},
		{given: 0, want: ErrInvalidCapacity},
	}

	for _, test := range tests {
		_, err := NewPool(test.given)

		assert.Equal(t, test.want, err)
	}
}

func TestMaxWorkerCount(t *testing.T) {
	tests := []struct {
		givenMaxCap    int
		givenWorkerCnt int
		wantWorkerCnt  int
	}{
		{givenMaxCap: 10, givenWorkerCnt: 5, wantWorkerCnt: 5},
		{givenMaxCap: 10, givenWorkerCnt: 20, wantWorkerCnt: 10},
		{givenMaxCap: 10, givenWorkerCnt: 0, wantWorkerCnt: 0},
		{givenMaxCap: 10, givenWorkerCnt: 10, wantWorkerCnt: 10},
	}

	for _, test := range tests {
		pool, _ := NewPool(test.givenMaxCap)

		var wg sync.WaitGroup
		wg.Add(test.givenWorkerCnt)
		for i := 0; i < test.givenWorkerCnt; i++ {
			go func(i int) {
				pool.Execute(func() {
					defer wg.Done()

					time.Sleep(1000 * time.Millisecond)
				})
			}(i)
		}

		wg.Wait()

		assert.LessOrEqual(t, len(pool.workers), test.wantWorkerCnt, "workers count should be less than or equal to pool capacity")

		pool.Destroy()
	}
}
