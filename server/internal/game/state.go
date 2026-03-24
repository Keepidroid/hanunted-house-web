package game

import "time"

type Phase string

type Trait string

type Side string

const (
	PhaseLobby     Phase = "lobby"
	PhasePreHaunt  Phase = "preHaunt"
	PhasePostHaunt Phase = "postHaunt"
	PhaseFinished  Phase = "finished"
)

const (
	TraitMight     Trait = "might"
	TraitSpeed     Trait = "speed"
	TraitKnowledge Trait = "knowledge"
	TraitSanity    Trait = "sanity"
)

const (
	SideHeroes  Side = "heroes"
	SideUnknown Side = "unknown"
)

type TraitTrack struct {
	Track      []int `json:"track"`
	Index      int   `json:"index"`
	StartIndex int   `json:"startIndex"`
}

type ExplorerState struct {
	ExplorerID string                `json:"explorerId"`
	PlayerID   string                `json:"playerId"`
	Seat       int                   `json:"seat"`
	Name       string                `json:"name"`
	Character  string                `json:"character"`
	Side       Side                  `json:"side"`
	Dead       bool                  `json:"dead"`
	Traits     map[Trait]*TraitTrack `json:"traits"`
}

type RuntimeState struct {
	RoomCode     string                    `json:"roomCode"`
	Phase        Phase                     `json:"phase"`
	TurnSeat     int                       `json:"turnSeat"`
	TurnNumber   int                       `json:"turnNumber"`
	StateVersion int64                     `json:"stateVersion"`
	Explorers    map[string]*ExplorerState `json:"explorers"`
}

type PlayerView struct {
	RoomCode string        `json:"roomCode"`
	Phase    Phase         `json:"phase"`
	Players  []SeatState   `json:"players"`
	TurnSeat int           `json:"turnSeat"`
	Version  int64         `json:"version"`
	Explorer *ExplorerView `json:"explorer,omitempty"`
}

type ExplorerView struct {
	ExplorerID string       `json:"explorerId"`
	Character  string       `json:"character"`
	Side       Side         `json:"side"`
	Dead       bool         `json:"dead"`
	Traits     []TraitValue `json:"traits"`
}

type TraitValue struct {
	Trait      Trait `json:"trait"`
	TrackIndex int   `json:"trackIndex"`
	Value      int   `json:"value"`
	Critical   bool  `json:"critical"`
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
