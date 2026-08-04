package valueobject_test

import (
	"encoding/json"
	"testing"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

func TestGeometry_Point_RoundTrip(t *testing.T) {
	g := vo.Geometry{Type: vo.Point, Position: &vo.Coordinates{X: 30.5, Y: 46.9}}

	raw, err := g.ToGeoJSON()
	if err != nil {
		t.Fatalf("ToGeoJSON: %v", err)
	}

	back, err := vo.GeometryFromGeoJSON(raw)
	if err != nil {
		t.Fatalf("GeometryFromGeoJSON: %v", err)
	}
	if back.Type != vo.Point {
		t.Errorf("type: got %s, want POINT", back.Type)
	}
	if back.Position == nil || back.Position.X != 30.5 || back.Position.Y != 46.9 {
		t.Errorf("position mismatch: %+v", back.Position)
	}
}

func TestGeometry_Rect_BecomesClosedPolygon(t *testing.T) {
	g := vo.Geometry{
		Type:      vo.Rect,
		Position:  &vo.Coordinates{X: 0, Y: 0},
		Dimension: &vo.Dimension{Width: 10, Length: 5},
	}

	raw, err := g.ToGeoJSON()
	if err != nil {
		t.Fatalf("ToGeoJSON: %v", err)
	}

	var parsed map[string]any
	json.Unmarshal(raw, &parsed)
	if parsed["type"] != "Polygon" {
		t.Errorf("geoJSON type: got %v, want Polygon", parsed["type"])
	}

	back, err := vo.GeometryFromGeoJSON(raw)
	if err != nil {
		t.Fatalf("GeometryFromGeoJSON: %v", err)
	}
	if back.Type != vo.Polygon {
		t.Errorf("type after round-trip: got %s, want POLYGON (RECT нормализуется в POLYGON)", back.Type)
	}
	if len(back.Polygon) != 5 {
		t.Errorf("expected closed ring with 5 points (4 угла + повтор первой), got %d", len(back.Polygon))
	}
	first, last := back.Polygon[0], back.Polygon[len(back.Polygon)-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("ring should be closed (first point == last point)")
	}
}

func TestGeometry_Polygon_AutoClosesRing(t *testing.T) {
	g := vo.Geometry{
		Type:    vo.Polygon,
		Polygon: [][]float64{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, // не замкнут
	}

	raw, err := g.ToGeoJSON()
	if err != nil {
		t.Fatalf("ToGeoJSON: %v", err)
	}

	back, err := vo.GeometryFromGeoJSON(raw)
	if err != nil {
		t.Fatalf("GeometryFromGeoJSON: %v", err)
	}
	if len(back.Polygon) != 5 {
		t.Errorf("expected ring auto-closed to 5 points, got %d", len(back.Polygon))
	}
}

func TestGeometry_Polygon_TooFewPoints_ReturnsError(t *testing.T) {
	g := vo.Geometry{Type: vo.Polygon, Polygon: [][]float64{{0, 0}, {1, 1}}}
	if _, err := g.ToGeoJSON(); err == nil {
		t.Error("expected error for polygon with fewer than 3 points")
	}
}

func TestGeometry_Point_WithoutPosition_ReturnsError(t *testing.T) {
	g := vo.Geometry{Type: vo.Point}
	if _, err := g.ToGeoJSON(); err == nil {
		t.Error("expected error for POINT without position")
	}
}

func TestGeometryFromGeoJSON_UnsupportedType_ReturnsError(t *testing.T) {
	raw := []byte(`{"type":"LineString","coordinates":[[0,0],[1,1]]}`)
	if _, err := vo.GeometryFromGeoJSON(raw); err == nil {
		t.Error("expected error for unsupported GeoJSON type")
	}
}
