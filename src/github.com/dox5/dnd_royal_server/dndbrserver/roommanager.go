package main

import (
  "fmt"
  "sync"
  "time"

  "github.com/dox5/dnd_royal_server/model"
)

const (
    UpdateRateHz float32 = 2
)

type activeRoom struct {
    room *model.Room
    roomLock sync.RWMutex
    shutdown chan bool
}

type RoomManager struct {
    rooms map[model.Identifier]*activeRoom
    managerLock sync.RWMutex
}

func roomAdvancer(room *activeRoom, updateRate float32) {
    accumulator := float64(0)
    periodSeconds := float64(float32(1) / updateRate)

    for {
        updateStart := time.Now()

        select {
        case <- room.shutdown:
            return
        default:
        }

        room.roomLock.Lock()
        for ; accumulator > periodSeconds ; accumulator -= periodSeconds {
            room.room.Update(float32(periodSeconds))
        }
        room.roomLock.Unlock()

        sleepFor := (periodSeconds / 10.0) * float64(time.Second)
        time.Sleep(time.Duration(sleepFor))

        updateEnd := time.Now()

        updateDuration := updateEnd.Sub(updateStart)
        accumulator += updateDuration.Seconds()
    }
}

func NewRoomManager() *RoomManager {
    rm := &RoomManager{}
    rm.rooms = make(map[model.Identifier]*activeRoom)
    return rm
}

func (rm *RoomManager) Create() *model.Room {
    gm := model.NewPlayer()
    r := model.NewRoom(gm)

    rm.managerLock.Lock()
    defer rm.managerLock.Unlock()
    active := &activeRoom{room: r, shutdown: make(chan bool)}

    // TODO: this should be controlable!
    numPlayers := 3
    for i := 0; i < numPlayers; i += 1 {
        pos := model.Vector{X: -50,
                            Y: float32(50 * i)}
        active.room.AddPlayerToken(pos)
    }

    rm.rooms[r.Id()] = active

    go roomAdvancer(active, UpdateRateHz)

    return r
}

func (rm *RoomManager) Count() int {
    rm.managerLock.RLock()
    defer rm.managerLock.RUnlock()
    return len(rm.rooms)
}

func (rm *RoomManager) getActiveRoom(roomId model.Identifier) (*activeRoom, error) {
    rm.managerLock.RLock()
    defer rm.managerLock.RUnlock()

    if activeRoom, ok := rm.rooms[roomId]; ok {
        return activeRoom, nil
    } else {
        return nil, fmt.Errorf("No room found with id %+v", roomId)
    }
}

func (rm *RoomManager) GetCurrentFog(roomId model.Identifier) (model.Circle, error) {
    activeRoom, err := rm.getActiveRoom(roomId)

    if err != nil {
        return model.Circle{}, err
    }

    activeRoom.roomLock.RLock()
    defer activeRoom.roomLock.RUnlock()

    return activeRoom.room.Fog().Current(), nil
}

type RoomUpdateCallback func (* model.Room) error
func (rm *RoomManager) WithExclusiveRoom(roomId model.Identifier,
                                         callback RoomUpdateCallback) error {
    room, err := rm.getActiveRoom(roomId)

    if err != nil {
        return err
    }

    room.roomLock.Lock()
    defer room.roomLock.Unlock()
    err = callback(room.room)

    return err
}

func (rm *RoomManager) WithSharedRoom(roomId model.Identifier,
                                      callback RoomUpdateCallback) error {
    room, err := rm.getActiveRoom(roomId)

    if err != nil {
        return err
    }

    room.roomLock.RLock()
    defer room.roomLock.RUnlock()
    err = callback(room.room)

    return err
}
