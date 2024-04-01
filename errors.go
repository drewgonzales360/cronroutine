package cronroutine

type ErrJobRunning struct{}

func (e ErrJobRunning) Error() string {
	return "job could not start job before deadline because another instance is still running"
}

type ErrPastStartingDeadline struct{}

func (e ErrPastStartingDeadline) Error() string {
	return "job could not start before starting deadline"
}

type ErrJobTimeout struct{}

func (e ErrJobTimeout) Error() string {
	return "job execution context deadline exceeded"
}
