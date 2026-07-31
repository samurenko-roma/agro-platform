package valueobject

import (
	"encoding/json"
	"fmt"
)

type GeometryType string

const (
	Rect    GeometryType = "RECT"
	Polygon GeometryType = "POLYGON"
	Point   GeometryType = "POINT"
)

type Geometry struct {
	Type GeometryType `json:"type"`

	Dimension *Dimension `json:"dimension,omitempty"`

	Position *Coordinates `json:"position,omitempty"`

	Polygon [][]float64 `json:"polygon,omitempty"`
}

type geoJSON struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

// ToGeoJSON конвертирует доменную геометрию в GeoJSON для
// PostGIS ST_GeomFromGeoJSON.
//   - POINT   → GeoJSON Point
//   - RECT    → GeoJSON Polygon (прямоугольник строится из Position как
//     нижнего левого угла + Dimension.Width/Length)
//   - POLYGON → GeoJSON Polygon (кольцо замыкается автоматически, если
//     клиент не замкнул его сам)
func (g Geometry) ToGeoJSON() ([]byte, error) {
	switch g.Type {
	case Point:
		if g.Position == nil {
			return nil, fmt.Errorf("geometry: POINT requires position")
		}
		return json.Marshal(map[string]any{
			"type":        "Point",
			"coordinates": []float64{g.Position.X, g.Position.Y},
		})

	case Rect:
		if g.Position == nil || g.Dimension == nil {
			return nil, fmt.Errorf("geometry: RECT requires position and dimension")
		}
		x0, y0 := g.Position.X, g.Position.Y
		x1, y1 := x0+g.Dimension.Width, y0+g.Dimension.Length
		ring := [][]float64{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}, {x0, y0}}
		return json.Marshal(map[string]any{
			"type":        "Polygon",
			"coordinates": [][][]float64{ring},
		})

	case Polygon:
		if len(g.Polygon) < 3 {
			return nil, fmt.Errorf("geometry: POLYGON requires at least 3 points")
		}
		ring := append([][]float64{}, g.Polygon...)
		first, last := ring[0], ring[len(ring)-1]
		if first[0] != last[0] || first[1] != last[1] {
			ring = append(ring, first)
		}
		return json.Marshal(map[string]any{
			"type":        "Polygon",
			"coordinates": [][][]float64{ring},
		})

	default:
		return nil, fmt.Errorf("geometry: unknown type %q", g.Type)
	}
}

// GeometryFromGeoJSON строит доменный VO из GeoJSON, который вернул
// PostGIS через ST_AsGeoJSON.
//
// Важно: различить RECT/POLYGON постфактум невозможно — PostGIS хранит
// только полигон как таковой, метаинформация "это был прямоугольник"
// теряется при записи. Поэтому всё, что пришло из БД как Polygon,
// нормализуется в GeometryType Polygon. Это осознанный компромисс:
// геометрически прямоугольник — частный случай полигона, для отображения
// на карте и для площади разницы нет; RECT остаётся удобной формой ТОЛЬКО
// на входе API (клиенту проще прислать угол+размеры, чем 5 точек контура).
func GeometryFromGeoJSON(raw []byte) (Geometry, error) {
	var g geoJSON
	if err := json.Unmarshal(raw, &g); err != nil {
		return Geometry{}, err
	}

	switch g.Type {
	case "Point":
		var coords []float64
		if err := json.Unmarshal(g.Coordinates, &coords); err != nil {
			return Geometry{}, err
		}
		if len(coords) < 2 {
			return Geometry{}, fmt.Errorf("geometry: invalid point coordinates")
		}
		return Geometry{
			Type:     Point,
			Position: &Coordinates{X: coords[0], Y: coords[1]},
		}, nil

	case "Polygon":
		var rings [][][]float64
		if err := json.Unmarshal(g.Coordinates, &rings); err != nil {
			return Geometry{}, err
		}
		if len(rings) == 0 {
			return Geometry{}, fmt.Errorf("geometry: empty polygon")
		}
		return Geometry{
			Type:    Polygon,
			Polygon: rings[0],
		}, nil

	default:
		return Geometry{}, fmt.Errorf("geometry: unsupported GeoJSON type %q", g.Type)
	}
}
