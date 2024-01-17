package main

import (
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	GroupShort string `json:"group_short"`
}

type Problem struct {
	Id    string `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Index int    `json:"index"`
}

type UserSubmit struct {
	Score   int    `json:"score"`
	Penalty int    `json:"penalty"`
	Verdict string `json:"verdict"`
	Time    int    `json:"time"`
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

func getSubmitsData(url string) (data *SubmitsData) {
	client := resty.New()
	res, err := client.R().SetResult(&data).Get(url)
	if err != nil {
		slog.Warn("Error while querying algocode: %v\n", err.Error())
		return nil
	}
	if res.StatusCode() != 200 {
		slog.Warn("Algocode returned code %v\n", res.StatusCode())
		return nil
	}
	return
}
