package productionunit

import "fmt"

// AreaUnit — единица измерения площади на входе API.
// Внутри домена и в БД площадь всегда хранится в квадратных метрах —
// это транспортный конвертер на границе, не доменное понятие сам по себе.
type AreaUnit string

const (
	SquareMeters AreaUnit = "M2"
	Hectares     AreaUnit = "HECTARE" // 1 га = 10 000 м²
	Sotka        AreaUnit = "SOTKA"   // 1 сотка (ар) = 100 м²
)

func (u AreaUnit) ToSquareMeters(value float64) (float64, error) {
	switch u {
	case SquareMeters, "":
		return value, nil
	case Hectares:
		return value * 10000, nil
	case Sotka:
		return value * 100, nil
	default:
		return 0, fmt.Errorf("production unit: unknown area unit %q", u)
	}
}

// AreaInput — площадь, заданная пользователем вручную при создании узла.
type AreaInput struct {
	Value float64  `json:"value" validate:"required,gt=0"`
	Unit  AreaUnit `json:"unit" validate:"required"`
}
