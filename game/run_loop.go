package game

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCmdMissingPlayer = errors.New("command missing player")
)

type CommandType int

const (
	GetState CommandType = iota
	AddPlayer
	UpdatePlayer
	RemovePlayer
	StartGame
	AnswerPrompt
)

type Command struct {
	Type     CommandType `json:"cmdType"`
	PlayerID uint64      `json:"playerId"`
	Player   *Player     `json:"player"`
	PromptID uint64      `json:"promptId"`
	Answer   *Answer     `json:"answer"`
}

type RunLoop struct {
	Instance *Instance
	mu       sync.Mutex
}

func NewGame(code string, prompts []string, config Config) *RunLoop {
	return &RunLoop{
		Instance: NewInstance(code, prompts, config),
	}
}

func (g *RunLoop) AdvanceState(currentState State, currentVotingPrompt int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	instance := g.Instance

	if instance.currentState == currentState {
		switch currentState {
		case VotingOnAnswers:
			if instance.currentVotingPrompt == currentVotingPrompt {
				instance.AdvanceState()
			}
		default:
			instance.AdvanceState()
		}
	}
}

func (g *RunLoop) ExecCommand(cmd Command) (*PublicGameState, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	instance := g.Instance

	switch cmd.Type {
	case AddPlayer:
		if cmd.Player == nil {
			return nil, ErrCmdMissingPlayer
		}

		err := instance.AddPlayer(cmd.Player)
		if err != nil {
			return nil, err
		}

	case UpdatePlayer:
		if cmd.Player == nil {
			return nil, ErrCmdMissingPlayer
		}

		err := instance.UpdatePlayer(cmd.Player)
		if err != nil {
			return nil, err
		}

	case RemovePlayer:
		if cmd.Player == nil {
			return nil, ErrCmdMissingPlayer
		}

		err := instance.RemovePlayer(cmd.Player.ID)
		if err != nil {
			return nil, err
		}

	case StartGame:
		err := instance.AdvanceState()
		if err != nil {
			return nil, err
		}

		// Set timeout for answering prompts.
		go func() {
			time.Sleep(60 * time.Second)
			g.AdvanceState(AnsweringPrompts, 0)
		}()

	case AnswerPrompt:
	}

	return instance.GetState(cmd.PlayerID), nil
}

func (g *RunLoop) GetPublicState(playerID uint64) *PublicGameState {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Instance.GetState(playerID)
}
