package database

import "math"

// calculateExoplanetMass рассчитывает массу экзопланеты по формуле радиальной скорости
func calculateExoplanetMass(starMass, velocityAmplitude, orbitalPeriod, inclination float64) float64 {
	const G = 6.67430e-11
	const solarMass = 1.989e30
	const dayInSeconds = 86400.0
	const jupiterMass = 1.898e27

	starMassKg := starMass * solarMass
	periodSeconds := orbitalPeriod * dayInSeconds
	inclinationRad := inclination * math.Pi / 180.0

	planetMass := (starMassKg * velocityAmplitude * math.Pow(periodSeconds, 1.0/3.0)) /
		(math.Pow(2*math.Pi*G, 2.0/3.0) * math.Sin(inclinationRad))

	planetMassJupiter := planetMass / jupiterMass

	if planetMassJupiter < 0.001 {
		planetMassJupiter = 0.001
	}
	if planetMassJupiter > 100 {
		planetMassJupiter = 100
	}
	return planetMassJupiter
}
