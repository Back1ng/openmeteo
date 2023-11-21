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

	weather, err := api.GetWeather(context.Background())
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
				for {
					weather, err := api.GetWeather(context.Background())
					if err != nil {
						fmt.Printf("Error: failed update weather: %v", err)
					} else {
						cache.Store(entity.Weather{
							Temp: weather.Temp,
						})

						break
					}
				}
			}
		}
	}()

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		if cache.IsValid() {
			weatherCache, _ := cache.Get()
			WriteResponse(w, fmt.Sprint(weatherCache.Temp))
		} else {
			weatherCache, _ := cache.Get()
			WriteResponse(w, fmt.Sprint(weatherCache.Temp))
			w.WriteHeader(203)
		}
	})

	fmt.Println("Initializing webserver :3334...")

	_ = http.ListenAndServe(":3334", nil)
}
