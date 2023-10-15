package task

type Status string

const (
	StatusPending  Status = "pending"
	StatusValid    Status = "valid"
	StatusUnknown  Status = "unknown"
	StatusExpired  Status = "expired"
	StatusExpiring Status = "expiring"
)
