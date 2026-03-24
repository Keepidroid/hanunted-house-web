package game

import (
	"math/rand"
)

type RollType string

const (
	RollTypeTrait  RollType = "trait"
	RollTypeAttack RollType = "attack"
	RollTypeHaunt  RollType = "haunt"
	RollTypeOther  RollType = "other"
)

type DiceRollResult struct {
	Type     RollType `json:"type"`
	Dice     []int    `json:"dice"`
	Total    int      `json:"total"`
	DiceUsed int      `json:"diceUsed"`
}

type Roller struct {
	rand *rand.Rand
}

func NewRoller(seed int64) *Roller {
	return &Roller{rand: rand.New(rand.NewSource(seed))}
}

func (r *Roller) Roll(rt RollType, diceCount int) DiceRollResult {
	if diceCount < 0 {
		diceCount = 0
	}
	res := DiceRollResult{Type: rt, Dice: make([]int, 0, diceCount), DiceUsed: diceCount}
	for i := 0; i < diceCount; i++ {
		n := r.rand.Intn(3) // 0/1/2 对应山屋骰子的0/1/2成功符号近似
		res.Dice = append(res.Dice, n)
		res.Total += n
	}
	return res
}

func TraitDiceCount(e *ExplorerState, trait Trait) int {
	if e == nil || e.Traits[trait] == nil {
		return 0
	}
	return e.Traits[trait].CurrentValue()
}
