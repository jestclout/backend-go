package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/jestclout/jestclout-go/game"
)

var (
	ErrCreateGame = errors.New("failed to create game")
)

type Handler struct {
	GameManager *game.Manager
	Logger      zerolog.Logger
}

func New(manager *game.Manager, ll zerolog.Logger) *Handler {
	return &Handler{
		GameManager: manager,
		Logger:      ll,
	}
}

func (h *Handler) PlayerIDFromRequest(r *http.Request) uint64 {
	playerHeader := r.Header.Get("X-Player-Id")

	playerID, err := strconv.ParseUint(playerHeader, 10, 64)
	if err != nil {
		return 0
	}

	return playerID
}

func (h *Handler) Router() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(h.LogRequest)

	api.HandleFunc("/game", h.CreateGame).Methods("POST")
	api.HandleFunc("/game/{gameCode}", h.GetGameState).Methods("GET")
	api.HandleFunc("/game/{gameCode}", h.ExecCommand).Methods("POST")

	return r
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	ll := hlog.FromRequest(r).With().Logger()

	state, err := h.GameManager.NewGame()
	if err != nil {
		ll.Error().Err(err).Msg("failed to create game")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(state)
	if err != nil {
		ll.Error().Err(err).Msg("failed to marshal game instance")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		ll.Error().Err(err).Msg("failed to write get state payload")
	}
}

func (h *Handler) GetGameState(w http.ResponseWriter, r *http.Request) {
	ll := hlog.FromRequest(r).With().Logger()

	vars := mux.Vars(r)
	code := vars["gameCode"]

	ll = ll.With().Str("code", code).Logger()

	playerID := h.PlayerIDFromRequest(r)

	state, err := h.GameManager.GetPublicGameState(code, playerID)
	if err != nil {
		ll.Error().Err(err).Msg("failed to retrieve game")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(state)
	if err != nil {
		ll.Error().Err(err).Msg("failed to marshal game instance")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		ll.Error().Err(err).Msg("failed to write get state payload")
	}
}

func (h *Handler) ExecCommand(w http.ResponseWriter, r *http.Request) {
	ll := hlog.FromRequest(r).With().Logger()

	vars := mux.Vars(r)
	code := vars["gameCode"]

	ll = ll.With().Str("code", code).Logger()

	var cmd game.Command
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		ll.Error().Msg("failed to unmarshal game command")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	state, err := h.GameManager.ExecCommand(code, cmd)
	if err != nil {
		ll.Error().Msg("failed to execute game command")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(state)
	if err != nil {
		ll.Error().Err(err).Msg("failed to marshal game instance")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		ll.Error().Err(err).Msg("failed to write get state payload")
	}
}
