package main

import (
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
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
	scale := 200.0 / (vertex.z + 200.0)
	projectedX := (vertex.x * scale) + (width / 2)
	projectedY := (vertex.y * scale) + (height /2)

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

func main() {
	width, height, err := getTerminalSize()
	if err != nil {
		fmt.Println("Failed to get terminal size:", err)
		return
	}

	fmt.Println("Terminal size:")
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	canvas := initArray(width, height)

	tri := Triangle{Vector3{20, 20, 20},
					Vector3{-20, 20, 20},
					Vector3{-20, -20, 20}}
	
	projectTriangle(tri, &canvas)
	
	drawCanvas(canvas)
}
