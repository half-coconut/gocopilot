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
	a := time.UnixMilli(1744275240000)
	t.Log(a.Format(time.DateTime))
}
