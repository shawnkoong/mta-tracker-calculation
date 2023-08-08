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
const r = 637810000.0

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
	//return 2 * r * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
}

// something wrong with this function not returning the correct value
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

	var sinLambda, cosLambda, sigma, sinSigma, cosSigma, sinAlpha, cosAlpha, cosSqAlpha, cos2SigmaM, C float64
	var iter = 0

	for {
		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)
		sinSigma = math.Sqrt(math.Pow(cosU2*sinLambda, 2) + math.Pow(cosU1*sinU2-sinU1*cosU2*cosLambda, 2))
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - (2*sinU1*sinU2)/cosSqAlpha
		C = f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))
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

func calculateDistance(p1, p2 Point) float64 {
	distance, err := vincentyDistance(p1, p2)
	if err != nil {
		return haversineDistance(p1, p2)
	}
	return distance
}

func calculateBearing(p1, p2 Point) float64 {
	piRad := math.Pi / 180
	lat1, lat2 := p1.Latitude*piRad, p2.Latitude*piRad
	long1, long2 := p1.Longitude*piRad, p2.Longitude*piRad
	y := math.Sin(long2-long1) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(long2-long1)
	return math.Atan2(y, x)
}

func pointWithProgress(point Point, distance float64, bearing float64) Point {
	piRad := math.Pi / 180
	lat1 := point.Latitude * piRad
	long1 := point.Longitude * piRad
	dr := distance / r
	lat2 := math.Asin(math.Sin(lat1)*math.Cos(dr) + math.Cos(lat1)*math.Sin(dr)*math.Cos(bearing))
	long2 := long1 + math.Atan2(math.Sin(bearing)*math.Sin(dr)*math.Cos(lat1), math.Cos(dr)-math.Sin(lat1)*math.Sin(lat2))
	return Point{Latitude: lat2 / piRad, Longitude: long2 / piRad}
}

func normalizeBearing(bearing float64) float64 {
	if bearing < 0 {
		bearing = bearing + 2*math.Pi
	}
	return bearing
}
