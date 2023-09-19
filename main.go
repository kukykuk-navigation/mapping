package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
)

var number_of_childrens = 0
var number_of_valid_childrens = 0

func main() {

	// default map limitations

	var LatLimitMin float64 = -90
	var LatLimitMax float64 = -LatLimitMin
	var LonLimitMin float64 = -180
	var LonLimitMax float64 = -LonLimitMin

	fmt.Printf("default map limitations:\n")
	fmt.Printf("Lat: %f -> %f\nLon: %f -> %f\n", LatLimitMin, LatLimitMax, LonLimitMin, LonLimitMax)

	// region
	/*
		var Region ROI
		Region.LatMin = 49.9130800
		Region.LatMax = 50.1565375
		Region.LonMin = 14.2118994
		Region.LonMax = 14.7879950
		Region.ZoomLevel = 10
	*/

	var Region ROI
	Region.LatMin = 50.3586956
	Region.LatMax = 50.4762797
	Region.LonMin = 14.0964572
	Region.LonMax = 14.3188444
	Region.ZoomLevel = 12

	fmt.Printf("mapping region:\n")
	fmt.Printf("Lat: %f -> %f\nLon: %f -> %f\nZoom level: %v\n", Region.LatMin, Region.LatMax, Region.LonMin, Region.LonMax, Region.ZoomLevel)

	// create global map

	var Map Tile
	Map.X = 0
	Map.Y = 0
	Map.Lat = (LatLimitMin + LatLimitMax) / 2
	Map.Lon = (LonLimitMin + LonLimitMax) / 2
	Map.LatMin = LatLimitMin
	Map.LatMax = LatLimitMax
	Map.LonMin = LonLimitMin
	Map.LonMax = LonLimitMax
	Map.ZoomLevel = 0
	Map.Identifier = "00"
	Map.Parent = nil

	// generate global map

	Map.generateChildren(Region)

	fmt.Printf("total tiles: %v\n", number_of_childrens+1)
	fmt.Printf("total valid tiles: %v\n", number_of_valid_childrens+1)

}

type ROI struct {
	LatMin, LatMax, LonMin, LonMax float64
	ZoomLevel                      uint
}

type Tile struct {
	X, Y                           uint
	Lat, Lon                       float64
	LatMin, LatMax, LonMin, LonMax float64
	ZoomLevel                      uint
	VerticalDistance               float64
	HorizontalDistance             float64
	Identifier                     string
	Children                       [2][2]*Tile
	Parent                         *Tile
}

func NewTile() *Tile {
	t := Tile{}
	return &t
}

func (t *Tile) printTile() {
	fmt.Printf("%+v\n", t)
}

func (t *Tile) containsRegion(r ROI) bool {

	// if the whole region is in the tile

	if r.LatMin > t.LatMin && r.LatMax < t.LatMax && r.LonMin > t.LonMin && r.LonMax < t.LonMax {
		return true
	}

	// the whole tile is covered in the

	if r.LatMin < t.LatMin && r.LatMax > t.LatMax && r.LonMin < t.LonMin && r.LonMax > t.LonMax {
		return true
	}

	// boundary conditions

	if (r.LatMin > t.LatMin && r.LatMin < t.LatMax) || (r.LatMax > t.LatMin && r.LatMax < t.LatMax) {

		if (r.LonMin > t.LonMin && r.LonMin < t.LonMax) || (r.LonMax > t.LonMin && r.LonMax < t.LonMax) || (r.LonMin < t.LonMin && r.LonMax > t.LonMax) || (r.LonMin > t.LonMin && r.LonMax < t.LonMax) {

			return true

		}

	}
	if (r.LonMin > t.LonMin && r.LonMin < t.LonMax) || (r.LonMax > t.LonMin && r.LonMax < t.LonMax) {

		if (r.LatMin > t.LatMin && r.LatMin < t.LatMax) || (r.LatMax > t.LatMin && r.LatMax < t.LatMax) || (r.LatMin < t.LatMin && r.LatMax > t.LatMax) || (r.LatMin > t.LatMin && r.LatMax < t.LatMax) {

			return true

		}

	}

	// else
	return false

}

