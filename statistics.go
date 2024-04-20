package main

import (
	"strconv"
)

//type People struct {
//	FullName    string `json:"FullName"`
//	NeedToSolve int    `json:"NeedToSolve"`
//}

type Stats struct {
	Count   int      `json:"Count"`
	Color   string   `json:"Color"`
	Peoples []string `json:"Peoples"`
}

//func writeToFile(name string, data any) {
//	bytes, err := json.MarshalIndent(data, "", "    ")
//	if err != nil {
//		fmt.Println("Meow")
//	}
//	fl, _ := os.Create(name)
//	defer fl.Close()
//	fl.Write(bytes)
//}

func statisticsFun(config *Config, userValues []*UserValues) map[int]*Stats {
	//data := make([]*People, 0)
	stat := make(map[int]*Stats)
	for _, el := range userValues {
		cnt, _ := strconv.Atoi(el.Values[0].Value)
		val, exs := stat[cnt]
		if !exs {
			stat[cnt] = &Stats{}
			val = stat[cnt]
			val.Peoples = make([]string, 0)
			val.Count = 0
		}
		stat[cnt].Color = config.GetColorByCount(cnt)
		stat[cnt].Peoples = append(stat[cnt].Peoples, el.FullName)
		stat[cnt].Count++
		//data = append(data, &People{
		//	FullName:    el.FullName,
		//	NeedToSolve: cnt,
		//})
	}
	return stat
}
