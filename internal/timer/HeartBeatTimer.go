package timer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// todo 之后再做
func HeartBeatTimer() {
	_, err := Timer.AddFunc("*/30 * * * * *", func() {
		fmt.Println("每30s执行一次的任务:", time.Now())
	})
	if err != nil {
		logrus.Error("定时任务(%v)添加失败:", "HeartBeatTimer", err)
		return
	}
}
