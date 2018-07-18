package main

import (
	"bytes"
	"testing"

	"github.com/faiface/pixel"
)

func TestLinePixels(t *testing.T) {
	tcs := []struct {
		nLines   float64
		index    float64
		expected []pixel.Vec
	}{
		{2, 0, []pixel.Vec{
			pixel.V(1, 0),
			pixel.V(1.5, 0),
			pixel.V(2, 0),
			pixel.V(2, 0.5),
			pixel.V(2, 1),
			pixel.V(2, 1.5),
			pixel.V(2, 2),
			pixel.V(2, 2.5),
			pixel.V(2, 3),
			pixel.V(1.5, 3),
			pixel.V(1, 3),
			pixel.V(1, 2.5),
			pixel.V(1, 2),
			pixel.V(1, 1.5),
			pixel.V(1, 1),
			pixel.V(1, 0.5),
		}},
	}

	for _, tc := range tcs {
		l := newLine(tc.nLines, tc.index)

		// Verify # points
		if len(l.points) != len(tc.expected) {
			t.Fatalf("got len(l.points)=%v, expected %v", len(l.points), len(tc.expected))
		}

		// Verify points
		for i := range l.points {
			if l.points[i] != tc.expected[i] {
				t.Errorf("got %v, want %v", l.points[i], tc.expected[i])
			}
		}
	}
}

func TestZMask(t *testing.T) {

	tcs := []struct {
		name     string
		nLines   float64
		index    float64
		expected []byte
	}{
		{"2 lines 1st", 2, 0, []byte{2, 2, 2, 1, 1, 0, 0, 2, 2, 2, 2, 1, 1, 0, 0, 2}},
		{"2 lines 2nd", 2, 1, []byte{2, 1, 1, 0, 0, 2, 2, 2, 2, 1, 1, 0, 0, 2, 2, 2}},
		{"3 lines: 2nd", 3, 1, []byte{2, 1, 1, 0, 0, 2, 2, 1, 1, 0, 0, 2, 2, 1, 1, 0, 0, 2, 2, 1, 1, 0, 0, 2}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := newLine(tc.nLines, tc.index)

			// Verify # points
			if len(l.zMask) != len(tc.expected) {
				t.Errorf("got len(l.points)=%v, expected %v", len(l.zMask), len(tc.expected))
			}

			// Verify points
			if !bytes.Equal(l.zMask, tc.expected) {
				t.Errorf("\n\tgot  %v\n\twant %v", l.zMask, tc.expected)
			}
		})
	}
}
