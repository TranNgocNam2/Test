package status

type Verification int16

const (
	Pending Verification = iota
	Verified
	Rejected
	Cancelled
)
