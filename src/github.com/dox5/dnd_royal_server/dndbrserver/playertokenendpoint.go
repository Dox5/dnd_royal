package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/dox5/dnd_royal_server/api"
    "github.com/dox5/dnd_royal_server/model"
)

func getTokens(rooms *RoomManager,
               logger *log.Logger,
               request *http.Request) (interface{}, error) {

    roomId, err := api.RoomIdFromRequest(request)

    if err != nil {
        return nil, err
    }

    var tokens []model.Token

    err = rooms.WithSharedRoom(roomId, func(room *model.Room) error {
        tokens = room.GetPlayerTokens()
        return nil
    })

    return tokens, err
}

func setPosition(rooms *RoomManager,
                 logger *log.Logger,
                 request *http.Request) (interface{}, error) {

    var tokenPosition struct {
        TokenId model.Identifier `json:",string"`
        Position model.Vector
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }
    err := api.ParseJsonRequest(request, &tokenPosition)

    if err != nil {
        return nil, err
    }

    err = rooms.WithExclusiveRoom(tokenPosition.RoomId,
                                  func(room *model.Room) error {
        if tokenPosition.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised access")
        }

        token, foundIt := room.GetPlayerToken(tokenPosition.TokenId)

        if(foundIt) {
            logger.Printf("Setting token position for token %+v to %+v for room %+v",
                          tokenPosition.TokenId,
                          tokenPosition.Position,
                          tokenPosition.RoomId)

            token.Position.X = tokenPosition.Position.X
            token.Position.Y = tokenPosition.Position.Y

            return nil
        }

        return fmt.Errorf("No token found with ID %+v", tokenPosition.TokenId)
    })

    return nil, err
}

func MakePlayerTokenEndpoint(rooms *RoomManager,
                             logger *log.Logger) *Endpoint {
    endpoint := NewEndpoint()

    endpoint.Register("/getTokens",
                      http.MethodGet,
                      func(request *http.Request) (interface{}, error) {
                          return getTokens(rooms, logger, request)
                      })

    endpoint.Register("/setTokenPosition",
                      http.MethodPost,
                      func(request *http.Request) (interface{}, error) {
                          return setPosition(rooms, logger, request)
                      })

    return endpoint
}
