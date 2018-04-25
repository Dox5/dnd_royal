package model_test

import (
    "testing"

    "github.com/dox5/dnd_royal_server/model"
)

type IdSet struct {
    idSet map[model.Identifier]bool
}

func (set *IdSet) Add(id model.Identifier) bool {
    _, found := set.idSet[id]
    set.idSet[id] = true
    return found
}

func TestNewRoomShouldHaveNoPlayerTokens(t *testing.T) {
    gm := model.NewPlayer()
    room := model.NewRoom(gm)

    if len(room.GetPlayerTokens()) > 0 {
        t.Errorf("Expected room to have no tokens but had %v tokens",
                 len(room.GetPlayerTokens()))
    }
}

func TestAddPlayerTokenToRoomShouldIncreaseLength(t *testing.T) {
    gm := model.NewPlayer()
    room := model.NewRoom(gm)

    pos := model.Vector{}
    room.AddPlayerToken(pos)

    if len(room.GetPlayerTokens()) != 1 {
        t.Errorf("Expected room to have 1 token but had %v tokens",
                 len(room.GetPlayerTokens()))
    }

    room.AddPlayerToken(pos)

    if len(room.GetPlayerTokens()) != 2 {
        t.Errorf("Expected room to have 2 tokens but had %v tokens",
                 len(room.GetPlayerTokens()))
    }
}

func TestAddPlayerTokenShouldHaveGivenPosition(t *testing.T) {
    gm := model.NewPlayer()
    room := model.NewRoom(gm)

    pos := model.Vector{X: 15, Y: 7}
    room.AddPlayerToken(pos)

    if len(room.GetPlayerTokens()) < 1 {
        t.Fatal("No player token added to room!")
    }

    fromRoom := room.GetPlayerTokens()[0]
    if fromRoom.Position != pos {
        t.Errorf("Expected token to have position %+v but it was %+v",
                 pos,
                 fromRoom.Position)
    }
}

func TestAddPlayerTokenShouldBeGivenUniqueIdendifier(t *testing.T) {
    gm := model.NewPlayer()
    room := model.NewRoom(gm)

    numTokens := 10

    for i := 0 ; i < numTokens; i += 1 {
        room.AddPlayerToken(model.Vector{})
    }

    usedIds := IdSet{make(map[model.Identifier]bool)}

    for _, token := range room.GetPlayerTokens() {
        alreadyUsed := usedIds.Add(token.Id)

        if alreadyUsed {
            t.Errorf("The token ID %v was already used!", token.Id)
        }
    }
}
