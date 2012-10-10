package main

import dbg "fmt"
import "flag"
import "fmt"
import "image"
import "log"
import "os"
import "path"
import "strings"

import "github.com/0xC3/progress/barcli"
import "github.com/mewrnd/blizzconv/configs/min"
import "github.com/mewrnd/blizzconv/configs/til"
import "github.com/mewrnd/blizzconv/images/cel"
import "github.com/mewrnd/blizzconv/images/imgconf"
import "github.com/mewrnd/blizzconv/mpq"
import "github.com/mewrnd/pkg/imgutil"

func init() {
	flag.Usage = usage
	flag.StringVar(&imgconf.IniPath, "celini", "cel.ini", "Path to an ini file containing image information.")
	flag.StringVar(&mpq.ExtractPath, "mpqdump", "mpqdump/", "Path to an extracted MPQ file.")
	flag.StringVar(&mpq.IniPath, "mpqini", "mpq.ini", "Path to an ini file containing relative path information.")
	flag.Parse()
	err := mpq.Init()
	if err != nil {
		log.Fatalln(err)
	}
	err = imgconf.Init()
	if err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [name.til]...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	for _, tilName := range flag.Args() {
		err := tilDump(tilName)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// bar represents the progress bar.
var bar *barcli.Bar

// dumpPrefix is the name of the dump directory.
const dumpPrefix = "_dump_/"

// tilDump creates a dump directory and dumps the TIL file's squares using the
// pillars constructed based on the MIN format, once for each image config
// (pal).
func tilDump(tilName string) (err error) {
	squares, err := til.Parse(tilName)
	if err != nil {
		return err
	}
	nameWithoutExt := tilName[:len(tilName)-len(path.Ext(tilName))]
	minName := nameWithoutExt + ".min"
	pillars, err := min.Parse(minName)
	if err != nil {
		return err
	}
	imgName := nameWithoutExt + ".cel"
	relPalPaths := imgconf.GetRelPalPaths(imgName)
	for _, relPalPath := range relPalPaths {
		conf, err := cel.GetConf(imgName, relPalPath)
		if err != nil {
			return err
		}
		var palDir string
		if len(relPalPaths) > 1 {
			dbg.Println("using pal:", relPalPath)
			palDir = path.Base(relPalPath) + "/"
		}
		bar, err = barcli.New(len(squares))
		if err != nil {
			return err
		}
		levelFrames, err := cel.DecodeAll(imgName, conf)
		if err != nil {
			return err
		}
		dumpDir := path.Clean(dumpPrefix+"_squares_/"+nameWithoutExt) + "/" + palDir
		// prevent directory traversal
		if !strings.HasPrefix(dumpDir, dumpPrefix) {
			return fmt.Errorf("path (%s) contains no dump prefix (%s).", dumpDir, dumpPrefix)
		}
		err = os.MkdirAll(dumpDir, 0755)
		if err != nil {
			return err
		}
		err = dumpSquares(squares, pillars, levelFrames, dumpDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// dumpPillars stores each pillar as a new png image, using the frames from a
// CEL image level file.
func dumpSquares(squares []til.Square, pillars []min.Pillar, levelFrames []image.Image, dumpDir string) (err error) {
	for squareNum, square := range squares {
		squarePath := dumpDir + fmt.Sprintf("square_%04d.png", squareNum)
		bar.Inc()
		img := square.Image(pillars, levelFrames)
		err = imgutil.WriteFile(squarePath, img)
		if err != nil {
			return err
		}
	}
	return nil
}
