package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	forecastBaseURL   = "https://api.open-meteo.com/v1/forecast"
	historicalBaseURL = "https://archive-api.open-meteo.com/v1/archive"
)

// hourlyVars — переменные для почасового запроса. Ключевые для агро.
var hourlyVars = []string{
	"temperature_2m",
	"relative_humidity_2m",
	"dew_point_2m",
	"apparent_temperature",
	"precipitation_probability",
	"precipitation",
	"rain",
	"snowfall",
	"weather_code",
	"cloud_cover",
	"surface_pressure",
	"wind_speed_10m",
	"wind_direction_10m",
	"wind_gusts_10m",
	"shortwave_radiation",
	"uv_index",
	"sunshine_duration",
	"soil_temperature_0cm",
	"soil_temperature_6cm",
	"soil_moisture_0_to_1cm",
	"soil_moisture_1_to_3cm",
	"et0_fao_evapotranspiration",
	"vapour_pressure_deficit",
	"visibility",
	"is_day",
}

var currentVars = []string{
	"temperature_2m",
	"relative_humidity_2m",
	"apparent_temperature",
	"precipitation",
	"rain",
	"weather_code",
	"cloud_cover",
	"surface_pressure",
	"wind_speed_10m",
	"wind_direction_10m",
	"wind_gusts_10m",
	"is_day",
}

type client struct {
	http    *http.Client
	timeout time.Duration
}

func newClient() *client {
	return &client{
		http:    &http.Client{Timeout: 15 * time.Second},
		timeout: 15 * time.Second,
	}
}

type forecastResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time   string             `json:"time"`
		Values map[string]float64 `json:"values"` // не совсем так — см. ниже
	} `json:"current"`
	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2m            []float64 `json:"temperature_2m"`
		RelativeHumidity2m       []float64 `json:"relative_humidity_2m"`
		DewPoint2m               []float64 `json:"dew_point_2m"`
		ApparentTemperature      []float64 `json:"apparent_temperature"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
		Precipitation            []float64 `json:"precipitation"`
		Rain                     []float64 `json:"rain"`
		Snowfall                 []float64 `json:"snowfall"`
		WeatherCode              []int     `json:"weather_code"`
		CloudCover               []float64 `json:"cloud_cover"`
		SurfacePressure          []float64 `json:"surface_pressure"`
		WindSpeed10m             []float64 `json:"wind_speed_10m"`
		WindDirection10m         []float64 `json:"wind_direction_10m"`
		WindGusts10m             []float64 `json:"wind_gusts_10m"`
		ShortwaveRadiation       []float64 `json:"shortwave_radiation"`
		UVIndex                  []float64 `json:"uv_index"`
		SunshineDuration         []float64 `json:"sunshine_duration"`
		SoilTemperature0cm       []float64 `json:"soil_temperature_0cm"`
		SoilTemperature6cm       []float64 `json:"soil_temperature_6cm"`
		SoilMoisture0to1cm       []float64 `json:"soil_moisture_0_to_1cm"`
		SoilMoisture1to3cm       []float64 `json:"soil_moisture_1_to_3cm"`
		ET0FaoEvapotranspiration []float64 `json:"et0_fao_evapotranspiration"`
		VapourPressureDeficit    []float64 `json:"vapour_pressure_deficit"`
		Visibility               []float64 `json:"visibility"`
		IsDay                    []int     `json:"is_day"`
	} `json:"hourly"`
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		WeatherCode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          string  `json:"time"`
	} `json:"current_weather"`
}

func (c *client) fetchForecast(ctx context.Context, lat, lon float64, days int) (*forecastResponse, error) {
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%.6f", lat))
	params.Set("longitude", fmt.Sprintf("%.6f", lon))
	params.Set("hourly", strings.Join(hourlyVars, ","))
	params.Set("current_weather", "true")
	params.Set("forecast_days", fmt.Sprintf("%d", days))
	params.Set("wind_speed_unit", "kmh")
	params.Set("timezone", "auto")

	return c.fetch(ctx, forecastBaseURL, params)
}

func (c *client) fetchHistorical(ctx context.Context, lat, lon float64, from, to time.Time) (*forecastResponse, error) {
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%.6f", lat))
	params.Set("longitude", fmt.Sprintf("%.6f", lon))
	params.Set("hourly", strings.Join(hourlyVars, ","))
	params.Set("start_date", from.Format("2006-01-02"))
	params.Set("end_date", to.Format("2006-01-02"))
	params.Set("wind_speed_unit", "kmh")
	params.Set("timezone", "auto")

	return c.fetch(ctx, historicalBaseURL, params)
}

func (c *client) fetch(ctx context.Context, baseURL string, params url.Values) (*forecastResponse, error) {
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("open-meteo responded %d: %s", resp.StatusCode, string(body))
	}

	var result forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
