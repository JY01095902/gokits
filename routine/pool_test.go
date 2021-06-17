package routine

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func runNewPool(cnt int, list []string, prefix string) {
	tasks := make(chan Task)
	pool := NewPool(10, tasks)
	pool.Run()

	for i := 0; i < cnt; i++ {
		n := i
		job := func() error {
			val := prefix + time.Now().String()
			list[n] = val

			return nil
		}

		tasks <- Task{
			Id:  strconv.Itoa(n + 1),
			Job: job,
		}
	}
	close(tasks)

	pool.Wait()
}

func BenchmarkNewPool(b *testing.B) {
	cnt := 100
	list := make([]string, cnt)
	prefix := "TEST"
	for n := 0; n < b.N; n++ {
		runNewPool(cnt, list, prefix)
	}
}

func TestNewPool(t *testing.T) {
	cnt := 100
	list := make([]string, cnt)
	prefix := "TEST"
	runNewPool(cnt, list, prefix)

	Convey("测试线程池", t, func() {
		So(len(list), ShouldEqual, cnt)
		for i := range list {
			So(list[i], ShouldStartWith, prefix)
		}
	})
}

func runNewPoolWithResultHandler(cnt int, list []string, errResults *[]TaskResult, prefix, msgPrefix, errMsg string) {
	tasks := make(chan Task)
	pool := NewPoolWithResultHandler(10, tasks, func(result TaskResult) {
		if result.Status == "FAILED" {
			*errResults = append(*errResults, result)
		}
	})
	pool.Run()

	for i := 0; i < cnt; i++ {
		n := i
		job := func() error {
			val := prefix + time.Now().String()
			list[n] = val

			if n%2 == 0 {
				return errors.New(errMsg)
			}

			return nil
		}

		tasks <- Task{
			Id:      strconv.Itoa(n + 1),
			Job:     job,
			Message: msgPrefix + " " + strconv.Itoa(n),
		}
	}
	close(tasks)

	pool.Wait()
}

func BenchmarkNewPoolWithResultHandler(b *testing.B) {
	cnt := 100
	list := make([]string, cnt)
	errResults := []TaskResult{}
	prefix := "TEST"
	msgPrefix := "number is"
	errMsg := "can not handle even number"
	for n := 0; n < b.N; n++ {
		runNewPoolWithResultHandler(cnt, list, &errResults, prefix, msgPrefix, errMsg)
	}
}

func TestNewPoolWithResultHandler(t *testing.T) {
	cnt := 100
	list := make([]string, cnt)
	errResults := []TaskResult{}
	prefix := "TEST"
	msgPrefix := "number is"
	errMsg := "can not handle even number"
	runNewPoolWithResultHandler(cnt, list, &errResults, prefix, msgPrefix, errMsg)

	Convey("测试带结果处理的线程池", t, func() {
		So(len(list), ShouldEqual, cnt)
		for i := range list {
			So(strings.HasPrefix(list[i], prefix), ShouldBeTrue)
		}
		So(len(errResults), ShouldEqual, cnt/2)
		for i := range errResults {
			So(errResults[i].Message, ShouldStartWith, msgPrefix)
			So(errResults[i].Error.Error(), ShouldEqual, errMsg)
		}
	})
}
