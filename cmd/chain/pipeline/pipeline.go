package pipeline

import (
    "container/list"
    "github.com/ssor/zlog"
    "time"
)

var (
    pipelineLogger = zlog.New("chain", "stream", "pipeline")

    defaultWorkCheckDuration = time.Second
    defaultPostponeDuration  = 5 * time.Second
)

// Pipeline is a sequence of stages
type Pipeline struct {
    Name              string `json:"name"`
    Steps             *list.List
    chanNewStep       chan *Step
    workCheckDuration time.Duration
}

// New returns a new pipeline
// 	name of the pipeline
// 	outBufferLen is the size of the output buffered channel
func New(name string) *Pipeline {
    return newPipeline(name)
}

func newPipeline(name string) *Pipeline {
    p := &Pipeline{
        Name:              name,
        Steps:             list.New(),
        chanNewStep:       make(chan *Step, 1024),
        workCheckDuration: defaultWorkCheckDuration,
    }
    return p
}

func (pipeline *Pipeline) WorkCount() int {
    return pipeline.Steps.Len()
}

func (pipeline *Pipeline) NewWork(w Work) {
    pipeline.chanNewStep <- newStep(w)
}

func (pipeline *Pipeline) doWork() {
    if pipeline.Steps.Len() <= 0 {
        return
    }

    ele := pipeline.Steps.Front()
    step := ele.Value.(*Step)
    if step.timeTick(time.Second) {
        result := step.Do(nil)
        if result.Error == nil {
            pipeline.Steps.Remove(ele)
            pipelineLogger.Successf("pipeline work ** %s ** done", step.work.WorkDescription())
        } else {
            pipelineLogger.Failedf("pipeline work ** %s ** failed for: %s", step.work.WorkDescription(), result.Error)
            step.Postpone(defaultPostponeDuration)
        }
    }
}

func (pipeline *Pipeline) Run() {
    ticker := time.NewTicker(pipeline.workCheckDuration)
    go func() {
        for {
            select {
            case step := <-pipeline.chanNewStep:
                pipeline.Steps.PushBack(step)
            case <-ticker.C:

                pipeline.doWork()
            }
        }
    }()
}
