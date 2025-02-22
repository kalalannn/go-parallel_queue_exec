package messages

const (
	ExecuteStart = "START: %s\n"
	ExecuteDone  = "DONE: %s\n"
)

const (
	ScheduleTag = "schedule"
	NextTag     = "next"
	StartTag    = "start"
	DoneTag     = "done"
)

const (
	WaitForWorkers       = "Waiting for workers to finish tasks..."
	WorkersDone          = "Workers done, count processed: %d\n"
	ExecutorDoneShutdown = "Executor done, shutting down..."
)
