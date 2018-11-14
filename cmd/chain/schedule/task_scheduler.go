package schedule

import (
    "container/list"
    "github.com/ssor/zlog"
    "time"
)

var (
    scheduleLogger = zlog.New("chain", "stream", "schedule")
)

func NewScheduler(taskExecutor func(height int) bool) *Scheduler {
    return &Scheduler{
        blocks:       list.New(),
        taskExecutor: taskExecutor,
    }
}

type Scheduler struct {
    latestBlock  int
    blocks       *list.List
    taskExecutor func(height int) bool
}

func (scheduler *Scheduler) Run() {
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for {
            <-ticker.C
            //scheduleLogger.Info("scheduler running ...")
            if scheduler.blocks.Len() <= 0 {
                continue
            }

            if scheduler.taskExecutor == nil {
                continue
            }

            head := scheduler.blocks.Front()
            ok := scheduler.taskExecutor(head.Value.(int))
            if ok {
                scheduler.blocks.Remove(head)
            }
        }
    }()
}

func (scheduler *Scheduler) EventHandler(height int) {
    if height <= scheduler.latestBlock {
        return
    }
    for i := scheduler.latestBlock + 1; i <= height; i++ {
        scheduler.blocks.PushBack(i)
        scheduler.latestBlock = i
    }
}
