package main

import (
	"fmt"
	"golang.org/x/term"
	"math"
	"time"
	"github.com/hschendel/stl"
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
    width := len(*canvas)
    height := len((*canvas)[0])

    projectedPoint1 := project(triangle.point1, float64(len((*canvas)[0])), float64(len(*canvas)))
	projectedPoint2 := project(triangle.point2, float64(len((*canvas)[0])), float64(len(*canvas)))
	projectedPoint3 := project(triangle.point3, float64(len((*canvas)[0])), float64(len(*canvas)))

    // Check bounds before accessing the canvas
    if int(projectedPoint1.x) >= 0 && int(projectedPoint1.x) < width &&
        int(projectedPoint1.y) >= 0 && int(projectedPoint1.y) < height {
        (*canvas)[int(projectedPoint1.y)][int(projectedPoint1.x)] = 1
    }
    
    if int(projectedPoint2.x) >= 0 && int(projectedPoint2.x) < width &&
        int(projectedPoint2.y) >= 0 && int(projectedPoint2.y) < height {
        (*canvas)[int(projectedPoint2.y)][int(projectedPoint2.x)] = 1
    }
    
    if int(projectedPoint3.x) >= 0 && int(projectedPoint3.x) < width &&
        int(projectedPoint3.y) >= 0 && int(projectedPoint3.y) < height {
        (*canvas)[int(projectedPoint3.y)][int(projectedPoint3.x)] = 1
    }

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

func convertSTLSolidToVector3Array(stlSolid stl.Solid) []Vector3 {
    var vertices []Vector3

    for _, triangle := range stlSolid.Triangles {
        // Extract vertices from the triangle
        point1 := Vector3{x: float64(triangle.Vertices[0][0]), y: float64(triangle.Vertices[0][1]), z: float64(triangle.Vertices[0][2])}
        point2 := Vector3{x: float64(triangle.Vertices[1][0]), y: float64(triangle.Vertices[1][1]), z: float64(triangle.Vertices[1][2])}
        point3 := Vector3{x: float64(triangle.Vertices[2][0]), y: float64(triangle.Vertices[2][1]), z: float64(triangle.Vertices[2][2])}

        // Append the vertices to the slice
        vertices = append(vertices, point1, point2, point3)
    }

    return vertices
}

// Function to clear the canvas
func clearCanvas(canvas [][]int) {
	for i := 0; i < len(canvas); i++ {
		for j := 0; j < len(canvas[i]); j++ {
			canvas[i][j] = 0
		}
	}
	fmt.Print("\033[H\033[2J") // Move cursor to top-left and clear the screen
}

// Function to print the canvas to the terminal
func printCanvas(canvas [][]int) {
	for i := 0; i < len(canvas); i++ {
		for j := 0; j < len(canvas[i]); j++ {
			if canvas[i][j] == 0 {
				fmt.Print(" ")
			} else if canvas[i][j] == 1 {
				fmt.Print("█")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	width, height, err := getTerminalSize()
	if err != nil {
		fmt.Println("Failed to get terminal size:", err)
		return
	}

	// Initialize the canvas with the terminal dimensions
	canvas := initArray(height, width) // Note the reversal of height and width
	clearCanvas(canvas) // Clear the canvas before drawing

	// Load the STL model using hschendel/stl
	stlFilePath := "model.stl"
	stlSolid, err := stl.ReadFile(stlFilePath)
	if err != nil {
		fmt.Println("Error reading STL file:", err)
		return
	}

	// Convert the STL data to your format
	model := convertSTLSolidToVector3Array(*stlSolid) // Dereference the pointer
	
	// Define a scaling factor (adjust as needed)
    scalingFactor := 1.00

    // Scale down the model's coordinates
    for i := 0; i < len(model); i++ {
        model[i].x *= scalingFactor
        model[i].y *= scalingFactor
        model[i].z *= scalingFactor
    }

	for {
		rotateModel(&model, 0.1, 0.1, 0.1)
		clearCanvas(canvas)

		projectModel(model, &canvas)
		printCanvas(canvas)

		time.Sleep(100 * time.Millisecond)
	}
}

