package game

import (
	"errors"
	"sync"

	"github.com/jestclout/jestclout-go/rand"
)

var (
	ErrGameNotFound = errors.New("game not found")
)

type Manager struct {
	games   map[string]*RunLoop
	config  Config
	prompts []string
	mu      sync.Mutex
}

func NewManager(prompts []string, config Config) *Manager {
	return &Manager{
		games:   make(map[string]*RunLoop),
		prompts: prompts,
		config:  config,
	}
}

func (m *Manager) NewGame() (*PublicGameState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for {
		code, err := rand.NewGameCode(m.config.CodeLength)
		if err != nil {
			return nil, err
		}

		if _, ok := m.games[code]; !ok {
			game := NewGame(code, m.prompts, m.config)

			m.games[code] = game

			return game.GetPublicState(0), nil
		}
	}
}

func (m *Manager) getGame(code string) (*RunLoop, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, ok := m.games[code]
	if !ok {
		return nil, ErrGameNotFound
	}

	return game, nil
}

func (m *Manager) GetPublicGameState(code string, playerID uint64) (*PublicGameState, error) {
	game, err := m.getGame(code)
	if err != nil {
		return nil, err
	}

	state := game.GetPublicState(playerID)

	return state, nil
}

func (m *Manager) ExecCommand(code string, cmd Command) (*PublicGameState, error) {
	game, err := m.getGame(code)
	if err != nil {
		return nil, err
	}

	state, err := game.ExecCommand(cmd)
	if err != nil {
		return nil, err
	}

	return state, nil
}
