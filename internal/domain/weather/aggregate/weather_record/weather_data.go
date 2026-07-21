package weatherrecord

// WeatherData — набор метеопоказателей в одной записи.
// Все поля опциональны: разные источники дают разный набор данных.
type WeatherData struct {
	// Температура и влажность
	Temperature          *float64 // °C
	TemperatureFeelsLike *float64 // °C
	Humidity             *float64 // %
	DewPoint             *float64 // °C

	// Осадки
	Precipitation     *float64 // мм
	PrecipitationProb *float64 // % (только для прогноза)
	Rain              *float64 // мм
	Snowfall          *float64 // см

	// Ветер
	WindSpeed     *float64 // км/ч
	WindDirection *float64 // градусы
	WindGusts     *float64 // км/ч

	// Давление и облачность
	PressureSea *float64 // гПа
	CloudCover  *float64 // %

	// Солнце
	SolarRadiation   *float64 // Вт/м²
	UVIndex          *float64
	SunshineDuration *float64 // секунды

	// Почва (актуально для агро)
	SoilTemperature0cm *float64 // °C
	SoilTemperature6cm *float64 // °C
	SoilMoisture0_1cm  *float64 // м³/м³
	SoilMoisture1_3cm  *float64 // м³/м³

	// Испарение (актуально для агро)
	Evapotranspiration *float64 // мм
	VapourPressureDef  *float64 // кПа (VPD)

	// Прочее
	WeatherCode *int     // WMO code
	Visibility  *float64 // м
	IsDay       *bool
}