func (t *Tile) generateChildren(in_Region ROI) {

	if t.ZoomLevel == in_Region.ZoomLevel {

		number_of_valid_childrens = number_of_valid_childrens + 1

		t.printTile()

		//getCustom(t.LatMin, t.LonMin, t.LatMax, t.LonMax)
		//getOutdoors(t.LatMin, t.LonMin, t.LatMax, t.LonMax)
		getSatellite(t.Identifier, t.VerticalDistance/t.HorizontalDistance, t.LatMin, t.LonMin, t.LatMax, t.LonMax)
		//getStreets(t.LatMin, t.LonMin, t.LatMax, t.LonMax)

		return
	}

	halfLat := (t.LatMin + t.LatMax) / 2
	halfLon := (t.LonMin + t.LonMax) / 2

	// children [0,0]

	t.Children[0][0] = NewTile()

	t.Children[0][0].X = 0
	t.Children[0][0].Y = 0

	t.Children[0][0].LatMin = t.LatMin
	t.Children[0][0].LatMax = halfLat
	t.Children[0][0].LonMin = t.LonMin
	t.Children[0][0].LonMax = halfLon

	t.Children[0][0].Lat = (t.Children[0][0].LatMin + t.Children[0][0].LatMax) / 2
	t.Children[0][0].Lon = (t.Children[0][0].LonMin + t.Children[0][0].LonMax) / 2

	t.Children[0][0].VerticalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[0][0].LatMin, Longitude: t.Children[0][0].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[0][0].LatMax, Longitude: t.Children[0][0].LonMin, Altitude: 0})
	t.Children[0][0].HorizontalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[0][0].LatMin, Longitude: t.Children[0][0].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[0][0].LatMin, Longitude: t.Children[0][0].LonMax, Altitude: 0})

	t.Children[0][0].Identifier = t.Identifier + "_00"

	t.Children[0][0].ZoomLevel = t.ZoomLevel + 1
	t.Children[0][0].Parent = t

	// children [0,1]

	t.Children[0][1] = NewTile()

	t.Children[0][1].X = 0
	t.Children[0][1].Y = 1

	t.Children[0][1].LatMin = t.LatMin
	t.Children[0][1].LatMax = halfLat
	t.Children[0][1].LonMin = halfLon
	t.Children[0][1].LonMax = t.LonMax

	t.Children[0][1].Lat = (t.Children[0][1].LatMin + t.Children[0][1].LatMax) / 2
	t.Children[0][1].Lon = (t.Children[0][1].LonMin + t.Children[0][1].LonMax) / 2

	t.Children[0][1].VerticalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[0][1].LatMin, Longitude: t.Children[0][1].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[0][1].LatMax, Longitude: t.Children[0][1].LonMin, Altitude: 0})
	t.Children[0][1].HorizontalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[0][1].LatMin, Longitude: t.Children[0][1].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[0][1].LatMin, Longitude: t.Children[0][1].LonMax, Altitude: 0})

	t.Children[0][1].Identifier = t.Identifier + "_01"

	t.Children[0][1].ZoomLevel = t.ZoomLevel + 1
	t.Children[0][1].Parent = t

	// children [1,0]

	t.Children[1][0] = NewTile()

	t.Children[1][0].X = 1
	t.Children[1][0].Y = 0

	t.Children[1][0].LatMin = halfLat
	t.Children[1][0].LatMax = t.LatMax
	t.Children[1][0].LonMin = t.LonMin
	t.Children[1][0].LonMax = halfLon

	t.Children[1][0].Lat = (t.Children[1][0].LatMin + t.Children[1][0].LatMax) / 2
	t.Children[1][0].Lon = (t.Children[1][0].LonMin + t.Children[1][0].LonMax) / 2

	t.Children[1][0].VerticalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[1][0].LatMin, Longitude: t.Children[1][0].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[1][0].LatMax, Longitude: t.Children[1][0].LonMin, Altitude: 0})
	t.Children[1][0].HorizontalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[1][0].LatMin, Longitude: t.Children[1][0].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[1][0].LatMin, Longitude: t.Children[1][0].LonMax, Altitude: 0})

	t.Children[1][0].Identifier = t.Identifier + "_10"

	t.Children[1][0].ZoomLevel = t.ZoomLevel + 1
	t.Children[1][0].Parent = t

	// children [1,1]

	t.Children[1][1] = NewTile()

	t.Children[1][1].X = 1
	t.Children[1][1].Y = 1

	t.Children[1][1].LatMin = halfLat
	t.Children[1][1].LatMax = t.LatMax
	t.Children[1][1].LonMin = halfLon
	t.Children[1][1].LonMax = t.LonMax

	t.Children[1][1].Lat = (t.Children[1][1].LatMin + t.Children[1][1].LatMax) / 2
	t.Children[1][1].Lon = (t.Children[1][1].LonMin + t.Children[1][1].LonMax) / 2

	t.Children[1][1].VerticalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[1][1].LatMin, Longitude: t.Children[1][1].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[1][1].LatMax, Longitude: t.Children[1][1].LonMin, Altitude: 0})
	t.Children[1][1].HorizontalDistance = Navcom_WGS842DDistance(Navcom_WGS84Point{Latitude: t.Children[1][1].LatMin, Longitude: t.Children[1][1].LonMin, Altitude: 0}, Navcom_WGS84Point{Latitude: t.Children[1][1].LatMin, Longitude: t.Children[1][1].LonMax, Altitude: 0})

	t.Children[1][1].Identifier = t.Identifier + "_11"

	t.Children[1][1].ZoomLevel = t.ZoomLevel + 1
	t.Children[1][1].Parent = t

	// generate grand children

	if t.Children[0][0].containsRegion(in_Region) {
		number_of_childrens = number_of_childrens + 1
		t.Children[0][0].generateChildren(in_Region)
	}
	if t.Children[0][1].containsRegion(in_Region) {
		number_of_childrens = number_of_childrens + 1
		t.Children[0][1].generateChildren(in_Region)
	}
	if t.Children[1][0].containsRegion(in_Region) {
		number_of_childrens = number_of_childrens + 1
		t.Children[1][0].generateChildren(in_Region)
	}
	if t.Children[1][1].containsRegion(in_Region) {
		number_of_childrens = number_of_childrens + 1
		t.Children[1][1].generateChildren(in_Region)
	}

}

