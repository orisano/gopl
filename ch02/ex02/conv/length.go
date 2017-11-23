package conv

import "fmt"

type Feet float64

func (f Feet) String() string { return fmt.Sprintf("%g ft", f) }
func (f Feet) ToMeter() Meter { return Meter(f * 0.3048) }

type Meter float64

func (m Meter) String() string { return fmt.Sprintf("%g m", m) }
func (m Meter) ToFeet() Feet   { return Feet(m * 3.28084) }
