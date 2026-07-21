package weatherrecord

type Kind string

const (
	Current    Kind = "CURRENT"
	Forecast   Kind = "FORECAST"
	Historical Kind = "HISTORICAL"
)
