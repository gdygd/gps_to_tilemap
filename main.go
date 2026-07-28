package main

// ref : https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames

import (
	"fmt"
	"io"
	"math"
	"os"

	"gopkg.in/ini.v1"
)

type TileInfo struct {
	Lv int
	X1 int
	X2 int
	Y1 int
	Y2 int
}

var Bound1 [2]float64
var Bound2 [2]float64

var rootpath string = ""
var destroot string = ""

var stz int = 6
var endz int = 15

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	fmt.Printf("initEnvVaiable...\n")

	cfg, err := ini.Load("./conf.ini")
	if err != nil {
		fmt.Printf("fail to read sysenvini.ini %v", err)
		return false
	}

	rootpath = cfg.Section("MAPINFO").Key("srcpath").String()
	destroot = cfg.Section("MAPINFO").Key("destpath").String()

	Bound1[0], _ = cfg.Section("MAPINFO").Key("lat1").Float64()
	Bound1[1], _ = cfg.Section("MAPINFO").Key("lon1").Float64()

	Bound2[0], _ = cfg.Section("MAPINFO").Key("lat2").Float64()
	Bound2[1], _ = cfg.Section("MAPINFO").Key("lon2").Float64()

	stz, _ = cfg.Section("MAPINFO").Key("startlv").Int()
	endz, _ = cfg.Section("MAPINFO").Key("endlv").Int()

	return true
}

func tiletoLong(x, z float64) float64 {

	val := (x/math.Pow(2, z)*360 - 180)
	return val
}

func tileToLat(y, z float64) float64 {
	n := math.Pi - 2*math.Pi*y/math.Pow(2, z)
	val := (180 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n))))
	return val
}

func tileToLatLon(z, x, y int) (float64, float64) {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(z))
	lat := 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lon := float64(x)/math.Exp2(float64(z))*360.0 - 180.0

	return lat, lon

}

func degToTiles(lat, lon float64, z int) (int, int) {
	n := math.Exp2(float64(z))
	x := int(math.Floor((lon + 180.0) / 360.0 * n))
	if float64(x) >= n {
		x = int(n - 1)
	}
	y := int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n))
	return x, y
}

func initTileInfo() []TileInfo {
	var t []TileInfo = make([]TileInfo, 20)

	var x1 int = 52
	var x2 int = 57
	var y1 int = 22
	var y2 int = 26

	for z := 6; z <= endz; z++ {
		t[z] = TileInfo{Lv: z, X1: x1, X2: x2, Y1: y1, Y2: y2}

		x1 = x1 * 2
		x2 = x2 * 2

		y1 = y1*2 + 1
		y2 = y2*2 + 1
	}

	return t
}

func copyfile(src, dstpath, filenm string) error {
	if _, err := os.Stat(dstpath); os.IsNotExist(err) {
		//fmt.Printf("make dir..%v, %v, %v\n", src, dstpath, filenm)
		err = os.MkdirAll(dstpath, 0775)
		if err != nil {
			return err
		}
	}

	dst := fmt.Sprintf("%s/%s", dstpath, filenm)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}

	return nil

}

func checkRange(lat, lon float64) bool {
	fmt.Printf("lat [%f], [%f], [%f]\n", lat, Bound1[0], Bound2[0])
	fmt.Printf("lon [%f], [%f], [%f]\n", lon, Bound1[1], Bound2[1])
	if lat <= Bound1[0] && lat >= Bound2[0] &&
		lon >= Bound1[1] && lon <= Bound2[1] {

		return true
	}

	return false
}

// func main() {
/*
	전체 타일을 읽고, 타일마다 gps 좌표를 확인해서
	해당 gps bounds에 포함이 되는지 확인 (시간 오래걸림 모든 타일을 읽기 때문)
*/
// 	initok := initEnvVaiable()
// 	if !initok {
// 		fmt.Printf("Map info initEnvVaiable fail..")
// 		return
// 	}

// 	t := initTileInfo()

// 	for z := stz; z <= endz; z++ {
// 		x := t[z].X1
// 		x2 := t[z].X2

// 		for x <= x2 {
// 			y := t[z].Y1
// 			y2 := t[z].Y2

// 			for y <= y2 {
// 				filepath := fmt.Sprintf("%d/%d/%d.png", z, x, y)
// 				path := fmt.Sprintf("%s/%s", rootpath, filepath)
// 				fmt.Printf("path : [%s]\n", path)

// 				// lon1 := tiletoLong(float64(x), float64(z))
// 				// lat1 := tileToLat(float64(y), float64(z))
// 				lat, lon := tileToLatLon(z, x, y)
// 				lat1, lon1 := tileToLatLon(z, x+1, y+1)

