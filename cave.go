package main

import (
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	width   = 100
	height  = 100
	random  = 50
	smooth  = 5
	seed    = uint32(0)
	seedSet = "abcdefghijklmnopqrstuvwxyz"
	ground  = color.RGBA{
		R: 35,
		G: 35,
		B: 35,
		A: 255,
	}
	wall = color.RGBA{
		R: 60,
		G: 60,
		B: 60,
		A: 255,
	}
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "help" {
		fmt.Println("Usage: cave width height random smooth seed")
		return
	}

	readArgs(args)
	rand.Seed(int64(seed))

	cave := makeMap()
	for i := 0; i < smooth; i++ {
		smoothMap(cave)
	}

	makeImage(cave)

}

func smoothMap(cave [][]int) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			c := countNeighbours(cave, i, j)
			if c > 4 {
				cave[i][j] = 1
			}
			if c < 4 {
				cave[i][j] = 0
			}
		}
	}
}

func countNeighbours(cave [][]int, x, y int) int {
	l := x - 1
	r := x + 1
	u := y + 1
	d := y - 1

	c := 0

	if x == 0 {
		l = width - 1
	}
	if x == width-1 {
		r = 0
	}
	if y == 0 {
		d = height - 1
	}
	if y == height-1 {
		u = 0
	}

	if cave[l][u] == 1 {
		c++
	}
	if cave[l][d] == 1 {
		c++
	}
	if cave[l][y] == 1 {
		c++
	}
	if cave[r][u] == 1 {
		c++
	}
	if cave[r][d] == 1 {
		c++
	}
	if cave[r][y] == 1 {
		c++
	}
	if cave[x][u] == 1 {
		c++
	}
	if cave[x][d] == 1 {
		c++
	}

	return c
}

func makeMap() [][]int {
	cave := make([][]int, width)
	for i := 0; i < width; i++ {
		cave[i] = make([]int, height)
	}

	/*
		//With wall
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				if i == 0 || i == width-1 || j == 0 || j == height-1 {
					cave[i][j] = 1
				} else {
					if (rand.Int() % 100) > random {
						cave[i][j] = 1
					} else {
						cave[i][j] = 0
					}
				}
			}
		}
	*/

	//Without wall
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if (rand.Int() % 100) > random {
				cave[i][j] = 1
			} else {
				cave[i][j] = 0
			}
		}
	}

	return cave
}

func makeImage(cave [][]int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if cave[i][j] == 1 {
				img.Set(i, j, wall)
			} else {
				img.Set(i, j, ground)
			}

		}
	}

	f, err := os.Create("cave.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func readArgs(args []string) {
	var err error
	if len(args) > 0 {
		width, err = strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if len(args) > 1 {
		height, err = strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if len(args) > 2 {
		random, err = strconv.Atoi(args[2])
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if len(args) > 3 {
		smooth, err = strconv.Atoi(args[3])
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if len(args) > 4 {
		h := fnv.New32a()
		h.Write([]byte(args[4]))
		seed = h.Sum32()
	} else {
		h := fnv.New32a()
		h.Write(randomSeed())
		seed = h.Sum32()
	}
}

func randomSeed() []byte {
	rand.Seed(time.Now().Unix())

	size := (rand.Int() % 10) + 10
	seed := make([]byte, size)

	for i := range seed {
		seed[i] = seedSet[rand.Int()%len(seedSet)]
	}
	return seed
}
