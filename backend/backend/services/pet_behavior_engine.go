package services

// PetBehaviorEngine 行为引擎
type PetBehaviorEngine struct{}

func NewPetBehaviorEngine() *PetBehaviorEngine {
	return &PetBehaviorEngine{}
}

// GenerateActionSequence 生成动作序列
func (e *PetBehaviorEngine) GenerateActionSequence(mood string) []string {
	sequences := map[string][]string{
		"happy":   {"wave", "dance", "sing"},
		"sad":     {"comfort", "hug", "cheer"},
		"excited": {"jump", "spin", "celebrate"},
	}
	if seq, ok := sequences[mood]; ok {
		return seq
	}
	return []string{"idle"}
}