func wget(url, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getSatellite(identifier string, ratio, lat1, lon1, lat2, lon2 float64) {

	var p1 float64 = 512
	var p2 float64 = p1 * ratio

	wget(fmt.Sprintf("https://api.mapbox.com/styles/v1/mapbox/satellite-v9/static/[%2.4f,%2.4f,%2.4f,%2.4f]/%vx%v?access_token=pk.eyJ1Ijoia3VreWt1ay1uYXZpZ2F0aW9uIiwiYSI6ImNsZDNpbGx6YzBqMm8zcXByNnBqbGN2b2gifQ.OkOCNA_26exIi_GkEAJy5A&attribution=false&logo=false", math.Min(lon1, lon2), math.Min(lat1, lat2), math.Max(lon1, lon2), math.Max(lat1, lat2), int(p1), int(p2)), "./data/mapbox/satellite/satellite_"+identifier+".png")

}

func getOutdoors(lat1, lon1, lat2, lon2 float64) {

	wget(fmt.Sprintf("https://api.mapbox.com/styles/v1/mapbox/outdoors-v12/static/[%2.4f,%2.4f,%2.4f,%2.4f]/512x512?access_token=pk.eyJ1Ijoia3VreWt1ay1uYXZpZ2F0aW9uIiwiYSI6ImNsZDNpbGx6YzBqMm8zcXByNnBqbGN2b2gifQ.OkOCNA_26exIi_GkEAJy5A&attribution=false&logo=false", math.Min(lon1, lon2), math.Min(lat1, lat2), math.Max(lon1, lon2), math.Max(lat1, lat2)), fmt.Sprintf("./data/mapbox/outdoors/outdoors_%2.4f_%2.4f_%2.4f_%2.4f.png", lat1, lon1, lat2, lon2))

}

func getStreets(lat1, lon1, lat2, lon2 float64) {

	wget(fmt.Sprintf("https://api.mapbox.com/styles/v1/mapbox/streets-v12/static/[%2.4f,%2.4f,%2.4f,%2.4f]/512x512?access_token=pk.eyJ1Ijoia3VreWt1ay1uYXZpZ2F0aW9uIiwiYSI6ImNsZDNpbGx6YzBqMm8zcXByNnBqbGN2b2gifQ.OkOCNA_26exIi_GkEAJy5A&attribution=false&logo=false", math.Min(lon1, lon2), math.Min(lat1, lat2), math.Max(lon1, lon2), math.Max(lat1, lat2)), fmt.Sprintf("./data/mapbox/streets/streets_%2.4f_%2.4f_%2.4f_%2.4f.png", lat1, lon1, lat2, lon2))

}
func getCustom(lat1, lon1, lat2, lon2 float64) {

	wget(fmt.Sprintf("https://api.mapbox.com/styles/v1/kukykuk-navigation/clmfdfvcb01ha01pfht9ze4er/static/[%2.4f,%2.4f,%2.4f,%2.4f]/512x512?access_token=pk.eyJ1Ijoia3VreWt1ay1uYXZpZ2F0aW9uIiwiYSI6ImNsZDNpbGx6YzBqMm8zcXByNnBqbGN2b2gifQ.OkOCNA_26exIi_GkEAJy5A&attribution=false&logo=false", math.Min(lon1, lon2), math.Min(lat1, lat2), math.Max(lon1, lon2), math.Max(lat1, lat2)), fmt.Sprintf("./data/mapbox/custom/custom_%2.4f_%2.4f_%2.4f_%2.4f.png", lat1, lon1, lat2, lon2))

}
