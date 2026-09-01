package jobs

type Job struct {
	Sport     Sport
	Operation Operation
	Date      string
}
type Sport string

const SportMLB Sport = "mlb"

type Operation string

const OperationGetSchedule Operation = "get_schedule"
