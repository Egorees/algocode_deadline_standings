package algocode

import (
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
)

func GetSubmitsData(url string) (data *SubmitsData) {
	client := resty.New()
	res, err := client.R().SetResult(&data).Get(url)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error while querying algocode: %v\n", err.Error()))
		return nil
	}
	if res.StatusCode() != 200 {
		slog.Warn(fmt.Sprintf("Algocode returned code %v\n", res.StatusCode()))
		return nil
	}
	return data
}
