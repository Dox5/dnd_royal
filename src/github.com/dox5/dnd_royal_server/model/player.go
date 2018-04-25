package model

type player struct {
    id Identifier
}

func NewPlayer() *player {
    return &player{id: MakeId()}
}

func (p *player) Id() Identifier {
    return p.id
}

