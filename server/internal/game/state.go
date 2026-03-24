package game

import "time"

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhasePreHaunt Phase = "preHaunt"
	PhaseFinished Phase = "finished"
)

type PlayerView struct {
	RoomCode string      `json:"roomCode"`
	Phase    Phase       `json:"phase"`
	Players  []SeatState `json:"players"`
	TurnSeat int         `json:"turnSeat"`
	Version  int64       `json:"version"`
}

type SeatState struct {
	Seat      int       `json:"seat"`
	PlayerID  string    `json:"playerId"`
	Nickname  string    `json:"nickname"`
	Connected bool      `json:"connected"`
	JoinedAt  time.Time `json:"joinedAt"`
	IsHost    bool      `json:"isHost"`
	IsCurrent bool      `json:"isCurrent"`
}
