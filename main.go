package main

import (
  "errors"
  "fmt"
  "image"
  "image/color"
  "image/png"
  "math"
  "math/rand"
  "os"
)

type Grid struct {
  size int
  grid []int
}

type Square struct {
  corners []Point
  midpoint Point
}

type Point struct {
  x int
  y int
}

func (p Point) Translate(x int, y int) Point {
  return Point{p.x + x, p.y + y}
}

func (g Grid) GetPoint(index int) Point {
  return Point{index % g.size, index / g.size}
}

func (g Grid) GetIndex(p Point) (int, error) {
  index := g.size * p.y + p.x
  if index < 0 || index > len(g.grid) {
    return 0, errors.New("Out of bounds.")
  }
  return index, nil
}

func (g Grid) GetSquare(topLeft Point, depth int) Square {
  size := g.CalculateDepthSize(depth);
  return Square{
    corners: []Point {
      topLeft,
      topLeft.Translate(size-1, 0),
      topLeft.Translate(0, size-1),
      topLeft.Translate(size-1, size-1),
    },
    midpoint: topLeft.Translate(size/2, size/2),
  }
}

func (g Grid) GetDiamond(midpoint Point, depth int) Square {
  size := g.CalculateDepthSize(depth-1);
  return Square{
    corners: []Point {
      midpoint.Translate(0, -1*(size-1)),
      midpoint.Translate(size-1, 0),
      midpoint.Translate(0, size-1),
      midpoint.Translate(-1*(size-1), 0),
    },
    midpoint: midpoint,
  }
}

func (g Grid) Calculate(depth int) {
  for stage := depth; stage > 0; stage-- {
    areas := g.GetAreas(stage)
    for _, square := range areas {
      g.CalculateSquare(square)
    }
    for _, square := range areas {
      g.CalculateDiamond(square, stage)
    }
  }
}

func (g Grid) GetAreas(depth int) []Square {
  step := g.CalculateDepthSize(depth)-1
  areas := make([]Square, 0)
  for y := 0; y < g.size-1; y += step {
    for x := 0; x < g.size-1; x += step {
      areas = append(areas, g.GetSquare(Point{x, y}, depth))
    }
  }
  return areas
}

func (g Grid) CalculateDiamond(s Square, depth int) {
  diamondMidpoints := g.GetDiamond(s.midpoint, depth)
  for _, midpoint := range diamondMidpoints.corners {
    diamond := g.GetDiamond(midpoint, depth)
    g.CalculateSquare(diamond)
  }
}

func (g Grid) CalculateSquare(s Square) {
  sum := 0
  amount := 0
  for _, p := range s.corners {
    if i, err := g.GetIndex(p); err == nil {
      sum += g.grid[i]
      amount++
    }
  }
  i, _ := g.GetIndex(s.midpoint)
  g.grid[i] = sum / amount
}

func (g Grid) CalculateDepthSize(depth int) int {
  return int(math.Pow(float64(2), float64(depth))) + 1
}

func (g Grid) String() string {
  var s string
  for i := range g.grid {
    s += fmt.Sprintf("%4d ", g.grid[i])
    if (i + 1) % g.size == 0 {
      s += "\n"
    }
  }
  return s
}

func (g Grid) CreateImage() *image.RGBA {
  img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{g.size, g.size}})
  for y := 0; y < g.size; y++ {
    for x := 0; x < g.size; x++ {
      gray := color.Gray{uint8(g.grid[y * g.size + x])}
      img.Set(x, y, gray)
    }
  }
  return img
}

func main() {
  base := 10
  size := int(math.Pow(float64(2), float64(base))) + 1
  grid := make([]int, size*size)
  surface := Grid{size, grid}
  square := surface.GetSquare(Point{0, 0}, base)
  for _, p := range square.corners {
    i, _ := surface.GetIndex(p)
    surface.grid[i] = rand.Intn(256)
  }
  surface.Calculate(base)
  img := surface.CreateImage()
  file, _ := os.Create("image.png")
  png.Encode(file, img)
}
