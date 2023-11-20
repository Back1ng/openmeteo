package openmeteo_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Back1ng/openmeteo/internal/entity"
	weather_cache "github.com/Back1ng/openmeteo/internal/weather-cache"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Requester struct {
	cache *weather_cache.Cache
}

func New(cache *weather_cache.Cache) Requester {
	return Requester{
		cache: cache,
	}
}

func (r *Requester) GetWeather() (*entity.Weather, error) {
	openMeteoURL, _ := url.Parse(
		fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
			56.837244,
			60.597647,
		),
	)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", openMeteoURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

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
