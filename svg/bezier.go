package svg

import "math"

const quadLengthApproximationInterval = 0.01
const cubicLengthApproximationInterval = 0.005

// A QuadraticBezier represents a 2nd degree Bezier curve
type QuadraticBezier struct {
	Start   Point
	Control Point
	End     Point
}

// Bounds computes the bounding box for the Bezier curve.
func (q *QuadraticBezier) Bounds() Rect {
	minX, maxX := quadraticBezierExtrema(q.Start.X, q.Control.X, q.End.X)
	minY, maxY := quadraticBezierExtrema(q.Start.Y, q.Control.Y, q.End.Y)
	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

// Length approximates the length of the curve.
func (q *QuadraticBezier) Length() float64 {
	var length float64
	for t := float64(0); t < 1; t += quadLengthApproximationInterval {
		segment := Line{q.Evaluate(t), q.Evaluate(t + quadLengthApproximationInterval)}
		length += segment.Length()
	}
	return length
}

// Evaluate gets a point on the bezier curve for a parameter between 0 and 1.
func (q *QuadraticBezier) Evaluate(t float64) Point {
	x := quadraticBezierPolynomial(q.Start.X, q.Control.X, q.End.X, t)
	y := quadraticBezierPolynomial(q.Start.Y, q.Control.Y, q.End.Y, t)
	return Point{x, y}
}

// From returns the curve's start point.
func (q *QuadraticBezier) From() Point {
	return q.Start
}

// To returns the curve's end point.
func (q *QuadraticBezier) To() Point {
	return q.End
}

func quadraticBezierExtrema(A, B, C float64) (min, max float64) {
	min = math.Min(A, C)
	max = math.Max(A, C)
	if t := (B - A) / (2*B - A - C); t >= 0 && t <= 1 {
		extreme := quadraticBezierPolynomial(A, B, C, t)
		min = math.Min(min, extreme)
		max = math.Max(max, extreme)
	}
	return
}

func quadraticBezierPolynomial(A, B, C, t float64) float64 {
	return math.Pow(1-t, 2)*A + 2*(1-t)*t*B + t*t*C
}

// A CubicBezier represents a 3rd degree Bezier curve.
type CubicBezier struct {
	Start    Point
	Control1 Point
	Control2 Point
	End      Point
}

// Bounds computes the bounding box for the Bezier curve.
func (c *CubicBezier) Bounds() Rect {
	minX := math.Min(c.Start.X, c.End.X)
	maxX := math.Max(c.Start.X, c.End.X)
	minY := math.Min(c.Start.Y, c.End.Y)
	maxY := math.Max(c.Start.Y, c.End.Y)

	xExtrema := cubicBezierExtrema(c.Start.X, c.Control1.X, c.Control2.X, c.End.X)
	yExtrema := cubicBezierExtrema(c.Start.Y, c.Control1.Y, c.Control2.Y, c.End.Y)
	for _, xValue := range xExtrema {
		minX = math.Min(minX, xValue)
		maxX = math.Max(maxX, xValue)
	}
	for _, yValue := range yExtrema {
		minY = math.Min(minY, yValue)
		maxY = math.Max(maxY, yValue)
	}

	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

// Length approximates the length of the curve.
func (c *CubicBezier) Length() float64 {
	var length float64
	for t := float64(0); t < 1; t += cubicLengthApproximationInterval {
		segment := Line{c.Evaluate(t), c.Evaluate(t + cubicLengthApproximationInterval)}
		length += segment.Length()
	}
	return length
}

// Evaluate gets a point on the bezier curve for a parameter between 0 and 1.
func (c *CubicBezier) Evaluate(t float64) Point {
	x := cubicBezierPolynomial(c.Start.X, c.Control1.X, c.Control2.X, c.End.X, t)
	y := cubicBezierPolynomial(c.Start.Y, c.Control1.Y, c.Control2.Y, c.End.Y, t)
	return Point{x, y}
}

// From returns the curve's start point.
func (c *CubicBezier) From() Point {
	return c.Start
}

// To returns the curve's end point.
func (c *CubicBezier) To() Point {
	return c.End
}

func cubicBezierExtrema(A, B, C, D float64) []float64 {
	// These coefficients result from taking the derivative of the cubic bezier
	// polynomial.
	a := 3*D - 9*C + 9*B - 3*A
	b := 6*A - 12*B + 6*C
	c := 3 * (B - A)
	discriminant := math.Pow(b, 2) - 4*a*c
	if discriminant < 0 {
		return []float64{}
	}

	solution1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	solution2 := (-b - math.Sqrt(discriminant)) / (2 * a)
	result := make([]float64, 0, 2)
	if solution1 >= 0 && solution1 <= 1 {
		result = append(result, cubicBezierPolynomial(A, B, C, D, solution1))
	}
	if solution2 >= 0 && solution2 <= 1 {
		result = append(result, cubicBezierPolynomial(A, B, C, D, solution2))
	}
	return result
}

func cubicBezierPolynomial(A, B, C, D, t float64) float64 {
	return A*math.Pow(1-t, 3) + 3*B*t*math.Pow(1-t, 2) + 3*C*(1-t)*t*t + D*t*t*t
}
