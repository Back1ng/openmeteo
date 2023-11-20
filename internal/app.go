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

func Run() {
	cache := weathercache.New()
	api := openmeteoapi.New(cache)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	weather, err := api.GetWeather(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cache.Store(entity.Weather{
		Temp: weather.Temp,
	})

	ticker := time.NewTicker(time.Second * 50)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Getting weather by timing")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				weather, err := api.GetWeather(ctx)
				if err != nil {
					fmt.Printf("Error: failed update weather: %v", err)
				}
				cancel()

				cache.Store(entity.Weather{
					Temp: weather.Temp,
				})
			}
		}
	}()

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		done := make(chan struct{})
		go func() {
			if cache.IsValid() {
				done <- struct{}{}
			}
		}()

		select {
		case <-ctx.Done():
		case <-done:
			weatherCache, _ := cache.Get()
			WriteResponse(w, fmt.Sprint(weatherCache.Temp))
			return
		}
	})

	_ = http.ListenAndServe(":3334", nil)
}
