package data_processors

type ProblemsFromContest struct {
	ContestTitle string
	Problems     []string
}

type UnsolvedData struct {
	total    int
	unsolved []*ProblemsFromContest
	// maybe we will need more data
}

type Value struct {
	Value string
	Color string
	// these values are for required problems that aren't solved
	SpecialValue string
	SpecialColor string
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
