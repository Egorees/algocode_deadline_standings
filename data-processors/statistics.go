package data_processors

import (
	"algocode_deadline_standings/configs"
	"strconv"
)

//func writeToFile(name string, data any) {
//	bytes, err := json.MarshalIndent(data, "", "    ")
//	if err != nil {
//		fmt.Println("Meow")
//	}
//	fl, _ := os.Create(name)
//	defer fl.Close()
//	fl.Write(bytes)
//}

func CreateStatistics(config *configs.Config, userValues []*UserValues) (map[int]*Stats, error) {
	if userValues == nil || config == nil {
		return nil, &DataError{Reason: "config or userValues is nil"}
	}
	stat := make(map[int]*Stats)
	for _, el := range userValues {
		cnt, _ := strconv.Atoi(el.Values[0].Value)
		if _, exs := stat[cnt]; !exs {
			stat[cnt] = &Stats{
				Peoples: make([]string, 0),
				Count:   0,
				Color:   "",
			}
		}
		if stat[cnt].Color == "" {
			stat[cnt].Color = config.GetColorByCount(cnt)
		}
		stat[cnt].Peoples = append(stat[cnt].Peoples, el.FullName)
		stat[cnt].Count++
	}
	return stat, nil
}
