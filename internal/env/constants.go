package env

import "math"

func PhysicsConstants() map[string]float64 {
	return map[string]float64{
		"c":     299792458,
		"g":     9.80665,
		"h":     6.62607015e-34,
		"k_b":   1.380649e-23,
		"N_A":   6.02214076e23,
		"R":     8.314462618,
		"sigma": 5.670374419e-8,
	}
}

func MathConstants() map[string]float64 {
	return map[string]float64{
		"pi":    math.Pi,
		"e":     math.E,
		"tau":   2 * math.Pi,
		"phi":   (1 + math.Sqrt(5)) / 2,
		"sqrt2": math.Sqrt2,
		"sqrt3": math.Sqrt(3),
		"ln2":   math.Ln2,
		"ln10":  math.Log(10),
		"inf":   math.Inf(1),
		"nan":   math.NaN(),
	}
}

func UnitConversions() map[string]float64 {
	return map[string]float64{
		"deg_to_rad": math.Pi / 180,
		"rad_to_deg": 180 / math.Pi,
		"km_to_mi":   0.621371,
		"mi_to_km":   1.60934,
		"lb_to_kg":   0.453592,
		"kg_to_lb":   2.20462,
		"ft_to_m":    0.3048,
		"m_to_ft":    3.28084,
		"in_to_cm":   2.54,
		"cm_to_in":   0.393701,
		"gal_to_l":   3.78541,
		"l_to_gal":   0.264172,
	}
}

func LoadConstants(e *Env, setName string) {
	switch setName {
	case "math":
		e.SetAll(MathConstants())
	case "physics":
		e.SetAll(PhysicsConstants())
	case "units":
		e.SetAll(UnitConversions())
	}
}

func AllConstantSets() []string {
	return []string{"math", "physics", "units"}
}
