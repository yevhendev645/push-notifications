package taskmgr

type Entity string

// TaskType defines a valid task type
type TaskType string

// TaskProcessor defines the interface a task handler must implement
type TaskProcessor interface {
	Process() error
	TaskType() TaskType
	Unmarshal([]byte) (TaskProcessor, error)
	GetRetries() int
	SetRetries(retries int)
	GetAttempts() int
	SetAttempts(retries int)
}

// HasRetries defines number of retries a task has left
type HasRetries struct {
	RetriesLeft int `json:"retriesLeft"`
	Attempts    int `json:"attempts"`
}

type InfoForLogging struct {
	Entity
	UserID string
}

// GetRetries gets number of retries left
func (r *HasRetries) GetRetries() int {
	return r.RetriesLeft
}

// SetRetries sets number of retries
func (r *HasRetries) SetRetries(left int) {
	r.RetriesLeft = left
}

// GetAttempts gets number of retries left
func (r *HasRetries) GetAttempts() int {
	return r.Attempts
}

// SetAttempts sets number of retries
func (r *HasRetries) SetAttempts(att int) {
	r.Attempts = att
}
