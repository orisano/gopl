package conv

import "fmt"

type Celsius float64

func (c Celsius) String() string           { return fmt.Sprintf("%g°C", c) }
func (c Celsius) ToFahrenheit() Fahrenheit { return Fahrenheit(c*9/5 + 32) }

type Fahrenheit float64

func (f Fahrenheit) String() string     { return fmt.Sprintf("%g°F", f) }
func (f Fahrenheit) ToCelsius() Celsius { return Celsius((f - 32) * 5 / 9) }
