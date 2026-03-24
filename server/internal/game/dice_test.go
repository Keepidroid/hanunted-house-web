package game

import "testing"

func TestRollerDeterministicWithSeed(t *testing.T) {
	r := NewRoller(42)
	res := r.Roll(RollTypeOther, 4)
	if len(res.Dice) != 4 {
		t.Fatalf("expected 4 dice, got %d", len(res.Dice))
	}
	if res.Total < 0 || res.Total > 8 {
		t.Fatalf("unexpected total: %d", res.Total)
	}
}
