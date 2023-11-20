package internal

import (
	"fmt"
	"github.com/Back1ng/openmeteo/internal/entity"
	gismeteo_api "github.com/Back1ng/openmeteo/internal/openmeteo-api"
	weather_cache "github.com/Back1ng/openmeteo/internal/weather-cache"
	"net/http"
)

func Run() {
	cache := weather_cache.New()

	api := gismeteo_api.New(cache)

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		if cache.IsValid() {
			weatherCache, ok := cache.Get()

			if ok {
				_, err := fmt.Fprintf(w, "Current weather: %v", weatherCache.Temp)
				if err != nil {
					fmt.Fprintf(w, "Error: %v", err)
					w.WriteHeader(500)
				}
				return
			}
		}

		weather, err := api.GetWeather()

		if err != nil {
			fmt.Fprint(w, err)
			w.WriteHeader(500)
			return
		}

		cache.Store(entity.Weather{
			Temp: weather.Temp,
		})

		weatherCache, ok := cache.Get()
		if ok {
			_, err := fmt.Fprintf(w, "Current weather: %v", weatherCache.Temp)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
				w.WriteHeader(500)
				return
			}
		}

		if !ok {
			fmt.Fprint(w, "Couldn't get the weather...")
		}
	})

	_ = http.ListenAndServe(":3334", nil)
}
