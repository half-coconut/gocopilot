package time

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	tm := time.NewTicker(time.Second)
	defer tm.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Log("超时了，或者被取消了")
			//return
			goto end
		case now := <-tm.C:
			t.Log(now.Unix())
		}
	}
end:
}

func TestTimer(t *testing.T) {
	tm := time.NewTimer(time.Second)
	defer tm.Stop()
	go func() {
		for now := range tm.C {
			t.Log(now.Unix())
		}
	}()
	time.Sleep(time.Second * 10)
}
