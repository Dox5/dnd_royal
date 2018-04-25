package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/dox5/dnd_royal_server/api"
    "github.com/dox5/dnd_royal_server/model"
)

type fogEndpoint struct {
    rooms *RoomManager
    logger *log.Logger
}

func NewFogEndpoint(logger * log.Logger, rooms *RoomManager) fogEndpoint {
    return fogEndpoint{logger: logger, rooms: rooms}
}

func (endpoint fogEndpoint) ServeHTTP(writer http.ResponseWriter,
                                      request *http.Request) {

    if endpoint.rooms == nil {
        http.Error(writer,
                   "rooms pointer should not be nil",
                   http.StatusInternalServerError)
        return
    }

    path := request.URL.Path
    var response interface{}
    var err error

    switch path {
    case "/resume":
        response, err = endpoint.resume(request)
    case "/pause":
        response, err = endpoint.pause(request)
    case "/paused":
        response, err = endpoint.paused(request)
    case "/setTarget":
        response, err = endpoint.setTarget(request)
    case "/setPeriod":
        response, err = endpoint.period(request)
    case "/location":
        response, err = endpoint.location(request)
    case "/getTarget":
        response, err = endpoint.getTarget(request)
    case "/advanceTime":
        response, err = endpoint.advanceTime(request)
    default:
        err = fmt.Errorf("Unknown endpoint for fog: %s", path)
    }

    api.FormatResponse(writer, response, err)
}

func (endpoint fogEndpoint) resume(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodPost {
        return nil, fmt.Errorf("resume is a POST endpoint")
    }

    var resumeRequest struct {
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }

    err := api.ParseJsonRequest(request, &resumeRequest)

    if err != nil {
        return nil, err
    }

    err = endpoint.rooms.WithExclusiveRoom(resumeRequest.RoomId,
                                           func(room *model.Room) error {
        if resumeRequest.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised")
        }

        endpoint.logger.Printf("Resumed room %+v", room.Id())
        room.Fog().Resume()
        return nil
    })

    return nil, err
}

func (endpoint fogEndpoint) pause(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodPost {
        return nil, fmt.Errorf("pause is a POST endpoint")
    }

    var pauseRequest struct {
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }

    err := api.ParseJsonRequest(request, &pauseRequest)

    if err != nil {
        return nil, err
    }

    err = endpoint.rooms.WithExclusiveRoom(pauseRequest.RoomId,
                                           func(room *model.Room) error {
        if pauseRequest.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised")
        }

        endpoint.logger.Printf("Paused room %+v", room.Id())
        room.Fog().Pause()
        return nil
    })

    return nil, err
}

func (endpoint fogEndpoint) paused(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodGet {
        return nil, fmt.Errorf("paused is a GET endpoint")
    }

    roomId, err := api.RoomIdFromRequest(request)

    if err != nil {
        return nil, err
    }

    response := struct {IsPaused bool}{}

    err = endpoint.rooms.WithSharedRoom(roomId, func(room *model.Room) error {
        response.IsPaused = room.Fog().Paused()
        return nil
    })

    return response, err
}

func (endpoint fogEndpoint) period(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodPost {
        return nil, fmt.Errorf("period is a POST endpoint")
    }

    var periodRequest struct {
        Period float32
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }
    err := api.ParseJsonRequest(request, &periodRequest)

    if err != nil {
        return nil, err
    }

    err = endpoint.rooms.WithExclusiveRoom(periodRequest.RoomId,
                                           func(room *model.Room) error {
        if periodRequest.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised access")
        }

        endpoint.logger.Printf("Setting period to %v for room %v",
                               periodRequest.Period,
                               periodRequest.RoomId)
        room.Fog().SetPeriod(periodRequest.Period)

        return nil
    })

    return nil, err
}

func (endpoint fogEndpoint) location(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodGet {
        return nil, fmt.Errorf("location is a GET endpoint")
    }

    roomId, err := api.RoomIdFromRequest(request)

    if err != nil {
        return nil, err
    }

    var fogState struct {
        Current model.Circle
        Target  model.Circle
        Rate   model.Rate
    }

    err = endpoint.rooms.WithSharedRoom(roomId, func(room *model.Room) error {
        fog := room.Fog()
        fogState.Current = fog.Current()
        fogState.Target  = fog.Target()

        if !fog.Paused() {
            fogState.Rate = fog.Rate()
        }
        return nil
    })

    return fogState, err
}

func (endpoint fogEndpoint) setTarget(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodPost {
        return nil, fmt.Errorf("target is a POST endpoint")
    }

    var targetRequest struct {
        FogTarget model.Circle
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }

    err := api.ParseJsonRequest(request, &targetRequest)

    if err != nil {
        return nil, err
    }

    err = endpoint.rooms.WithExclusiveRoom(targetRequest.RoomId,
                                           func(room *model.Room) error {
        if targetRequest.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised")
        }

        endpoint.logger.Printf("Setting target %+v", targetRequest)
        room.Fog().SetTarget(targetRequest.FogTarget)
        return nil
    })

    return nil, err
}

func (endpoint fogEndpoint) getTarget(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodGet {
        return nil, fmt.Errorf("getTarget is a GET endpoint")
    }

    roomId, err := api.RoomIdFromRequest(request)

    if err != nil {
        return nil, err
    }

    var targetLocation struct {
        Target model.Circle
    }

    err = endpoint.rooms.WithSharedRoom(roomId, func(room *model.Room) error {
        targetLocation.Target = room.Fog().Target()
        return nil
    })

    return targetLocation, err
}

func (endpoint fogEndpoint) advanceTime(request *http.Request) (interface{}, error) {
    if request.Method != http.MethodPost {
        return nil, fmt.Errorf("advanceTime is a POST endpoint")
    }

    var advanceRequest struct {
        Amount float32
        RoomId model.Identifier `json:",string"`
        GameMasterId model.Identifier `json:",string"`
    }
    err := api.ParseJsonRequest(request, &advanceRequest)

    if err != nil {
        return nil, err
    }

    err = endpoint.rooms.WithExclusiveRoom(advanceRequest.RoomId,
                                           func(room *model.Room) error {
        if advanceRequest.GameMasterId != room.GameMaster().Id() {
            return fmt.Errorf("Unautherised access")
        }

        endpoint.logger.Printf("Advancing game time for room %v by %v",
                               advanceRequest.RoomId,
                               advanceRequest.Amount)
        room.Fog().Advance(advanceRequest.Amount)

        return nil
    })

    return nil, err
}
