package goptuna

import (
	"encoding/json"
	"math"
)

// Distribution represents a parameter that can be optimized.
type Distribution interface {
	// ToExternalRepr to convert internal representation of a parameter value into external representation.
	ToExternalRepr(float64) interface{}
	// Single to test whether the range of this distribution contains just a single value.
	Single() bool
	// Contains to check a parameter value is contained in the range of this distribution.
	Contains(float64) bool
}

var _ Distribution = &UniformDistribution{}

// UniformDistribution is a uniform distribution in the linear domain.
type UniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High float64 `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low float64 `json:"low"`
}

// UniformDistributionName is the identifier name of UniformDistribution
const UniformDistributionName = "UniformDistribution"

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *UniformDistribution) ToExternalRepr(ir float64) interface{} {
	return ir
}

// Single to test whether the range of this distribution contains just a single value.
func (d *UniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *UniformDistribution) Contains(ir float64) bool {
	if d.Single() {
		return ir == d.Low
	}
	return d.Low <= ir && ir < d.High
}

var _ Distribution = &LogUniformDistribution{}

// LogUniformDistribution is a uniform distribution in the log domain.
type LogUniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High float64 `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low float64 `json:"low"`
}

// LogUniformDistributionName is the identifier name of LogUniformDistribution
const LogUniformDistributionName = "LogUniformDistribution"

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *LogUniformDistribution) ToExternalRepr(ir float64) interface{} {
	return ir
}

// Single to test whether the range of this distribution contains just a single value.
func (d *LogUniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *LogUniformDistribution) Contains(ir float64) bool {
	if d.Single() {
		return ir == d.Low
	}
	return d.Low <= ir && ir < d.High
}

var _ Distribution = &IntUniformDistribution{}

// IntUniformDistribution is a uniform distribution on integers.
type IntUniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High int `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low int `json:"low"`
}

// IntUniformDistributionName is the identifier name of IntUniformDistribution
const IntUniformDistributionName = "IntUniformDistribution"

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *IntUniformDistribution) ToExternalRepr(ir float64) interface{} {
	return int(ir)
}

// Single to test whether the range of this distribution contains just a single value.
func (d *IntUniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *IntUniformDistribution) Contains(ir float64) bool {
	value := int(ir)
	if d.Single() {
		return value == d.Low
	}
	return d.Low <= value && value < d.High
}

var _ Distribution = &DiscreteUniformDistribution{}

// DiscreteUniformDistribution is a discretized uniform distribution in the linear domain.
type DiscreteUniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High float64 `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low float64 `json:"low"`
	// Q is a discretization step.
	Q float64 `json:"q"`
}

// DiscreteUniformDistributionName is the identifier name of DiscreteUniformDistribution
const DiscreteUniformDistributionName = "DiscreteUniformDistribution"

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *DiscreteUniformDistribution) ToExternalRepr(ir float64) interface{} {
	return math.Floor((ir-d.Low)/d.Q+0.5)*d.Q + d.Low
}

// Single to test whether the range of this distribution contains just a single value.
func (d *DiscreteUniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *DiscreteUniformDistribution) Contains(ir float64) bool {
	if d.Single() {
		return ir == d.Low
	}
	if d.Low > ir || ir > d.High {
		return false
	}

	eps := 1e-6
	if math.Mod(ir-d.Low+eps, d.Q) <= 2*eps {
		return true
	}
	if math.Mod(d.High-ir+eps, d.Q) <= 2*eps {
		return true
	}
	return false
}

var _ Distribution = &CategoricalDistribution{}

// CategoricalDistribution is a distribution for categorical parameters
type CategoricalDistribution struct {
	// Choices is a candidates of parameter values
	Choices []string `json:"choices"`
}

// CategoricalDistributionName is the identifier name of CategoricalDistribution
const CategoricalDistributionName = "CategoricalDistribution"

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *CategoricalDistribution) ToExternalRepr(ir float64) interface{} {
	return d.Choices[int(ir)]
}

// Single to test whether the range of this distribution contains just a single value.
func (d *CategoricalDistribution) Single() bool {
	return len(d.Choices) == 1
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *CategoricalDistribution) Contains(ir float64) bool {
	index := int(ir)
	return 0 <= index && index < len(d.Choices)
}

// ToExternalRepresentation converts to external representation
func ToExternalRepresentation(distribution interface{}, ir float64) (interface{}, error) {
	switch d := distribution.(type) {
	case UniformDistribution:
		return d.ToExternalRepr(ir), nil
	case LogUniformDistribution:
		return d.ToExternalRepr(ir), nil
	case IntUniformDistribution:
		return d.ToExternalRepr(ir), nil
	case DiscreteUniformDistribution:
		return d.ToExternalRepr(ir), nil
	case CategoricalDistribution:
		return d.ToExternalRepr(ir), nil
	default:
		return nil, ErrUnknownDistribution
	}
}

// DistributionIsSingle whether the distribution contains just a single value.
func DistributionIsSingle(distribution interface{}) (bool, error) {
	switch d := distribution.(type) {
	case UniformDistribution:
		return d.Single(), nil
	case LogUniformDistribution:
		return d.Single(), nil
	case IntUniformDistribution:
		return d.Single(), nil
	case DiscreteUniformDistribution:
		return d.Single(), nil
	case CategoricalDistribution:
		return d.Single(), nil
	default:
		return false, ErrUnknownDistribution
	}
}

// DistributionToJSON serialize a distribution to JSON format.
func DistributionToJSON(distribution interface{}) ([]byte, error) {
	var ir struct {
		Name  string      `json:"name"`
		Attrs interface{} `json:"attributes"`
	}
	switch distribution.(type) {
	case UniformDistribution:
		ir.Name = UniformDistributionName
	case LogUniformDistribution:
		ir.Name = LogUniformDistributionName
	case IntUniformDistribution:
		ir.Name = IntUniformDistributionName
	case DiscreteUniformDistribution:
		ir.Name = DiscreteUniformDistributionName
	case CategoricalDistribution:
		ir.Name = CategoricalDistributionName
	default:
		return nil, ErrUnknownDistribution
	}
	ir.Attrs = distribution
	return json.Marshal(&ir)
}

// JSONToDistribution deserialize a distribution in JSON format.
func JSONToDistribution(jsonBytes []byte) (interface{}, error) {
	var x struct {
		Name  string      `json:"name"`
		Attrs interface{} `json:"attributes"`
	}
	err := json.Unmarshal(jsonBytes, &x)
	if err != nil {
		return nil, err
	}
	switch x.Name {
	case UniformDistributionName:
		var y UniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case LogUniformDistributionName:
		var y LogUniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case IntUniformDistributionName:
		var y IntUniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case DiscreteUniformDistributionName:
		var y DiscreteUniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case CategoricalDistributionName:
		var y CategoricalDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	}
	return nil, ErrUnknownDistribution
}
