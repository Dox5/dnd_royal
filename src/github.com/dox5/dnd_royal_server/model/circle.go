package model

import (
    "math"
)

func float32Equal(a float32, b float32) bool {
    const EPSILON float64 = 0.00001

    diff := math.Abs(float64(a - b))

    if a == b {
        return true
    } else if a == 0 || b == 0 || diff <= math.SmallestNonzeroFloat64 {
        return diff < (EPSILON * math.SmallestNonzeroFloat64)
    } else {
        absA := math.Abs(float64(a))
        absB := math.Abs(float64(b))
        return diff / math.Min(absA + absB,
                               math.MaxFloat64) < EPSILON
    }
}

type Vector struct {
    X float32
    Y float32
}

func (v1 Vector) Add(v2 Vector) Vector {
    return Vector{v1.X + v2.X,
                  v1.Y + v2.Y}
}

func (v1 Vector) Sub(v2 Vector) Vector {
    return Vector{v1.X - v2.X,
                  v1.Y - v2.Y}
}

func (v Vector) MultiplyScalar(amount float32) Vector {
    return Vector{X: v.X * amount,
                  Y: v.Y * amount}
}

func (v Vector) DivideScalar(amount float32) Vector {
    return Vector{X: v.X / amount,
                  Y: v.Y / amount}
}

func (v Vector) AddScalar(amount float32) Vector {
    return Vector{X: v.X + amount,
                  Y: v.Y + amount}
}

func (v Vector) SubScalar(amount float32) Vector {
    return Vector{X: v.X - amount,
                  Y: v.Y - amount}
}

func (v Vector) Magnatude() float32 {
    squareXY := math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2)
    return float32(math.Sqrt(squareXY))
}

func (v1 Vector) Equal(v2 Vector) bool {
    return float32Equal(v1.X, v2.X) && float32Equal(v1.Y, v2.Y)
}

type Circle struct {
    Centre Vector
    Radius float32
}

func (c1 Circle) DistanceTo(c2 Circle) Vector {
    return c2.Centre.Sub(c1.Centre)
}

func (c1 Circle) Equal(c2 Circle) bool {
    return c1.Centre.Equal(c2.Centre) && float32Equal(c1.Radius, c2.Radius)
}
