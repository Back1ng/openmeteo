package internal

import (
	"context"
	"fmt"
	"github.com/Back1ng/openmeteo/internal/entity"
	openmeteoapi "github.com/Back1ng/openmeteo/internal/openmeteo-api"
	weathercache "github.com/Back1ng/openmeteo/internal/weather-cache"
	"log"
	"net/http"
	"time"
)

func WriteResponse(w http.ResponseWriter, s string) {
	_, err := fmt.Fprintf(w, "Current weather: %v", s)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		w.WriteHeader(500)
	}
}

func GetWeather(api openmeteoapi.Requester, cache *weathercache.Cache, w http.ResponseWriter) bool {
	if cache.IsValid() {
		weatherCache, _ := cache.Get()

		WriteResponse(w, fmt.Sprint(weatherCache.Temp))

		return true
	}

	weather, err := api.GetWeather()
	if err != nil {
		weatherCache, _ := cache.Get()

		WriteResponse(w, fmt.Sprint(weatherCache.Temp))

		return true
	}

	cache.Store(entity.Weather{
		Temp: weather.Temp,
	})

	WriteResponse(w, fmt.Sprint(weather.Temp))

	return true
}

func Run() {
	cache := weathercache.New()
	api := openmeteoapi.New(cache)

	weather, err := api.GetWeather()
	if err != nil {
		log.Fatal(err)
	}

	cache.Store(entity.Weather{
		Temp: weather.Temp,
	})

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		done := make(chan bool)
		go func() {
			done <- GetWeather(api, cache, w)
		}()

		select {
		case <-ctx.Done():
			weatherCache, _ := cache.Get()

			WriteResponse(w, fmt.Sprint(weatherCache.Temp))
		case <-done:
			return
		}
	})

	_ = http.ListenAndServe(":3334", nil)
}
