package task

type Task struct {
	Duration int    // millisecond (const memory size)
	ID       string // unique identifier
}

func NewTask(id string, duration int) Task {
	return Task{
		Duration: duration,
		ID:       id,
	}
}
