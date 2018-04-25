package model

import (
    "math"
)

type Fog struct {
    target      Circle
    current     Circle
    period      float32
    advance     bool
    advanceRate Rate
}

func NewFog(initial Circle) * Fog {
    fog := &Fog{target: Circle{},
                current: initial,
                period: 1,
                advance: false,
                advanceRate: Rate{}}

    return fog
}

func (f *Fog) recalculateRate() {
    radiusDelta := f.target.Radius - f.current.Radius
    radiusRate := radiusDelta / f.period
    f.advanceRate.Radius = radiusRate

    translation := f.current.DistanceTo(f.target)
    f.advanceRate.Translation = translation.DivideScalar(f.period)
}

func (f *Fog) SetTarget(target Circle) {
    f.target = target
    f.recalculateRate()
}

func (f *Fog) SetPeriod(period float32) {
    f.period = period
    f.recalculateRate()
}

func (f *Fog) Advance(timeDelta float32) {
    if !f.advance {
        return
    }

    deltaRadius := timeDelta * f.advanceRate.Radius
    deltaTranslation := f.advanceRate.Translation.MultiplyScalar(timeDelta)

    fullStep := Circle{Centre: f.current.Centre.Add(deltaTranslation),
                       Radius: f.current.Radius + deltaRadius}

    // Deal with the circle
    fullStepDistance := f.current.DistanceTo(fullStep).Magnatude()
    targetDistance := f.current.DistanceTo(f.target).Magnatude()

    arrived := true
    if fullStepDistance < targetDistance {
        // Safe to do the full move
        f.current.Centre = fullStep.Centre
        // Not arrived
        arrived = false
    } else {
        // Would land on or go past target
        f.current.Centre = f.target.Centre
    }

    // Deal with the radius
    deltaStepRadius   := math.Abs(float64(f.current.Radius - fullStep.Radius))
    deltaTargetRadius := math.Abs(float64(f.current.Radius - f.target.Radius))

    if deltaStepRadius < deltaTargetRadius {
        // Safe to do the full step
        f.current.Radius = fullStep.Radius
        // Not arrived
        arrived = false
    } else {
        // Would be >= target radius
        f.current.Radius = f.target.Radius
    }

    if arrived {
        f.Pause()
    }
}

func (f *Fog) Target() Circle {
    return f.target
}

func (f *Fog) Current() Circle {
    return f.current
}

func (f *Fog) Rate() Rate {
    return f.advanceRate
}

func (f *Fog) Pause() {
    f.advance = false
}

func (f *Fog) Resume() {
    f.advance = true
}

func (f *Fog) Paused() bool {
    return !f.advance
}
