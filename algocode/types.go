package algocode

import (
	"encoding/json"
	"log/slog"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	GroupShort string `json:"group_short"`
}

type Problem struct {
	Id    int    `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Index int    `json:"index"`
}

func (problem *Problem) UnmarshalJSON(b []byte) error { // not sure if this should be places here
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		slog.Warn("Json unmarshalling error:", err)
		panic(err)
	}
	problem.Id, _ = data["id"].(int)
	problem.Long, _ = data["long"].(string)
	problem.Short, _ = data["short"].(string)
	problem.Index, _ = data["index"].(int)
	return nil
}

type UserSubmit struct {
	Score   int    `json:"score"`
	Penalty int    `json:"penalty"`
	Verdict string `json:"verdict"`
//	Time    int    `json:"time"`
}

type Contest struct {
	Id          int                      `json:"id"`
	Date        string                   `json:"date"`
	EjudgeId    int                      `json:"ejudge_id"`
	Title       string                   `json:"title"`
	Coefficient float64                  `json:"coefficient"`
	Problems    []*Problem               `json:"problems"`
	Users       map[string][]*UserSubmit `json:"users"`
}

type SubmitsData struct {
	Users    []*User    `json:"users"`
	Contests []*Contest `json:"contests"`
}
