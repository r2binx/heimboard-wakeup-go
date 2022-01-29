package schedule

type Schedule struct {
	Time   int64  `json:"time"`
	Action string `json:"action"`
}

func New(time int64, action string) *Schedule {
	return &Schedule{
		Time:   time,
		Action: action,
	}
}
