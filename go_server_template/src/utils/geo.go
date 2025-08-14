package utils

import (
	"template/global"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func PointFormat(wkt string) string {
	if wkt == "" {
		return ""
	}
	format := strings.ReplaceAll(wkt, "POINT(", "")
	format = strings.ReplaceAll(format, " ", ",")
	format = strings.ReplaceAll(format, ")", "")
	return format
}

func Geom(field string) string {
	return fmt.Sprintf("st_setsrid(st_transform(st_setsrid(%s,4087),4326),0)", field)
}

func GeoBbox(field string) string {
	return fmt.Sprintf(`case
		when coalesce(%s,'') <> ''
	 then
		concat_ws(',',st_xmin(st_transform(st_makeenvelope(
		                        split_part(%s,',',1)::numeric,
		                        split_part(%s,',',2)::numeric,
		                        split_part(%s,',',3)::numeric,
		                        split_part(%s,',',4)::numeric ,
		                4087),4326)) ::text,  st_ymin(st_transform(st_makeenvelope(
		                        split_part(%s,',',1)::numeric,
		                        split_part(%s,',',2)::numeric,
		                        split_part(%s,',',3)::numeric,
		                        split_part(%s,',',4)::numeric ,
		                4087),4326)) ::text,  st_xmax(st_transform(st_makeenvelope(
		                        split_part(%s,',',1)::numeric,
		                        split_part(%s,',',2)::numeric,
		                        split_part(%s,',',3)::numeric,
		                        split_part(%s,',',4)::numeric ,
		                4087),4326)) ::text,  st_ymax(st_transform(st_makeenvelope(
		                        split_part(%s,',',1)::numeric,
		                        split_part(%s,',',2)::numeric,
		                        split_part(%s,',',3)::numeric,
		                        split_part(%s,',',4)::numeric ,
		                4087),4326)) ::text)
	 else
		''
	 end`, field, field, field, field, field, field, field, field, field, field, field, field, field, field, field, field, field)
}

func Distance(x1, y1, x2, y2 float64) float64 {
	distance := math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
	return float64(distance) / global.MeterDeg
}

func Distance2(point1, point2 string) float64 {
	var x1, y1, x2, y2 float64
	if point1 != "" {
		geo := strings.Split(point1, ",")
		x1, _ = strconv.ParseFloat(geo[0], 64)
		y1, _ = strconv.ParseFloat(geo[1], 64)
	}
	if point2 != "" {
		geo := strings.Split(point2, ",")
		x2, _ = strconv.ParseFloat(geo[0], 64)
		y2, _ = strconv.ParseFloat(geo[1], 64)
	}
	if x1 == 0 && y1 == 0 || x2 == 0 && y2 == 0 {
		return 0
	}
	return Distance(x1, y1, x2, y2)
}

func FormatGeom(str string) (float64, float64) {
	var x, y float64
	if str != "" && strings.Contains(str, ",") {
		geo := strings.Split(str, ",")
		x, _ = strconv.ParseFloat(geo[0], 64)
		y, _ = strconv.ParseFloat(geo[1], 64)
	}
	return x, y
}
