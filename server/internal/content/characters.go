package content

import "fmt"

type CharacterTemplate struct {
	ID             string
	Name           string
	MightTrack     []int
	SpeedTrack     []int
	KnowledgeTrack []int
	SanityTrack    []int
	MightStart     int
	SpeedStart     int
	KnowledgeStart int
	SanityStart    int
}

var defaultCharacters = []CharacterTemplate{
	{
		ID:             "ox-bellows",
		Name:           "Ox Bellows",
		MightTrack:     []int{2, 2, 3, 4, 5, 6, 6, 7},
		SpeedTrack:     []int{2, 2, 2, 3, 4, 5, 5, 6},
		KnowledgeTrack: []int{2, 2, 3, 3, 4, 5, 5, 6},
		SanityTrack:    []int{2, 2, 3, 3, 4, 5, 5, 6},
		MightStart:     3,
		SpeedStart:     5,
		KnowledgeStart: 3,
		SanityStart:    3,
	},
	{
		ID:             "vivian-lopez",
		Name:           "Vivian Lopez",
		MightTrack:     []int{2, 2, 3, 3, 4, 5, 5, 6},
		SpeedTrack:     []int{2, 2, 3, 4, 4, 5, 6, 7},
		KnowledgeTrack: []int{2, 3, 4, 4, 5, 5, 6, 7},
		SanityTrack:    []int{2, 2, 3, 4, 5, 5, 6, 6},
		MightStart:     3,
		SpeedStart:     3,
		KnowledgeStart: 2,
		SanityStart:    3,
	},
	{
		ID:             "father-rhinehardt",
		Name:           "Father Rhinehardt",
		MightTrack:     []int{2, 2, 3, 3, 4, 4, 5, 6},
		SpeedTrack:     []int{2, 3, 3, 4, 4, 5, 5, 6},
		KnowledgeTrack: []int{3, 4, 4, 5, 5, 6, 7, 7},
		SanityTrack:    []int{3, 4, 4, 5, 5, 6, 6, 7},
		MightStart:     4,
		SpeedStart:     4,
		KnowledgeStart: 0,
		SanityStart:    2,
	},
}

func CharacterByIndex(index int) (CharacterTemplate, error) {
	if index < 0 || index >= len(defaultCharacters) {
		return CharacterTemplate{}, fmt.Errorf("character index out of range")
	}
	return defaultCharacters[index], nil
}

func CharacterPoolSize() int { return len(defaultCharacters) }
