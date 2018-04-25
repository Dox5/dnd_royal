package api

import (
  "github.com/dox5/dnd_royal_server/model"
)

type RoomCreateResponse struct {
    RoomId model.Identifier `json:",string"`
    GameMasterId model.Identifier `json:",string"`
}

type RoomStateResponse struct {
    Fog model.Circle
}
