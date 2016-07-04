package expr

import (
	"fmt"
)

var (
	// Zero value
	Zero = CONST(0).Accumulator()
)

// CONST returns an Accumulator that always has a constant value.
func CONST(value float64) Expr {
	return &constant{value}
}

type constantAccumulator struct {
	value float64
}

func (a *constantAccumulator) Update(params Params) {
}

func (a *constantAccumulator) Get() float64 {
	return a.value
}

type constant struct {
	value float64
}

func (e *constant) Accumulator() Accumulator {
	return &constantAccumulator{e.value}
}

func (e *constant) DependsOn() []string {
	return []string{}
}

func (e *constant) String() string {
	return fmt.Sprintf("%f", e.value)
}
