package pipeline

import "time"

type Step struct {
    work         Work
    durationLeft time.Duration
}

func newStep(w Work) *Step {
    step := &Step{
        durationLeft: 0,
        work:         w,
    }
    return step
}

func (step *Step) timeTick(d time.Duration) bool {
    step.durationLeft -= d
    if step.durationLeft <= 0 {
        return true
    }
    return false
}

func (step *Step) Postpone(d time.Duration) {
    step.durationLeft = d
}

func (step *Step) Do(request *Request) *Result {
    return step.work.Exec(request)
}
