package game

import (
	"errors"
)

var (
	ErrPlayerLimitReached  = errors.New("player limit reached")
	ErrPlayerMinimumNotmet = errors.New("players below minimum")
	ErrGameInProgress      = errors.New("game in progress")
	ErrPlayerNotFound      = errors.New("player not found")
)

type Player struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type State int

const (
	WaitingForPlayers State = iota
	AnsweringPrompts
	VotingOnAnswers
	ScoringRound
	ScoringGame
	Done
)

type Prompt struct {
	ID   uint64 `json:"id"`
	Text string `json:"text"`
}

type Answer struct {
	ID   uint64 `json:"id"`
	Text string `json:"text"`
}

type PublicGameState struct {
	Code         string            `json:"code"`
	Players      []*Player         `json:"players"`
	CurrentState State             `json:"currentState"`
	Prompt       *Prompt           `json:"prompt"`
	Answers      []*Answer         `json:"answers"`
	Votes        map[uint64]uint64 `json:"votes"`
}

type UserAnswer struct {
	PlayerID uint64
	Prompt   *Prompt `json:"prompt"`
	Answer   *Answer `json:"answers"`
}

type Round struct {
	Prompts     []*Prompt
	UserAnswers map[uint64][]*UserAnswer
}

type Instance struct {
	code                string
	players             []*Player
	currentState        State
	currentRound        int
	currentVotingPrompt int
	rounds              []*Round
	possiblePrompts     []string
	nextPlayerID        uint64
	nextPromptID        uint64
	nextAnswerID        uint64
	config              Config
}

func NewInstance(code string, prompts []string, config Config) *Instance {
	return &Instance{
		code:                code,
		players:             make([]*Player, 0, config.MaxPlayers),
		currentState:        WaitingForPlayers,
		currentRound:        0,
		currentVotingPrompt: 0,
		rounds:              []*Round{},
		possiblePrompts:     prompts,
		nextPlayerID:        1,
		nextPromptID:        1,
		nextAnswerID:        1,
		config:              config,
	}
}

func (i *Instance) GetState(playerID uint64) *PublicGameState {
	gameState := &PublicGameState{
		Code:         i.code,
		Players:      i.players,
		CurrentState: i.currentState,
	}

	if playerID != 0 {
		switch i.currentState {
		case AnsweringPrompts:
			round := i.rounds[i.currentRound]
			userAnswers := round.UserAnswers[playerID]
			for _, userAnswer := range userAnswers {
				if userAnswer.Answer == nil {
					gameState.Prompt = userAnswer.Prompt

					break
				}
			}

			// case VotingOnAnswers:
			// 	round := i.Rounds[i.CurrentRound]
			// 	userAnswers := round.UserAnswers[playerID]
		}

	}

	return gameState
}

func (i *Instance) getNextUserID() uint64 {
	id := i.nextPlayerID
	i.nextPlayerID++

	return id
}

func (i *Instance) getNextPromptID() uint64 {
	id := i.nextPromptID
	i.nextPromptID++

	return id
}

func (i *Instance) getNextAnswerID() uint64 {
	id := i.nextAnswerID
	i.nextAnswerID++

	return id
}

func (i *Instance) AddPlayer(p *Player) error {
	if len(i.players) >= i.config.MaxPlayers {
		return ErrPlayerLimitReached
	}

	if i.currentState != WaitingForPlayers {
		return ErrGameInProgress
	}

	p.ID = i.getNextUserID()

	i.players = append(i.players, p)

	return nil
}

func (i *Instance) UpdatePlayer(p *Player) error {
	if i.currentState != WaitingForPlayers {
		return ErrGameInProgress
	}

	for _, player := range i.players {
		if player.ID == p.ID {
			player.Name = p.Name
			return nil
		}
	}

	return ErrPlayerNotFound
}

func (i *Instance) RemovePlayer(id uint64) error {
	if i.currentState != WaitingForPlayers {
		return ErrGameInProgress
	}

	for j, player := range i.players {
		if player.ID == id {
			i.players = append(i.players[:j], i.players[j+1:]...)
			return nil
		}
	}

	return ErrPlayerNotFound
}

func (i *Instance) getRandomPrompt() string {
	return i.possiblePrompts[0]
}

func (i *Instance) createRound() *Round {
	var prompts []*Prompt
	for j := 0; j < len(i.players); j++ {
		pt := &Prompt{
			ID:   i.getNextPromptID(),
			Text: i.getRandomPrompt(),
		}

		prompts = append(prompts, pt)
	}

	userAnswers := make(map[uint64][]*UserAnswer)
	for j, p := range i.players {

		answers := []*UserAnswer{
			{
				PlayerID: p.ID,
				Prompt:   prompts[j%len(prompts)],
			},
			{
				PlayerID: p.ID,
				Prompt:   prompts[(j-1)%len(prompts)],
			},
		}

		userAnswers[p.ID] = answers
	}

	return &Round{
		Prompts:     prompts,
		UserAnswers: userAnswers,
	}
}

func (i *Instance) startNewRound() {
	i.currentState = AnsweringPrompts

	round := i.createRound()
	i.rounds = append(i.rounds, round)
	i.currentRound = i.currentRound + 1
}

func (i *Instance) AdvanceState() error {
	switch i.currentState {
	case WaitingForPlayers:
		if len(i.players) <= i.config.MinPlayers {
			return ErrPlayerMinimumNotmet
		}

		i.startNewRound()
		return nil

	case AnsweringPrompts:
		return nil

	default:
		return nil
	}

}

func (i *Instance) AddAnswer(playerID uint64, promptID uint64, answer string) {

}

func (i *Instance) Vote(playerID uint64, answerID uint64) {

}

func (i *Instance) AdvanceVoting() {

}

func (i *Instance) ScoreRound() {

}
