package middleware

import (
	"errors"
	"fmt"
)

type StateSetter interface {
	StateSet(*State)
}

type StateRef struct{ *State }

func (cr *StateRef) StateSet(c *State) { cr.State = c }

// State web application store
type State struct {
	S map[string]interface{} `vugu:"data"`
}

// NewState returns a new State store
func NewState() *State {
	return &State{
		S: make(map[string]interface{}),
	}
}

// Set a value in the store
func (st *State) Set(key string, value interface{}) error {
	fmt.Printf("[STATE] Set: %s\n", key)

	st.S[key] = value
	return nil
}

// Get a value from the store
func (st *State) Get(key string) (interface{}, error) {
	fmt.Printf("[STATE] Get: %s\n", key)

	if val, ok := st.S[key]; ok {
		return val, nil
	}
	return nil, errors.New("key not found")
}

// Unset removes a value from the store
func (st *State) Unset(key string) {
	delete(st.S, key)
}
