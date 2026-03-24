package httpapi

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"hanunted-house/server/internal/game"
	"hanunted-house/server/internal/room"
	"hanunted-house/server/internal/ws"
)

type API struct {
	manager *room.Manager
	hub     *ws.Hub
	logger  *slog.Logger
	roller  *game.Roller
}

func New(manager *room.Manager, hub *ws.Hub, logger *slog.Logger) *API {
	return &API{manager: manager, hub: hub, logger: logger, roller: game.NewRoller(time.Now().UnixNano())}
}

func (a *API) Router() *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		api.POST("/rooms", a.createRoom)
		api.POST("/rooms/:roomCode/join", a.joinRoom)
		api.POST("/rooms/:roomCode/rejoin", a.rejoin)
		api.POST("/rooms/:roomCode/start", a.startGame)
		api.GET("/rooms/:roomCode", a.getRoom)
		api.POST("/rooms/:roomCode/commands/damage", a.damage)
		api.POST("/rooms/:roomCode/commands/roll", a.roll)
	}
	r.GET("/ws/:roomCode", a.ws)
	return r
}

type createRoomInput struct {
	Nickname string `json:"nickname"`
}
type joinRoomInput struct {
	Nickname string `json:"nickname"`
}
type tokenInput struct {
	PlayerToken string `json:"playerToken"`
}

type damageInput struct {
	PlayerToken string             `json:"playerToken"`
	Kind        string             `json:"kind"`
	Allocation  map[game.Trait]int `json:"allocation"`
}

type rollInput struct {
	PlayerToken string        `json:"playerToken"`
	Trait       game.Trait    `json:"trait"`
	DiceCount   int           `json:"diceCount"`
	Type        game.RollType `json:"type"`
}

func (a *API) createRoom(c *gin.Context) {
	var in createRoomInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := a.manager.CreateRoom(in.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (a *API) joinRoom(c *gin.Context) {
	var in joinRoomInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := a.manager.JoinRoom(c.Param("roomCode"), in.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (a *API) rejoin(c *gin.Context) {
	var in tokenInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player, err := a.manager.Rejoin(c.Param("roomCode"), in.PlayerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"playerId": player.PlayerID, "seat": player.Seat})
	a.hub.Broadcast(strings.ToUpper(c.Param("roomCode")))
}

func (a *API) startGame(c *gin.Context) {
	var in tokenInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.manager.StartGame(c.Param("roomCode"), in.PlayerToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code := strings.ToUpper(c.Param("roomCode"))
	a.hub.Broadcast(code)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (a *API) getRoom(c *gin.Context) {
	rm, err := a.manager.GetRoomByCode(c.Param("roomCode"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roomCode": rm.RoomCode})
}

func (a *API) damage(c *gin.Context) {
	var in damageInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.manager.ApplyDamageByToken(c.Param("roomCode"), in.PlayerToken, in.Kind, in.Allocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.hub.Broadcast(strings.ToUpper(c.Param("roomCode")))
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (a *API) roll(c *gin.Context) {
	var in rollInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	diceCount := in.DiceCount
	if in.Type == game.RollTypeTrait || in.Type == game.RollTypeAttack {
		roomData, err := a.manager.GetRoomByCode(c.Param("roomCode"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		roomDataState := roomData.State
		if roomDataState == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "game not started"})
			return
		}
		player, ok := roomData.PlayerByToken(in.PlayerToken)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid player token"})
			return
		}
		explorer := roomDataState.Explorers[player.PlayerID]
		diceCount = game.TraitDiceCount(explorer, in.Trait)
	}
	res := a.roller.Roll(in.Type, diceCount)
	c.JSON(http.StatusOK, res)
}

func (a *API) ws(c *gin.Context) {
	playerID := c.Query("playerId")
	if strings.TrimSpace(playerID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playerId required"})
		return
	}
	roomCode := strings.ToUpper(c.Param("roomCode"))
	a.hub.Handle(c.Writer, c.Request, roomCode, playerID)
}
