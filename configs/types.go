package configs

type DeadlineTasks = map[string][]string

type DeadlineData struct {
	Tasks DeadlineTasks `yaml:"deadline"`
}

type UnsolvedBorder struct {
	Count int    `json:"count"`
	Color string `json:"color"`
}

type Config struct {
	CacheTime         float64           `yaml:"cache_time"`
	FullSolveText     string            `yaml:"full_solve_text"`
	UnsolvedBorders   []*UnsolvedBorder `yaml:"unsolved_borders"`
	ServerAddressPort string            `yaml:"server_address_port"`
	SubmitsLink       string            `yaml:"submits_link"`
	DeadlineFilepath  string            `yaml:"deadline_filepath"`
	ReleaseMode       bool              `yaml:"release_mode"`
}
