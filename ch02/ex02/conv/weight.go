package conv

import "fmt"

type Pound float64

func (p Pound) String() string       { return fmt.Sprintf("%glb", p) }
func (p Pound) ToKiloGram() KiloGram { return KiloGram(p * 0.45359237) }

type KiloGram float64

func (kg KiloGram) String() string { return fmt.Sprintf("%gkg", kg) }
func (kg KiloGram) ToPound() Pound { return Pound(kg / 0.45359237) }
