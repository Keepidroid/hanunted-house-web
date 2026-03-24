package game

import "testing"

func TestTraitTrackGainLoseHeal(t *testing.T) {
	track := &TraitTrack{Track: []int{2, 3, 4, 5, 6}, Index: 2, StartIndex: 2}
	track.Gain(1)
	if track.Index != 1 {
		t.Fatalf("expected index 1 after gain, got %d", track.Index)
	}
	track.Lose(10, false)
	if track.Index != track.CriticalIndex() {
		t.Fatalf("pre-haunt lose should clamp to critical")
	}
	track.HealToStart()
	if track.Index != 2 {
		t.Fatalf("heal should reset to start")
	}
}

func TestDamageAllocationValidation(t *testing.T) {
	e := &ExplorerState{Traits: map[Trait]*TraitTrack{
		TraitMight:     {Track: []int{2, 3, 4}, Index: 1, StartIndex: 1},
		TraitSpeed:     {Track: []int{2, 3, 4}, Index: 2, StartIndex: 1},
		TraitKnowledge: {Track: []int{2, 3, 4}, Index: 1, StartIndex: 1},
		TraitSanity:    {Track: []int{2, 3, 4}, Index: 1, StartIndex: 1},
	}}

	err := e.ApplyDamage("physical", map[Trait]int{TraitSpeed: 1}, false)
	if err == nil {
		t.Fatalf("expected error when allocating on critical speed while might available")
	}

	err = e.ApplyDamage("physical", map[Trait]int{TraitMight: 1}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
