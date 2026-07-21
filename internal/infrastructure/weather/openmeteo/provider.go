package openmeteo

import (
	"context"
	"fmt"
	"time"

	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
	weatherprovider "github.com/samurenkoroma/agro-platform/internal/domain/weather/provider"
)

type Provider struct {
	client *client
}

func NewProvider() *Provider {
	return &Provider{client: newClient()}
}

func (p *Provider) Name() weatherrecord.Source {
	return weatherrecord.SourceOpenMeteo
}

func (p *Provider) FetchCurrent(ctx context.Context, location *weatherlocation.WeatherLocation) (*weatherrecord.WeatherData, error) {
	resp, err := p.client.fetchForecast(ctx, location.Latitude, location.Longitude, 1)
	if err != nil {
		return nil, fmt.Errorf("open-meteo FetchCurrent: %w", err)
	}

	data := &weatherrecord.WeatherData{
		Temperature:   float64Ptr(resp.CurrentWeather.Temperature),
		WindSpeed:     float64Ptr(resp.CurrentWeather.Windspeed),
		WindDirection: float64Ptr(resp.CurrentWeather.Winddirection),
		WeatherCode:   intPtr(resp.CurrentWeather.WeatherCode),
	}
	isDay := resp.CurrentWeather.IsDay == 1
	data.IsDay = &isDay
	return data, nil
}

func (p *Provider) FetchForecast(ctx context.Context, location *weatherlocation.WeatherLocation, days int) ([]weatherprovider.ForecastPoint, error) {
	if days < 1 || days > 16 {
		days = 7
	}
	resp, err := p.client.fetchForecast(ctx, location.Latitude, location.Longitude, days)
	if err != nil {
		return nil, fmt.Errorf("open-meteo FetchForecast: %w", err)
	}
	return p.mapHourlyToForecast(resp), nil
}

func (p *Provider) FetchHistorical(ctx context.Context, location *weatherlocation.WeatherLocation, from, to time.Time) ([]weatherprovider.HistoricalPoint, error) {
	resp, err := p.client.fetchHistorical(ctx, location.Latitude, location.Longitude, from, to)
	if err != nil {
		return nil, fmt.Errorf("open-meteo FetchHistorical: %w", err)
	}
	return p.mapHourlyToHistorical(resp), nil
}

func (p *Provider) mapHourlyToForecast(resp *forecastResponse) []weatherprovider.ForecastPoint {
	points := make([]weatherprovider.ForecastPoint, 0, len(resp.Hourly.Time))
	for i, timeStr := range resp.Hourly.Time {
		t, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			continue
		}
		points = append(points, weatherprovider.ForecastPoint{
			Time: t,
			Data: p.mapHourlyIndex(resp, i),
		})
	}
	return points
}

func (p *Provider) mapHourlyToHistorical(resp *forecastResponse) []weatherprovider.HistoricalPoint {
	points := make([]weatherprovider.HistoricalPoint, 0, len(resp.Hourly.Time))
	for i, timeStr := range resp.Hourly.Time {
		t, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			continue
		}
		points = append(points, weatherprovider.HistoricalPoint{
			Time: t,
			Data: p.mapHourlyIndex(resp, i),
		})
	}
	return points
}

func (p *Provider) mapHourlyIndex(resp *forecastResponse, i int) weatherrecord.WeatherData {
	h := resp.Hourly
	d := weatherrecord.WeatherData{}

	if i < len(h.Temperature2m) {
		d.Temperature = float64Ptr(h.Temperature2m[i])
	}
	if i < len(h.RelativeHumidity2m) {
		d.Humidity = float64Ptr(h.RelativeHumidity2m[i])
	}
	if i < len(h.DewPoint2m) {
		d.DewPoint = float64Ptr(h.DewPoint2m[i])
	}
	if i < len(h.ApparentTemperature) {
		d.TemperatureFeelsLike = float64Ptr(h.ApparentTemperature[i])
	}
	if i < len(h.PrecipitationProbability) {
		d.PrecipitationProb = float64Ptr(h.PrecipitationProbability[i])
	}
	if i < len(h.Precipitation) {
		d.Precipitation = float64Ptr(h.Precipitation[i])
	}
	if i < len(h.Rain) {
		d.Rain = float64Ptr(h.Rain[i])
	}
	if i < len(h.Snowfall) {
		d.Snowfall = float64Ptr(h.Snowfall[i])
	}
	if i < len(h.WeatherCode) {
		d.WeatherCode = intPtr(h.WeatherCode[i])
	}
	if i < len(h.CloudCover) {
		d.CloudCover = float64Ptr(h.CloudCover[i])
	}
	if i < len(h.SurfacePressure) {
		d.PressureSea = float64Ptr(h.SurfacePressure[i])
	}
	if i < len(h.WindSpeed10m) {
		d.WindSpeed = float64Ptr(h.WindSpeed10m[i])
	}
	if i < len(h.WindDirection10m) {
		d.WindDirection = float64Ptr(h.WindDirection10m[i])
	}
	if i < len(h.WindGusts10m) {
		d.WindGusts = float64Ptr(h.WindGusts10m[i])
	}
	if i < len(h.ShortwaveRadiation) {
		d.SolarRadiation = float64Ptr(h.ShortwaveRadiation[i])
	}
	if i < len(h.UVIndex) {
		d.UVIndex = float64Ptr(h.UVIndex[i])
	}
	if i < len(h.SunshineDuration) {
		d.SunshineDuration = float64Ptr(h.SunshineDuration[i])
	}
	if i < len(h.SoilTemperature0cm) {
		d.SoilTemperature0cm = float64Ptr(h.SoilTemperature0cm[i])
	}
	if i < len(h.SoilTemperature6cm) {
		d.SoilTemperature6cm = float64Ptr(h.SoilTemperature6cm[i])
	}
	if i < len(h.SoilMoisture0to1cm) {
		d.SoilMoisture0_1cm = float64Ptr(h.SoilMoisture0to1cm[i])
	}
	if i < len(h.SoilMoisture1to3cm) {
		d.SoilMoisture1_3cm = float64Ptr(h.SoilMoisture1to3cm[i])
	}
	if i < len(h.ET0FaoEvapotranspiration) {
		d.Evapotranspiration = float64Ptr(h.ET0FaoEvapotranspiration[i])
	}
	if i < len(h.VapourPressureDeficit) {
		d.VapourPressureDef = float64Ptr(h.VapourPressureDeficit[i])
	}
	if i < len(h.Visibility) {
		d.Visibility = float64Ptr(h.Visibility[i])
	}
	if i < len(h.IsDay) {
		isDay := h.IsDay[i] == 1
		d.IsDay = &isDay
	}
	return d
}

func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int             { return &v }
