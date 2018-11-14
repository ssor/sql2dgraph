package pipeline

// Result is returned by a step to dispatch data to the next step or stage
type Result struct {
    Error error
    // dispatch any type
    Data interface{}
    // dispatch key value pairs
    KeyVal map[string]interface{}
}

// Request is the result dispatched in a previous step.
type Request struct {
    Data   interface{}
    KeyVal map[string]interface{}
}

// Work is the unit of work which can be concurrently or sequentially staged with other steps
type Work interface {
    // Exec is invoked by the pipeline when it is run
    Exec(*Request) *Result
    // Cancel is invoked by the pipeline when one of the concurrent steps set Result{Error:err}
    //Cancel() error
    WorkDescription() string
}
