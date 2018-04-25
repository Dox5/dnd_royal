package model

type Room struct {
    id Identifier
    fog Fog
    gameMaster *player
    mapAsset string
    playerTokens []Token
}

func NewRoom(gameMaster *player) *Room {
    return &Room{id: MakeId(),
                 gameMaster: gameMaster,
                 mapAsset: "CoolJenniMap",
                 playerTokens: make([]Token, 0, 3)}
}

func (r *Room) GameMaster() * player {
    return r.gameMaster
}

func (r *Room) Id() Identifier {
    return r.id
}

func (r *Room) Fog() * Fog {
    return &r.fog
}

func (r *Room) MapAsset() string {
    return r.mapAsset
}

func (r *Room) Update(timeDelta float32) {
    r.fog.Advance(timeDelta)
}

func (r *Room) AddPlayerToken(position Vector) {
    token := Token {
        Id: Identifier(len(r.playerTokens)),
        Position: position }

    r.playerTokens = append(r.playerTokens, token)
}

func (r *Room) GetPlayerTokens() []Token {
    return r.playerTokens
}

func (r *Room) GetPlayerToken(id Identifier) (*Token, bool) {
    for i := 0; i < len(r.playerTokens); i += 1 {
        token := &r.playerTokens[i]
        if token.Id == id {
            return token, true
        }
    }
    return nil, false
}
