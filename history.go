package cronroutine

import "time"

type History struct {
	jobID string
	ranAt time.Time
	err   error
}

func (h *History) JobID() string    { return h.jobID }
func (h *History) RanAt() time.Time { return h.ranAt }
func (h *History) Error() error     { return h.err }
