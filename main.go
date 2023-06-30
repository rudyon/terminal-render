package main

import (
	"time"
	"math"
	"fmt"
	"golang.org/x/term"
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
				fmt.Print("#")
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
	cam_dist := 200.0
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

func rotateZ(vertex *Vector3, theta float64) {
	sin_theta := math.Sin(theta)
	cos_theta := math.Cos(theta)

	x := vertex.x * cos_theta - vertex.y * sin_theta
	y := vertex.y * cos_theta + vertex.x * sin_theta

	vertex.x = x
	vertex.y = y

}

func rotateX(vertex *Vector3, theta float64) {
	sin_theta := math.Sin(theta)
	cos_theta := math.Cos(theta)

	y := vertex.y * cos_theta - vertex.z * sin_theta
	z := vertex.z * cos_theta + vertex.y * sin_theta

	vertex.y = y
	vertex.z = z
}

func rotateY(vertex *Vector3, theta float64) {
	sin_theta := math.Sin(theta)
	cos_theta := math.Cos(theta)

	x := vertex.x * cos_theta + vertex.z * sin_theta
	z := vertex.z * cos_theta - vertex.x * sin_theta

	vertex.x = x
	vertex.z = z
}

func main() {
	width, height, err := getTerminalSize()
	if err != nil {
		fmt.Println("Failed to get terminal size:", err)
		return
	}

	tri := Triangle{Vector3{1, 1, 1},
					Vector3{-1, 1, 1},
					Vector3{-1, -1, 1}}
	tri2 := Triangle{Vector3{-1, -1, 1},
					Vector3{1, -1, 1},
					Vector3{1, 1, 1}}

	for {
		rotateZ(&tri.point1, .1)
		rotateZ(&tri.point2, .1)
		rotateZ(&tri.point3, .1)
		rotateZ(&tri2.point1, .1)
		rotateZ(&tri2.point2, .1)
		rotateZ(&tri2.point3, .1)
		
		canvas := initArray(width, height)
		projectTriangle(tri, &canvas)
		projectTriangle(tri2, &canvas)
		
		fmt.Printf("\x1bc") // clears screen
		drawCanvas(canvas)
		time.Sleep(100 * time.Millisecond)
	}
}
