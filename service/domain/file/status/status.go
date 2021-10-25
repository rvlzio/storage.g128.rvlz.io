package status

type Status struct {
	status string
}

var (
	Instantiated = Status{"instantiated"}
	Verifying    = Status{"verifying"}
)
