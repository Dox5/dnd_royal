package model_test

import (
    "testing"

    "github.com/dox5/dnd_royal_server/model"
)

func TestNewFogShouldHaveTargetAndCurrentSet(t *testing.T) {
    expectedTarget := model.Circle{Centre: model.Vector{0, 0}, Radius: 10}
    expectedCurrent := model.Circle{Centre: model.Vector{0, 0}, Radius: 25}

    actual := model.NewFog(expectedCurrent)
    actual.SetTarget(expectedTarget)

    if actual.Target() != expectedTarget {
        t.Errorf("Expected target circle to be %+v but it was %+v",
                 expectedTarget,
                 actual.Target())
    }

    if actual.Current() != expectedCurrent {
        t.Errorf("Expected initial circle to be %+v but it was %+v",
                 expectedCurrent,
                 actual.Target())
    }
}

func TestAdvanceFullPeriodShouldMatchTarget(t *testing.T) {
    target := model.Circle{Centre: model.Vector{10, 0}, Radius: 10}
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}

    period := float32(30)

    fog := model.NewFog(initial)
    fog.Resume()
    fog.SetTarget(target)
    fog.SetPeriod(period)

    fog.Advance(period)

    if fog.Current() != target {
        t.Errorf("After full time expected initial circle %+v to be the target %+v",
                 fog.Current(),
                 target)
    }
}

func TestPartialAdvanceShouldMoveCorrectAmount(t *testing.T) {
    // Moving 1/3 of the distance
    totalPeriod := float32(60)
    partialPeriod := float32(20)

    initial := model.Circle{Centre: model.Vector{0, 21}, Radius: 12}

    fog := model.NewFog(initial)
    fog.Resume()
    fog.SetPeriod(totalPeriod)

    fog.Advance(partialPeriod)

    expected := model.Circle{Centre: model.Vector{0, 14}, Radius: 8}

    if expected != fog.Current() {
        t.Errorf("After 1/3 period expected to have circle %+v but was actually %+v",
               expected,
               fog.Current())
    }
}

func TestFogWithNoPeriodShouldNotChangeOnAdvance(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{5, 17}, Radius: 13}
    fog := model.NewFog(initial)
    fog.Resume()

    if fog.Current() != initial {
        t.Errorf("Expected initial fog to be set correctly. Expected %+v but it was %+v",
                 initial,
                 fog.Current())
    }

    fog.Advance(150)

    if fog.Current() != initial {
        t.Errorf("Expected fog to be unchanged after advance. Expected %+v but it was %+v",
                 initial,
                 fog.Current())
    }
}

func TestChangePeriodAfterAdvanceShouldProduceCorrectMove(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{80, 80}, Radius: 40}
    fog := model.NewFog(initial)
    fog.Resume()
    fog.SetPeriod(10)

    // Move half way
    fog.Advance(5)

    halfWay := model.Circle{Centre: model.Vector{40, 40}, Radius: 20}
    oneStep := model.Circle{Centre: halfWay.Centre.SubScalar(10),
                            Radius: halfWay.Radius - 5}

    // Would've been another Advance(5) to the target, but now needs Advance(4)
    fog.SetPeriod(4)
    fog.Advance(1)

    if fog.Current() != oneStep {
        t.Errorf("Expected fog to be %+v but it was %+v",
                 oneStep,
                 fog.Current())
    }
}

func TestNewlyCreatedFogShouldBePaused(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    fog := model.NewFog(initial)

    if !fog.Paused() {
        t.Error("Expected fog to be paused be it wasn't!")
    }

    fog.SetPeriod(2)
    fog.Advance(1)

    if fog.Current() != initial {
        t.Error("Expected fog to have not advanced while paused, but it did")
    }
}

func TestResumeShouldEnableAdvancement(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    fog := model.NewFog(initial)

    fog.SetPeriod(2)
    fog.Resume()
    fog.Advance(1)

    if fog.Current() == initial {
        t.Error("Expected fog to  advance after resume but it was unchanged")
    }
}

func TestPauseAfterResumeShouldStopAdvancement(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    target := model.Circle{Centre: model.Vector{0, 0}, Radius: 10}

    fog := model.NewFog(initial)
    fog.SetTarget(target)

    fog.SetPeriod(2)
    fog.Resume()
    fog.Advance(1)

    fog.Pause()
    fog.Advance(1)

    if fog.Current() == target {
        t.Error("Expected to have not reached target yet")
    }
}

func TestSubSecondAdvanceShouldWork(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    target := model.Circle{Centre: model.Vector{0, 0}, Radius: 10}

    fog := model.NewFog(initial)
    fog.SetTarget(target)

    fog.SetPeriod(1)
    fog.Resume()

    for i := 0 ; i < 4 ; i += 1 {
        fog.Advance(0.25)
    }

    if fog.Current() != target {
        t.Errorf("Expected fog to reach target (%+v) but it was actually %+v",
                 target,
                 fog.Current())
    }
}

func TestAdvancePastTargetShouldStop(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    target := model.Circle{Centre: model.Vector{0, 0}, Radius: 40}

    fog := model.NewFog(initial)
    fog.SetTarget(target)

    fog.SetPeriod(5)
    fog.Resume()

    fog.Advance(10)

    if fog.Current() != target {
        t.Errorf("Expcted to stop at target but did not!\n" +
                 "Target: %+v\n" +
                 "Current: %+v\n",
                 target,
                 fog.Current())
    }
}

func TestAdvanceInSmallStepsGetsCloseToTarget(t *testing.T) {
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}
    target := model.Circle{Centre: model.Vector{0, 0}, Radius: 40}

    fog := model.NewFog(initial)
    fog.SetTarget(target)

    period := 1
    steps  := 1001 // Rounding error means it's going to take a few more steps
    fog.SetPeriod(float32(period))
    fog.Resume()

    stepSize := float64(period) / float64(steps)
    for i := 0; i < steps ; i += 1 {
        fog.Advance(float32(stepSize))
    }

    if !fog.Current().Equal(target) {
        t.Errorf("Expcted to stop at target but did not!\n" +
                 "Target: %+v\n" +
                 "Current: %+v\n",
                 target,
                 fog.Current())
    }
}

func TestArriveAtTargetShouldPause(t *testing.T) {
    target := model.Circle{Centre: model.Vector{10, 0}, Radius: 10}
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}

    period := float32(30)

    fog := model.NewFog(initial)
    fog.SetTarget(target)
    fog.SetPeriod(period)
    fog.Resume()

    fog.Advance(period)

    if !fog.Paused() {
        t.Errorf("After arriving at target, should go back to paused state")
    }
}

func TestSetTargetShouldAdvanceRate(t *testing.T) {
    target := model.Circle{Centre: model.Vector{10, 0}, Radius: 10}
    target2 := model.Circle{Centre: model.Vector{-10, 0}, Radius: 50}
    initial := model.Circle{Centre: model.Vector{0, 0}, Radius: 50}

    fog := model.NewFog(initial)
    fog.SetTarget(target)
    fog.SetPeriod(30)

    initialRate := fog.Rate()
    fog.SetTarget(target2)

    if initialRate == fog.Rate() {
        t.Errorf("Expected setting target to change rate but it stayed the same")
    }


}
