package room

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"hanunted-house/server/internal/content"
	"hanunted-house/server/internal/game"
)

const (
	MinSeats = 3
	MaxSeats = 6
)

type Status string

const (
	StatusLobby   Status = "lobby"
	StatusRunning Status = "running"
)

type Player struct {
	PlayerID    string
	PlayerToken string
	Nickname    string
	Seat        int
	Connected   bool
	JoinedAt    time.Time
}

type Room struct {
	RoomID    string
	RoomCode  string
	Status    Status
	HostSeat  int
	TurnSeat  int
	Players   map[string]*Player
	seatIndex map[int]string
	Version   int64
	State     *game.RuntimeState
	mu        sync.RWMutex
}

type Manager struct {
	mu      sync.RWMutex
	byCode  map[string]*Room
	randSrc *rand.Rand
}

func NewManager() *Manager {
	return &Manager{byCode: make(map[string]*Room), randSrc: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

type CreateRoomResult struct {
	RoomCode    string `json:"roomCode"`
	PlayerID    string `json:"playerId"`
	PlayerToken string `json:"playerToken"`
	Seat        int    `json:"seat"`
}

type JoinRoomResult = CreateRoomResult

func (m *Manager) CreateRoom(nickname string) (CreateRoomResult, error) {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return CreateRoomResult{}, errors.New("nickname required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	code := m.generateCodeLocked()
	room := &Room{RoomID: uuid.NewString(), RoomCode: code, Status: StatusLobby, HostSeat: 0, TurnSeat: 0, Players: map[string]*Player{}, seatIndex: map[int]string{}, Version: 1}
	player := &Player{PlayerID: uuid.NewString(), PlayerToken: uuid.NewString(), Nickname: nickname, Seat: 0, Connected: true, JoinedAt: time.Now()}
	room.Players[player.PlayerID] = player
	room.seatIndex[0] = player.PlayerID
	m.byCode[code] = room

	return CreateRoomResult{RoomCode: code, PlayerID: player.PlayerID, PlayerToken: player.PlayerToken, Seat: 0}, nil
}

func (m *Manager) JoinRoom(code, nickname string) (JoinRoomResult, error) {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return JoinRoomResult{}, errors.New("nickname required")
	}

	room, err := m.GetRoomByCode(code)
	if err != nil {
		return JoinRoomResult{}, err
	}

	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Status != StatusLobby {
		return JoinRoomResult{}, errors.New("room already started")
	}

	seat, ok := room.firstEmptySeat()
	if !ok {
		return JoinRoomResult{}, errors.New("room is full")
	}

	player := &Player{PlayerID: uuid.NewString(), PlayerToken: uuid.NewString(), Nickname: nickname, Seat: seat, Connected: true, JoinedAt: time.Now()}
	room.Players[player.PlayerID] = player
	room.seatIndex[seat] = player.PlayerID
	room.Version++

	return JoinRoomResult{RoomCode: room.RoomCode, PlayerID: player.PlayerID, PlayerToken: player.PlayerToken, Seat: seat}, nil
}

func (m *Manager) Rejoin(code, token string) (*Player, error) {
	room, err := m.GetRoomByCode(code)
	if err != nil {
		return nil, err
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	for _, p := range room.Players {
		if p.PlayerToken == token {
			p.Connected = true
			room.Version++
			return p, nil
		}
	}
	return nil, errors.New("invalid player token")
}

func (m *Manager) StartGame(code, hostToken string) error {
	room, err := m.GetRoomByCode(code)
	if err != nil {
		return err
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Status != StatusLobby {
		return errors.New("game already started")
	}
	hostID, ok := room.seatIndex[room.HostSeat]
	if !ok || room.Players[hostID].PlayerToken != hostToken {
		return errors.New("only host can start")
	}
	if len(room.Players) < MinSeats {
		return errors.New("need at least 3 players")
	}

	state := &game.RuntimeState{RoomCode: room.RoomCode, Phase: game.PhasePreHaunt, TurnSeat: room.HostSeat, TurnNumber: 1, StateVersion: 1, Explorers: map[string]*game.ExplorerState{}}
	for _, player := range room.Players {
		char, err := content.CharacterByIndex(player.Seat % content.CharacterPoolSize())
		if err != nil {
			return err
		}
		explorerID := uuid.NewString()
		state.Explorers[player.PlayerID] = &game.ExplorerState{
			ExplorerID: explorerID,
			PlayerID:   player.PlayerID,
			Seat:       player.Seat,
			Name:       player.Nickname,
			Character:  char.Name,
			Side:       game.SideHeroes,
			Traits: map[game.Trait]*game.TraitTrack{
				game.TraitMight:     {Track: char.MightTrack, Index: char.MightStart, StartIndex: char.MightStart},
				game.TraitSpeed:     {Track: char.SpeedTrack, Index: char.SpeedStart, StartIndex: char.SpeedStart},
				game.TraitKnowledge: {Track: char.KnowledgeTrack, Index: char.KnowledgeStart, StartIndex: char.KnowledgeStart},
				game.TraitSanity:    {Track: char.SanityTrack, Index: char.SanityStart, StartIndex: char.SanityStart},
			},
		}
	}

	room.Status = StatusRunning
	room.State = state
	room.TurnSeat = room.HostSeat
	room.Version++
	return nil
}

func (m *Manager) Disconnect(code, playerID string) {
	room, err := m.GetRoomByCode(code)
	if err != nil {
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if p, ok := room.Players[playerID]; ok {
		p.Connected = false
		room.Version++
	}
}

func (m *Manager) GetRoomByCode(code string) (*Room, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, ok := m.byCode[code]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (r *Room) PlayerView(playerID string) game.PlayerView {
	r.mu.RLock()
	defer r.mu.RUnlock()
	players := make([]game.SeatState, 0, len(r.Players))
	for seat := 0; seat < MaxSeats; seat++ {
		id, ok := r.seatIndex[seat]
		if !ok {
			continue
		}
		p := r.Players[id]
		players = append(players, game.SeatState{Seat: p.Seat, PlayerID: p.PlayerID, Nickname: p.Nickname, Connected: p.Connected, JoinedAt: p.JoinedAt, IsHost: p.Seat == r.HostSeat, IsCurrent: p.PlayerID == playerID})
	}
	phase := game.PhaseLobby
	if r.Status == StatusRunning {
		phase = game.PhasePreHaunt
	}
	view := game.PlayerView{RoomCode: r.RoomCode, Phase: phase, Players: players, TurnSeat: r.TurnSeat, Version: r.Version}
	if r.State != nil {
		if e := r.State.Explorers[playerID]; e != nil {
			traits := make([]game.TraitValue, 0, 4)
			for _, t := range []game.Trait{game.TraitMight, game.TraitSpeed, game.TraitKnowledge, game.TraitSanity} {
				if tr := e.Traits[t]; tr != nil {
					traits = append(traits, game.TraitValue{Trait: t, TrackIndex: tr.Index, Value: tr.CurrentValue(), Critical: tr.IsCritical()})
				}
			}
			view.Explorer = &game.ExplorerView{ExplorerID: e.ExplorerID, Character: e.Character, Side: e.Side, Dead: e.Dead, Traits: traits}
		}
	}
	return view
}

func (r *Room) firstEmptySeat() (int, bool) {
	for i := 0; i < MaxSeats; i++ {
		if _, ok := r.seatIndex[i]; !ok {
			return i, true
		}
	}
	return 0, false
}

func (m *Manager) generateCodeLocked() string {
	for {
		const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
		buf := make([]byte, 6)
		for i := range buf {
			buf[i] = letters[m.randSrc.Intn(len(letters))]
		}
		code := string(buf)
		if _, exists := m.byCode[code]; !exists {
			return code
		}
	}
}

func (r *Room) PlayerByToken(token string) (*Player, bool) {
	for _, p := range r.Players {
		if p.PlayerToken == token {
			return p, true
		}
	}
	return nil, false
}

func (m *Manager) ApplyDamageByToken(code, token, kind string, allocation map[game.Trait]int) error {
	room, err := m.GetRoomByCode(code)
	if err != nil {
		return err
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.State == nil {
		return errors.New("game not started")
	}
	p, ok := room.PlayerByToken(token)
	if !ok {
		return errors.New("invalid player token")
	}
	explorer := room.State.Explorers[p.PlayerID]
	if explorer == nil {
		return errors.New("explorer not found")
	}
	if err := explorer.ApplyDamage(kind, allocation, room.State.Phase == game.PhasePostHaunt); err != nil {
		return err
	}
	room.State.StateVersion++
	room.Version++
	return nil
}
