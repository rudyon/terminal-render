package main

import (
	"fmt"
	"golang.org/x/term"
	"math"
	"time"
)

type Vector2 struct {
	x float64
	y float64
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Triangle struct {
	point1 Vector3
	point2 Vector3
	point3 Vector3
}

func getTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(0)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func initArray(rows, cols int) [][]int {
	arr := make([][]int, rows)

	for i := 0; i < rows; i++ {
		arr[i] = make([]int, cols)
	}

	return arr
}

func drawCanvas(canvas [][]int) {
	for i := 0; i < len(canvas[0]); i++ {
		for j := 0; j < len(canvas); j++ {
			if canvas[j][i] == 0 {
				fmt.Print(" ")
			} else if canvas[j][i] == 1 {
				fmt.Print("█")
			}
		}
		fmt.Print("\n")
	}
}

func drawLine(x0, y0, x1, y1 int, canvas *[][]int) {
	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))
	sx := -1
	sy := -1

	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}

	err := dx - dy

	for {
		(*canvas)[y0][x0] = 1

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func project(vertex Vector3, width, height float64) Vector2 {
	fov := 1200.0
	cam_dist := 100.0
	scale := fov / (vertex.z + cam_dist)

	projectedX := (vertex.x * scale) + (width / 2)
	projectedY := (vertex.y * scale) + (height / 2)

	return Vector2{projectedX, projectedY}
}

func projectTriangle(triangle Triangle, canvas *[][]int) {
	projectedPoint1 := project(triangle.point1, float64(len((*canvas)[0])), float64(len(*canvas)))
	projectedPoint2 := project(triangle.point2, float64(len((*canvas)[0])), float64(len(*canvas)))
	projectedPoint3 := project(triangle.point3, float64(len((*canvas)[0])), float64(len(*canvas)))

	(*canvas)[int(projectedPoint1.y)][int(projectedPoint1.x)] = 1
	(*canvas)[int(projectedPoint2.y)][int(projectedPoint2.x)] = 1
	(*canvas)[int(projectedPoint3.y)][int(projectedPoint3.x)] = 1

	drawLine(int(projectedPoint1.x), int(projectedPoint1.y), int(projectedPoint2.x), int(projectedPoint2.y), canvas)
	drawLine(int(projectedPoint2.x), int(projectedPoint2.y), int(projectedPoint3.x), int(projectedPoint3.y), canvas)
	drawLine(int(projectedPoint3.x), int(projectedPoint3.y), int(projectedPoint1.x), int(projectedPoint1.y), canvas)
}

func projectModel(model []Vector3, canvas *[][]int) {
	for i := 0; i+2 < len(model); i += 3 {
		triangle := Triangle{
			point1: model[i],
			point2: model[i+1],
			point3: model[i+2],
		}
		projectTriangle(triangle, canvas)
	}
}

func rotateModel(model *[]Vector3, x, y, z float64) {
	for i := 0; i < len(*model); i++ {
		rotateX(&((*model)[i]), x)
		rotateY(&((*model)[i]), y)
		rotateZ(&((*model)[i]), z)
	}
}

func rotateZ(vertex *Vector3, theta float64) {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	x := vertex.x
	y := vertex.y

	vertex.x = x*cosTheta - y*sinTheta
	vertex.y = y*cosTheta + x*sinTheta
}

func rotateX(vertex *Vector3, theta float64) {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	y := vertex.y
	z := vertex.z

	vertex.y = y*cosTheta - z*sinTheta
	vertex.z = z*cosTheta + y*sinTheta
}

func rotateY(vertex *Vector3, theta float64) {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	x := vertex.x
	z := vertex.z

	vertex.x = x*cosTheta + z*sinTheta
	vertex.z = z*cosTheta - x*sinTheta
}

func main() {
	width, height, err := getTerminalSize()
	if err != nil {
		fmt.Println("Failed to get terminal size:", err)
		return
	}

	
	
	model := []Vector3{
		// Front face
		{-1, -1, 1},   // Front bottom left
		{1, -1, 1},    // Front bottom right
		{-1, 1, 1},    // Front top left
		{1, 1, 1},     // Front top right
		{-1, 1, 1},    // Front top left
		{1, -1, 1},    // Front bottom right

		// Back face
		{-1, -1, -1},  // Back bottom left
		{-1, 1, -1},   // Back top left
		{1, -1, -1},   // Back bottom right
		{-1, 1, -1},   // Back top left
		{1, 1, -1},    // Back top right
		{1, -1, -1},   // Back bottom right

		// Left face
		{-1, -1, -1},  // Back bottom left
		{-1, 1, -1},   // Back top left
		{-1, -1, 1},   // Front bottom left
		{-1, 1, -1},   // Back top left
		{-1, 1, 1},    // Front top left
		{-1, -1, 1},   // Front bottom left

		// Right face
		{1, -1, 1},    // Front bottom right
		{1, 1, 1},     // Front top right
		{1, -1, -1},   // Back bottom right
		{1, 1, 1},     // Front top right
		{1, 1, -1},    // Back top right
		{1, -1, -1},   // Back bottom right

		// Top face
		{-1, 1, 1},    // Front top left
		{1, 1, 1},     // Front top right
		{-1, 1, -1},   // Back top left
		{1, 1, 1},     // Front top right
		{1, 1, -1},    // Back top right
		{-1, 1, -1},   // Back top left

		// Bottom face
		{-1, -1, 1},   // Front bottom left
		{-1, -1, -1},  // Back bottom left
		{1, -1, 1},    // Front bottom right
		{-1, -1, -1},  // Back bottom left
		{1, -1, -1},   // Back bottom right
		{1, -1, 1},    // Front bottom right
	}

	for {
		rotateModel(&model, 0.1, 0.1, 0.1)

		canvas := initArray(width, height)
		projectModel(model, &canvas)

		fmt.Printf("\x1bc") // clears screen
		drawCanvas(canvas)
		time.Sleep(100 * time.Millisecond)
	}
}
