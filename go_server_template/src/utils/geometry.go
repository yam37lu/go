package utils

import "fmt"

var (
	CharSet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func CalculatePolygonGravityCenter(coordinates [][]float64) (float64, float64, error) {
	area := 0.0
	gravityLat := 0.0
	gravityLng := 0.0
	for _, coordinate := range coordinates {
		if len(coordinate) != 2 {
			return 0.0, 0.0, fmt.Errorf("coordinate[%v] format error", coordinate)
		}
	}
	for index, coordinate := range coordinates {
		lng := coordinate[0]
		lat := coordinate[1]
		nextLng := coordinates[(index+1)%len(coordinates)][0]
		nextLat := coordinates[(index+1)%len(coordinates)][1]

		tempArea := (nextLat*lng - nextLng*lat) / 2.0
		area += tempArea

		gravityLng += tempArea * (lng + nextLng) / 3
		gravityLat += tempArea * (lat + nextLat) / 3
	}
	return gravityLng / area, gravityLat / area, nil
}

func MapSubDivisionID(coordinates [][]float64) (string, error) {
	lng, lat, err := CalculatePolygonGravityCenter(coordinates)
	if err != nil {
		return "", err
	}
	// lng := 114.5625
	// lat := 39.375

	a := int(lat/4.0) + 1
	b := int(lng/6.0) + 31
	minuteLat := int((lat - float64(int(lat))) * 60.0)
	//secondLat := int(((lat-float64(int(lat)))*60.0 - float64(minuteLat)) * 60)
	c := 48 - (int(lat)*60+minuteLat)%240/5
	minuteLng := int((lng - float64(int(lng))) * 60.0)
	secondLng := int(((lng-float64(int(lng)))*60.0 - float64(minuteLng)) * 60)
	d := (int(lng)*3600+minuteLng*60+secondLng)%21600/450 + 1

	return fmt.Sprintf("%s%dF%03d%03d", CharSet[a-1], b, c, d), nil
}
