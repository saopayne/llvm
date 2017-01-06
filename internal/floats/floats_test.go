package floats

import (
	"math"
	"testing"
)

func TestFloat16Float32(t *testing.T) {
	golden := []struct {
		in   string
		want float32
	}{
		{in: "3C00", want: 1},
		{in: "4000", want: 2},
		{in: "C000", want: -2},
		{in: "7BFE", want: 65472},
		{in: "7BFF", want: 65504},
		{in: "FBFF", want: -65504},
		{in: "0000", want: 0},
		{in: "8000", want: float32(math.Copysign(0, -1))},
		{in: "7C00", want: float32(math.Inf(1))},
		{in: "FC00", want: float32(math.Inf(-1))},
		{in: "5B8F", want: 241.875},
		{in: "48C8", want: 9.5625},
	}
	for _, g := range golden {
		f := NewFloat16FromString(g.in)
		got := f.Float32()
		if got != g.want {
			t.Errorf("float32 mismatch for binary16 0x%04X; expected %v, got %v", g.in, g.want, got)
		}
	}
}

func TestFloat16Float64(t *testing.T) {
	golden := []struct {
		in   uint16
		want float64
	}{
		{in: 0x3C00, want: 1},
		{in: 0x4000, want: 2},
		{in: 0xC000, want: -2},
		{in: 0x7BFE, want: 65472},
		{in: 0x7BFF, want: 65504},
		{in: 0xFBFF, want: -65504},
		{in: 0x0000, want: 0},
		{in: 0x8000, want: math.Copysign(0, -1)},
		{in: 0x7C00, want: math.Inf(1)},
		{in: 0xFC00, want: math.Inf(-1)},
		{in: 0x5B8F, want: 241.875},
		{in: 0x48C8, want: 9.5625},
	}
	for _, g := range golden {
		f := NewFloat16FromBits(g.in)
		got := f.Float64()
		if got != g.want {
			t.Errorf("float64 mismatch for binary16 0x%04X; expected %v, got %v", g.in, g.want, got)
		}
	}
}