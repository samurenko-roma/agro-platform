package weatherrecord

// Source — откуда пришли данные.
// Расширяемо: добавить новый источник = новая константа + новый адаптер.
type Source string

const (
	SourceOpenMeteo Source = "OPEN_METEO"
	SourceSensor    Source = "SENSOR" // датчик с фермы
	SourceCustom    Source = "CUSTOM" // любой другой внешний сервис
)
