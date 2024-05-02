package data_processors

type TasksFromContest struct {
	ContestTitle string
	Tasks        []string
}

type UnsolvedData struct {
	total    int
	unsolved []*TasksFromContest
	// maybe we will need more data
}

type Value struct {
	Value string
	Color string
}

type UserValues struct {
	FirstName  string
	SecondName string
	FullName   string
	Values     []*Value
}

type CriterionTitle struct {
	Title    string
	EjudgeId int
}

type Stats struct {
	Count   int      `json:"Count"`
	Color   string   `json:"Color"`
	Peoples []string `json:"Peoples"`
}
