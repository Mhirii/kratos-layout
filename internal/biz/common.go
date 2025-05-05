package biz

type Pagination struct {
	Page int32 `json:"page"`
	Size int32 `json:"size"`
}

type Difficulty int

const (
	EASY Difficulty = iota
	MEDIUM
	HARD
	EXPERT
)

func (d Difficulty) String(difficulty int) string {
	switch difficulty {
	case 0:
		return "Easy"
	case 1:
		return "Medium"
	case 2:
		return "Hard"
	case 3:
		return "Expert"
	default:
		return "Unknown"
	}
}

func DifficultyFromString(difficulty string) Difficulty {
	switch difficulty {
	case "Easy":
		return EASY
	case "Medium":
		return MEDIUM
	case "Hard":
		return HARD
	case "Expert":
		return EXPERT
	default:
		return EASY
	}
}

type Audit struct {
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	DeletedBy string `json:"deleted_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
