package productionunit

import (
	"context"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type CapacityDTO struct {
	PlantCapacity     *int     `json:"plantCapacity"`
	WaterVolumeLiters *float64 `json:"waterVolumeLiters"`
	AreaM2            *float64 `json:"areaM2"`
	TrayCount         *int     `json:"trayCount"`
	ChannelCount      *int     `json:"channelCount"`
}

type ClimateDTO struct {
	TemperatureMin *float64 `json:"temperatureMin"`
	TemperatureMax *float64 `json:"temperatureMax"`
	HumidityMin    *float64 `json:"humidityMin"`
	HumidityMax    *float64 `json:"humidityMax"`
	CO2Min         *float64 `json:"co2Min"`
	CO2Max         *float64 `json:"co2Max"`
	LightPPFDMin   *float64 `json:"lightPpfdMin"`
	LightPPFDMax   *float64 `json:"lightPpfdMax"`
}

type DTO struct {
	ID         vo.ID          `json:"id"`
	ParentID   *vo.ID         `json:"parentId"`
	Type       string         `json:"type"`
	Area       float64        `json:"area"`
	Status     string         `json:"status"`
	Code       string         `json:"code"`
	Geometry   *vo.Geometry   `json:"geometry"`
	Properties map[string]any `json:"properties"`
	Children   []*DTO         `json:"children"`
}

type Projection interface {
	Get(ctx context.Context, id vo.ID) (*DTO, error)
	ListRoots(ctx context.Context, ownerId vo.ID) ([]*DTO, error)
	Tree(ctx context.Context, rootID *vo.ID) (*DTO, error)
	Children(ctx context.Context, parentID vo.ID) ([]*DTO, error)
}
