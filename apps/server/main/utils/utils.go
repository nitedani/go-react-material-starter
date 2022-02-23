package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/reactivex/rxgo/v2"
)

func Body(c echo.Context) any {
	parsed := new(any)
	c.Bind(&parsed)
	return parsed
}

func parseJson(response *resty.Response) map[string]interface{} {
	parsed := map[string]interface{}{}
	err := json.Unmarshal(response.Body(), &parsed)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return parsed
}

func doRequest(url string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {

		client := resty.New()
		res, _ := client.R().
			SetHeader("Accept", "application/json").
			Get(url)

		obj := parseJson(res)
		next <- rxgo.Of(obj)

	}})
}
