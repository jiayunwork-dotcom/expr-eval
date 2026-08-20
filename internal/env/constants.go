package env

import "math"

// Predefined constant sets that can be loaded into an environment.

// PhysicsConstants returns common physical constants.
func PhysicsConstants() map[string]float64 {
	return map[string]float64{
		"c":    299792458,     // speed of light (m/s)
		"g":    9.80665,       // standard gravity (m/s^2)
		"h":    6.62607015e-34, // Planck constant (J*s)
		"k_b":  1.380649e-23,  // Boltzmann constant (J/K)
		"N_A":  6.02214076e23, // Avogadro number
		"R":    8.314462618,   // gas constant (J/(mol*K))
		"sigma": 5.670374419e-8, // Stefan-Boltzmann (W/(m^2*K^4))
	}
}

// MathConstants returns extended mathematical constants.
func MathConstants() map[string]float64 {
	return map[string]float64{
		"pi":     math.Pi,
		"e":      math.E,
		"tau":    2 * math.Pi,
		"phi":    (1 + math.Sqrt(5)) / 2, // golden ratio
		"sqrt2":  math.Sqrt2,
		"sqrt3":  math.Sqrt(3),
		"ln2":    math.Ln2,
		"ln10":   math.Log(10),
		"inf":    math.Inf(1),
		"nan":    math.NaN(),
	}
}

// UnitConversions returns common unit conversion factors.
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

// LoadConstants loads a named constant set into the environment.
// Supported sets: "math", "physics", "units".
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

// AllConstantSets returns all available constant set names.
func AllConstantSets() []string {
	return []string{"math", "physics", "units"}
}
