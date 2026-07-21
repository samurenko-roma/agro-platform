package repository

type WeatherProvider interface {
	Locations() WeatherLocationRepository
	Records() WeatherRecordRepository
}
