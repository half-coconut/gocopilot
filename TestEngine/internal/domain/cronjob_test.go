package domain

import (
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := "*/15 * * * *"
	a := NextTimeV1(c)
	t.Log(a.Format(time.DateTime))

}

func TestTime(t *testing.T) {
	a := time.UnixMilli(1744360220933)
	t.Log(a.Format(time.DateTime))
}

func TestTricker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		t.Log("执行了。。。")
	}
}
