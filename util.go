package main

import (
	"errors"
	"math"
)

// constants used for vincenty's formulae, https://en.wikipedia.org/wiki/Vincenty%27s_formulae
const a = 6378137.0
const b = 6356752.314245
const f = 1 / 298.257223563
const tolerance = 1e-12 // tolerance for convergence
const maxIter = 20

// radius of Earth in meters
const r = 6378100

// https://en.wikipedia.org/wiki/Haversine_formula
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func haversineDistance(p1, p2 Point) float64 {
	piRad := math.Pi / 180
	lat1, lat2 := p1.Latitude*piRad, p2.Latitude*piRad
	long1, long2 := p1.Longitude*piRad, p2.Longitude*piRad
	h := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(long2-long1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func vincentyDistance(p1, p2 Point) (float64, error) {
	piRad := math.Pi / 180
	lat1, lat2 := p1.Latitude*piRad, p2.Latitude*piRad
	long1, long2 := p1.Longitude*piRad, p2.Longitude*piRad
	L := long2 - long1
	U1 := math.Atan((1 - f) * math.Tan(lat1))
	U2 := math.Atan((1 - f) * math.Tan(lat2))
	lambda := L
	sinU1 := math.Sin(U1)
	cosU1 := math.Sin(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	var sinLambda, cosLambda, sigma, sinSigma, cosSigma, sinAlpha, cosAlpha, cos2SigmaM, C float64
	var iter = 0

	for {
		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)
		sinSigma = math.Sqrt(math.Pow(cosU2*sinLambda, 2) + math.Pow(cosU1*sinU2-sinU1*cosU2*cosLambda, 2))
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cos2SigmaM = cosSigma - (2*sinU1*sinU2)/(cosAlpha*cosAlpha)
		C = f / 16 * cosAlpha * cosAlpha * (4 + f*(4-3*cosAlpha*cosAlpha))
		lambdaPrev := lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
		if math.Abs(lambdaPrev-lambda) <= tolerance || iter == maxIter {
			break
		}
		iter++
	}
	if iter == maxIter {
		return 0, errors.New("vincenty failed to converge")
	}
	u2 := cosAlpha * cosAlpha * (a*a - b*b) / (b * b)
	A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
	B := u2 / 1024 * (256 + u2*(-128+u2*(74-47*u2)))
	deltaSigma := B*sinSigma*(cos2SigmaM) + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM))
	return b * A * (sigma - deltaSigma), nil
}