// 				fmt.Printf("lat, lon :[%v, %v], [%v, %v]\n", lat, lon, lat1, lon1)

// 				if _, err := os.Stat(path); os.IsNotExist(err) {
// 					fmt.Printf("Unexist z/x/y file ..%d, %d, %d\n", z, x, y)
// 				} else {
// 					destpath := fmt.Sprintf("%s/%d/%d", destroot, z, x)
// 					filenm := fmt.Sprintf("%d.png", y)

// 					// check range
// 					chk1 := checkRange(lat, lon)
// 					chk2 := checkRange(lat1, lon1)

// 					// if checkRange(lat, lon) {
// 					// 	fmt.Printf("copy..file %d, %d, %d \n", z, x, y)
// 					// 	copyfile(path, destpath, filenm)
// 					// }
// 					if chk1 || chk2 {
// 						fmt.Printf("copy..file %d, %d, %d \n", z, x, y)
// 						copyfile(path, destpath, filenm)
// 					}
// 				}

// 				y++
// 			}

// 			x++
// 		}
// 	}
// }

type TileBound struct {
	xSt  int
	xEnd int

	ySt  int
	yEnd int
}

func min(n1, n2 int) int {

	if n1 > n2 {
		return n2
	}
	return n1
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func main() {
	/*
		// gps좌표 -> 타일번호 추출
		// z별로 X축에 해당하는 타일범위 Y축에 해당하는 타일 범위를 구함
		// 추출한 타일번호만 copy하기 때문에 속도 빠름
	*/

	var mpTileBound map[int]TileBound = make(map[int]TileBound)

	initok := initEnvVaiable()
	if !initok {
		fmt.Printf("Map info initEnvVaiable fail..\n")
		return
	}

	//t := initTileInfo()

	for z := stz; z <= endz; z++ {
		lat1 := Bound1[0]
		lon1 := Bound1[1]

		lat2 := Bound2[0]
		lon2 := Bound2[1]

		// X축
		Xx1, Xy1 := degToTiles(lat1, lon1, z)
		Xx2, Xy2 := degToTiles(lat1, lon2, z)

		fmt.Printf(" [%d, %d] - [%d, %d] \n", Xx1, Xy1, Xx2, Xy2)
		Xx1lat, Xx1lon := tileToLatLon(z, Xx1, Xy1)
		Xx2lat, Xx2lon := tileToLatLon(z, Xx2, Xy2)
		fmt.Printf(" [X] x1 latlon[%v, %v] - x2 latlon[%v, %v] \n", Xx1lat, Xx1lon, Xx2lat, Xx2lon)

		fmt.Println("")

		// Y축
		Yx1, Yy1 := degToTiles(lat1, lon1, z)
		Yx2, Yy2 := degToTiles(lat2, lon1, z)

		fmt.Printf(" [%d, %d] - [%d, %d] \n", Yx1, Yy1, Yx2, Yy2)

		Yx1lat, Yx1lon := tileToLatLon(z, Yx1, Yy1)
		Yx2lat, Yx2lon := tileToLatLon(z, Yx2, Yy2)
		fmt.Printf(" [Y] x1 latlon[%v, %v] - x2 latlon[%v, %v] \n", Yx1lat, Yx1lon, Yx2lat, Yx2lon)

		tileinfo := TileBound{}
		tileinfo.xSt = min(Xx1, Xx2)
		tileinfo.xEnd = max(Xx1, Xx2)
		tileinfo.ySt = min(Yy1, Yy2)
		tileinfo.yEnd = max(Yy1, Yy2)
		mpTileBound[z] = tileinfo

		fmt.Printf(">> z:%d X : (%d - %d), Y : [%d - %d] \n", z, tileinfo.xSt, tileinfo.xEnd, tileinfo.ySt, tileinfo.yEnd)

	}

	// file copy
	for z := stz; z <= endz; z++ {
		t := mpTileBound[z]

		for x := t.xSt; x <= t.xEnd; x++ {

			for y := t.ySt; y <= t.yEnd; y++ {
				filepath := fmt.Sprintf("%d/%d/%d.png", z, x, y)
				path := fmt.Sprintf("%s/%s", rootpath, filepath)

				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Printf("Unexist z/x/y file ..%d, %d, %d\n", z, x, y)
				} else {

					destpath := fmt.Sprintf("%s/%d/%d", destroot, z, x)
					filenm := fmt.Sprintf("%d.png", y)
					copyfile(path, destpath, filenm)

					fmt.Printf("\t >> copy file %d, %d, %d \n", z, x, y)
				}
			}
		}
	}
}
