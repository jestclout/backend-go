package game

type Config struct {
	CodeLength      int
	MinPlayers      int
	MaxPlayers      int
	Rounds          int
	PromptsPerRound int
}

func DefaultConfig() Config {
	return Config{
		CodeLength:      4,
		MinPlayers:      3,
		MaxPlayers:      8,
		Rounds:          2,
		PromptsPerRound: 2,
	}
}
