package productionunit_test

import (
	"testing"

	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

func TestAreaUnit_ToSquareMeters(t *testing.T) {
	cases := []struct {
		unit  pu.AreaUnit
		value float64
		want  float64
	}{
		{pu.SquareMeters, 100, 100},
		{"", 100, 100}, // пустая единица трактуется как м²
		{pu.Hectares, 50, 500000},
		{pu.Sotka, 20, 2000},
	}
	for _, c := range cases {
		got, err := c.unit.ToSquareMeters(c.value)
		if err != nil {
			t.Fatalf("unit=%s: unexpected error: %v", c.unit, err)
		}
		if got != c.want {
			t.Errorf("unit=%s value=%v: got %v, want %v", c.unit, c.value, got, c.want)
		}
	}
}

func TestAreaUnit_UnknownUnit_ReturnsError(t *testing.T) {
	_, err := pu.AreaUnit("ACRE").ToSquareMeters(1)
	if err == nil {
		t.Error("expected error for unknown area unit")
	}
}
