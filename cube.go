package main

import (
	"fmt"
	"math"
	"time"
)

var A, B, C float64

const width = 160
const height = 44

var zBuffer [width * height]float64
var buffer [width * height]byte
var bgCode byte = ' '
var distanceFromCam float64 = 100
var horizontalOffset float64
var K1 float64 = 40
var incrementSpeed float64 = 0.8

func calculateX(i, j, k float64) float64 {

	return 0 +
		math.Floor(j)*math.Sin(A)*math.Sin(B)*math.Cos(C) -
		math.Floor(k)*math.Cos(A)*math.Sin(B)*math.Cos(C) +
		math.Floor(j)*math.Cos(A)*math.Sin(C) +
		math.Floor(k)*math.Sin(A)*math.Sin(C) +
		math.Floor(i)*math.Cos(B)*math.Cos(C)
}
func calculateY(i, j, k float64) float64 {
	return 0 +
		math.Floor(j)*math.Cos(A)*math.Cos(C) +
		math.Floor(k)*math.Sin(A)*math.Cos(C) -
		math.Floor(j)*math.Sin(A)*math.Sin(B)*math.Sin(C) +
		math.Floor(k)*math.Cos(A)*math.Sin(B)*math.Sin(C) -
		math.Floor(i)*math.Cos(B)*math.Sin(C)
}
func calculateZ(i, j, k float64) float64 {
	return 0 +
		math.Floor(k)*math.Cos(A)*math.Cos(B) -
		math.Floor(j)*math.Sin(A)*math.Cos(B) +
		math.Floor(i)*math.Sin(B)
}

func calculateForSurface(horizontalOffset, cubeX, cubeY, cubeZ float64, ch byte) {
	x := calculateX(cubeX, cubeY, cubeZ)
	y := calculateY(cubeX, cubeY, cubeZ)
	z := calculateZ(cubeX, cubeY, cubeZ) + distanceFromCam
	ooz := 1 / z
	xp := width/2 + int(horizontalOffset+K1*ooz*x*2)
	yp := height/2 + int(K1*ooz*y)
	idx := xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
		}
	}
}
func calculateForAllSurfaces(cubeWidth, horizontalOffset float64) {
	for cubeX := -cubeWidth; cubeX < cubeWidth; cubeX += incrementSpeed {
		for cubeY := -cubeWidth; cubeY < cubeWidth; cubeY += incrementSpeed {
			calculateForSurface(horizontalOffset, cubeX, cubeY, -cubeWidth, '$')
			calculateForSurface(horizontalOffset, cubeWidth, cubeY, cubeX, '.')
			calculateForSurface(horizontalOffset, -cubeWidth, cubeY, -cubeX, ':')
			calculateForSurface(horizontalOffset, -cubeX, cubeY, cubeWidth, ';')
			calculateForSurface(horizontalOffset, cubeX, -cubeWidth, -cubeY, '#')
			calculateForSurface(horizontalOffset, cubeX, cubeWidth, cubeY, '*')
		}
	}
}

func fillArray[T comparable](a []T, value T) {
	for i := range a {
		a[i] = value
	}
}
func main() {
	// Clear screen
	fmt.Println("\x1b[2J")
	var cubeWidth float64
	for {
		fillArray(buffer[:], bgCode)
		fillArray(zBuffer[:], 0)

		cubeWidth = 20.0
		horizontalOffset = -2 * cubeWidth
		calculateForAllSurfaces(cubeWidth, horizontalOffset)

		// Draw the cube
		fmt.Println("\x1b[H")
		for k := 0; k < width*height; k++ {
			fmt.Printf("%c", buffer[k])
			if k%width == width-1 {
				fmt.Println()
			}
		}

		A += 0.01
		B += 0.08
		C += 0.02

		time.Sleep(time.Millisecond * 8)
	}
}
