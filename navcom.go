package mapping

import (
	"math"
	"math/cmplx"

	"gocv.io/x/gocv"
)

// navigation - computing

func r2d(in float64) float64 {
	return in * (180 / math.Pi)
}

func d2r(in float64) float64 {
	return in * (math.Pi / 180)
}

type Navcom_NEDPoint struct {
	N float64
	E float64
	D float64
}

type Navcom_ECEFPoint struct {
	X float64
	Y float64
	Z float64
}

type Navcom_WGS84Point struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

func Navcom_WGS84ReferenceEllipse() (a, b, f, e, e2 float64) {

	var inverse_f float64

	inverse_f = 298.257223563

	a = 6378137.0
	f = float64(1) / inverse_f
	e = math.Sqrt(f * (float64(2) - f))
	e2 = f * (float64(2) - f)
	b = a * (float64(1) - f)

	return a, b, f, e, e2
}

func Navcom_normalRadius(lat float64) (nr float64) {

	lat = d2r(lat)

	a, _, _, _, e2 := Navcom_WGS84ReferenceEllipse()

	sinlat2 := math.Sin(lat) * math.Sin(lat)

	normal_radius := (a) / (math.Sqrt(float64(1) - (e2 * sinlat2)))

	return normal_radius
}

func Navcom_WGS842ECEF(in_point Navcom_WGS84Point) (out_point Navcom_ECEFPoint) {

	// convert degrees to radians

	lat_rad := d2r(in_point.Latitude)
	lon_rad := d2r(in_point.Longitude)

	// prepare constatns for computation

	_, _, _, _, e2 := Navcom_WGS84ReferenceEllipse()

	clat := math.Cos(lat_rad)
	clon := math.Cos(lon_rad)

	slat := math.Sin(lat_rad)
	slon := math.Sin(lon_rad)

	N := Navcom_normalRadius(in_point.Latitude)

	// compute ECEF: X,Y,Z

	out_point = Navcom_ECEFPoint{}

	out_point.X = (N + in_point.Altitude) * clat * clon
	out_point.Y = (N + in_point.Altitude) * clat * slon
	out_point.Z = (N*(float64(1)-e2) + in_point.Altitude) * slat

	return out_point
}

func Navcom_WGS84TangentPlaneRotationMatrix(in_point Navcom_WGS84Point) gocv.Mat {

	lat := d2r(in_point.Latitude)
	lon := d2r(in_point.Longitude)

	R := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV64F)

	R.SetDoubleAt(0, 0, -math.Sin(lat)*math.Cos(lon))
	R.SetDoubleAt(0, 1, -math.Sin(lat)*math.Sin(lon))
	R.SetDoubleAt(0, 2, math.Cos(lat))
	R.SetDoubleAt(1, 0, -math.Sin(lon))
	R.SetDoubleAt(1, 1, math.Cos(lon))
	R.SetDoubleAt(1, 2, 0)
	R.SetDoubleAt(2, 0, -math.Cos(lat)*math.Cos(lon))
	R.SetDoubleAt(2, 1, -math.Cos(lat)*math.Sin(lon))
	R.SetDoubleAt(2, 2, -math.Sin(lat))

	return R

}

func Navcom_WGS842NED(in_target, in_reference Navcom_WGS84Point) Navcom_NEDPoint {

	// convert WGS84 to ECEF

	target_ECEF := Navcom_WGS842ECEF(in_target)
	reference_ECEF := Navcom_WGS842ECEF(in_reference)

	// compute difference in ECEF

	dECEF := gocv.NewMatWithSize(3, 1, gocv.MatTypeCV64F)
	dECEF.SetDoubleAt(0, 0, target_ECEF.X-reference_ECEF.X)
	dECEF.SetDoubleAt(1, 0, target_ECEF.Y-reference_ECEF.Y)
	dECEF.SetDoubleAt(2, 0, target_ECEF.Z-reference_ECEF.Z)

	// compute R

	R := Navcom_WGS84TangentPlaneRotationMatrix(in_reference)

	// apply R to compute difference in NED

	dNED := R.MultiplyMatrix(dECEF)

	// create return value

	to_return := Navcom_NEDPoint{
		N: dNED.GetDoubleAt(0, 0),
		E: dNED.GetDoubleAt(1, 0),
		D: dNED.GetDoubleAt(2, 0),
	}

	// close matrices

	dECEF.Close()
	R.Close()
	dNED.Close()

	return to_return

}

func Navcom_ECEF3DDistance(in_p1, in_p2 Navcom_ECEFPoint) float64 {

	dECEF := gocv.NewMatWithSize(3, 1, gocv.MatTypeCV64F)
	defer dECEF.Close()

	dECEF.SetDoubleAt(0, 0, in_p1.X-in_p2.X)
	dECEF.SetDoubleAt(1, 0, in_p1.Y-in_p2.Y)
	dECEF.SetDoubleAt(2, 0, in_p1.Z-in_p2.Z)

	// compute distance

	return math.Sqrt(math.Pow(dECEF.GetDoubleAt(0, 0), 2) + math.Pow(dECEF.GetDoubleAt(0, 1), 2) + math.Pow(dECEF.GetDoubleAt(0, 2), 2))

}

func Navcom_NED3DDistance(in_p1, in_p2 Navcom_NEDPoint) float64 {

	dNED := gocv.NewMatWithSize(3, 1, gocv.MatTypeCV64F)
	defer dNED.Close()

	dNED.SetDoubleAt(0, 0, in_p1.N-in_p2.N)
	dNED.SetDoubleAt(1, 0, in_p1.E-in_p2.E)
	dNED.SetDoubleAt(2, 0, in_p1.D-in_p2.D)

	// compute distance

	return math.Sqrt(math.Pow(dNED.GetDoubleAt(0, 0), 2) + math.Pow(dNED.GetDoubleAt(0, 1), 2) + math.Pow(dNED.GetDoubleAt(0, 2), 2))

}

func Navcom_WGS842DDistance(in_target, in_reference Navcom_WGS84Point) float64 {

	// convert WGS84 to ECEF

	target_ECEF := Navcom_WGS842ECEF(in_target)
	reference_ECEF := Navcom_WGS842ECEF(in_reference)

	// compute difference in ECEF

	dECEF := gocv.NewMatWithSize(3, 1, gocv.MatTypeCV64F)
	dECEF.SetDoubleAt(0, 0, target_ECEF.X-reference_ECEF.X)
	dECEF.SetDoubleAt(1, 0, target_ECEF.Y-reference_ECEF.Y)
	dECEF.SetDoubleAt(2, 0, target_ECEF.Z-reference_ECEF.Z)

	// compute R

	R := Navcom_WGS84TangentPlaneRotationMatrix(in_reference)

	// apply R to compute difference in NED

	dNED := R.MultiplyMatrix(dECEF)

	return math.Sqrt(math.Pow(dNED.GetDoubleAt(0, 0), 2) + math.Pow(dNED.GetDoubleAt(0, 1), 2))

}

func Navcom_relativeHeading(in_target, in_reference float64) float64 {

	f64h1rad := d2r(in_target)
	f64h2rad := d2r(in_reference)

	c128h1Reverse := cmplx.Rect(1, -f64h1rad)
	c128h2 := cmplx.Rect(1, f64h2rad)
	c128Diff := c128h1Reverse * c128h2

	return r2d(cmplx.Phase(c128Diff))
}

func Navcom_normalizeCourse(in_course float64) float64 {

	complex_course := cmplx.Rect(1, d2r(in_course))
	f := r2d(cmplx.Phase(complex_course))
	if f < 0 {
		return 360 + f
	} else {
		return f
	}

}
