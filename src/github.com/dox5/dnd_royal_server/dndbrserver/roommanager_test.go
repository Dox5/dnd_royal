package main_test

import (
  "testing"

  "github.com/dox5/dnd_royal_server/dndbrserver"
)

func TestCreateRoomShouldStoreRoom(t *testing.T) {
    rooms := main.NewRoomManager()
    createdRoom := rooms.Create()

    if createdRoom == nil {
        t.Error("Created room should not be nil")
    }

    if rooms.Count() != 1 {
        t.Errorf("Expected room count to be 1 but it was %v", rooms.Count())
    }
}

func TestCreateRoomShouldAttachGameMaster(t *testing.T) {
    rooms := main.NewRoomManager()
    createdRoom := rooms.Create()

    if createdRoom == nil {
        t.Fatal("Created room should not be nil")
    }

    if createdRoom.GameMaster() == nil {
        t.Error("Created room should have game master assigned")
    }
}

func TestGetFogForRoomThatDoesNotExistShouldReturnError(t *testing.T) {
    rooms := main.NewRoomManager()

    _, err := rooms.GetCurrentFog(15)

    if err == nil {
        t.Error("No room exists, should get an error")
    }
}

func TestGetFogForRoomShouldReturnFog(t *testing.T) {
    rooms := main.NewRoomManager()

    room := rooms.Create()

    _, err := rooms.GetCurrentFog(room.Id())

    if err != nil {
        t.Errorf("Room exists, should be able to get room. Got error %s", err)
    }
}
