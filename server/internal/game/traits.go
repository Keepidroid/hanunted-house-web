package game

import "fmt"

func (t *TraitTrack) CurrentValue() int {
	if len(t.Track) == 0 || t.Index < 0 || t.Index >= len(t.Track) {
		return 0
	}
	return t.Track[t.Index]
}

func (t *TraitTrack) CriticalIndex() int {
	if len(t.Track) == 0 {
		return 0
	}
	return len(t.Track) - 1
}

func (t *TraitTrack) SkullIndex() int {
	return len(t.Track)
}

func (t *TraitTrack) IsCritical() bool {
	return t.Index >= t.CriticalIndex()
}

func (t *TraitTrack) Gain(steps int) {
	if steps <= 0 {
		return
	}
	t.Index -= steps
	if t.Index < 0 {
		t.Index = 0
	}
}

func (t *TraitTrack) Lose(steps int, postHaunt bool) {
	if steps <= 0 {
		return
	}
	t.Index += steps
	if !postHaunt && t.Index > t.CriticalIndex() {
		t.Index = t.CriticalIndex()
	}
}

func (t *TraitTrack) HealToStart() {
	if t.Index > t.StartIndex {
		t.Index = t.StartIndex
	}
}

func (e *ExplorerState) ApplyDamage(kind string, allocation map[Trait]int, postHaunt bool) error {
	allowed := map[Trait]bool{}
	switch kind {
	case "physical":
		allowed[TraitMight], allowed[TraitSpeed] = true, true
	case "mental":
		allowed[TraitKnowledge], allowed[TraitSanity] = true, true
	case "general":
		allowed[TraitMight], allowed[TraitSpeed], allowed[TraitKnowledge], allowed[TraitSanity] = true, true, true, true
	default:
		return fmt.Errorf("invalid damage kind")
	}

	for trait, steps := range allocation {
		if !allowed[trait] {
			return fmt.Errorf("trait %s not allowed for %s damage", trait, kind)
		}
		if steps < 0 {
			return fmt.Errorf("damage steps cannot be negative")
		}
		track := e.Traits[trait]
		if track == nil {
			return fmt.Errorf("missing trait track %s", trait)
		}
		if !postHaunt && track.IsCritical() {
			for t, tr := range e.Traits {
				if allowed[t] && !tr.IsCritical() {
					return fmt.Errorf("cannot allocate to critical while another legal trait is above critical")
				}
			}
		}
	}

	for trait, steps := range allocation {
		e.Traits[trait].Lose(steps, postHaunt)
	}

	if postHaunt {
		for _, trait := range e.Traits {
			if trait.Index >= trait.SkullIndex() {
				e.Dead = true
				break
			}
		}
	}

	return nil
}

func (e *ExplorerState) Heal(trait Trait) {
	if t := e.Traits[trait]; t != nil {
		t.HealToStart()
	}
}
