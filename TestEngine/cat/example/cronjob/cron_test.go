package cronjob

import (
	cron "github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestCronExpression(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	//expr.AddJob("@every 1s", myJob{})
	expr.AddFunc("@every 3s", func() {
		t.Log("开始长任务了")
		time.Sleep(time.Second * 12)
		t.Log("结束长任务了")
	})
	expr.Start()
	time.Sleep(time.Second * 10)
	stop := expr.Stop()
	t.Log("已经发出停止信号")
	<-stop.Done()
	t.Log("彻底结束")
}

type myJob struct {
}

func (m myJob) Run() {
	log.Println("运行了")
}
