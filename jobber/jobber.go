package jobber

type JobSettings struct {
	Name   string `json:"name"`
	Worker string `json:"worker"`
	Repeat string `json:"name"`
}

type Jobber struct {
}
