package main

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

// var Bound1 = [2]float64{37.4966612, 126.8259632}
// var Bound2 = [2]float64{37.2095794, 127.1892489}

var Bound1 [2]float64
var Bound2 [2]float64

var rootpath string = ""
var destroot string = ""

var stz int = 6
var endz int = 15

// var stz int = 19
// var endz int = 19

// var stz int
// var endz int

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	fmt.Printf("initEnvVaiable...")

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

	// fmt.Printf("Bound1 %v", Bound1)
	// fmt.Printf("Bound2 %v", Bound2)

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
		fmt.Printf("make dir..%v, %v, %v\n", src, dstpath, filenm)
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

func main() {
	initok := initEnvVaiable()
	if !initok {
		fmt.Printf("Map info initEnvVaiable fail..")
		return
	}

	t := initTileInfo()

	for z := stz; z <= endz; z++ {
		x := t[z].X1
		x2 := t[z].X2

		for x <= x2 {
			y := t[z].Y1
			y2 := t[z].Y2

			for y <= y2 {
				filepath := fmt.Sprintf("%d/%d/%d.png", z, x, y)
				path := fmt.Sprintf("%s/%s", rootpath, filepath)
				fmt.Printf("path : [%s]\n", path)

				lon := tiletoLong(float64(x), float64(z))
				lat := tileToLat(float64(y), float64(z))
				fmt.Printf("lat, lon :%v, %v\n", lat, lon)

				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Printf("Unexist z/x/y file ..%d, %d, %d\n", z, x, y)
				} else {
					destpath := fmt.Sprintf("%s/%d/%d", destroot, z, x)
					filenm := fmt.Sprintf("%d.png", y)

					// check range
					if checkRange(lat, lon) {
						fmt.Printf("copy..file %d, %d, %d \n", z, x, y)
						copyfile(path, destpath, filenm)
					}
				}

				y++
			}

			x++
		}
	}
}
