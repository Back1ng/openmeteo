package openmeteo_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Back1ng/openmeteo/internal/entity"
	weathercache "github.com/Back1ng/openmeteo/internal/weather-cache"
	"io"
	"net/http"
	"net/url"
)

type Requester struct {
	cache *weathercache.Cache
}

func New(cache *weathercache.Cache) Requester {
	return Requester{
		cache: cache,
	}
}

func (r *Requester) GetWeather(ctx context.Context) (*entity.Weather, error) {
	openMeteoURL, _ := url.Parse(
		fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
			56.837244,
			60.597647,
		),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", openMeteoURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	weather := struct {
		Current struct {
			Temperature float32 `json:"temperature_2m"`
		} `json:"current"`
	}{}

	responseBody, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(responseBody, &weather)
	if err != nil {
		return nil, err
	}

	return &entity.Weather{
		Temp: weather.Current.Temperature,
	}, nil
}
