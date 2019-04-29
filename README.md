# FSM
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/dadi/fsm)
[![Coverage](http://gocover.io/_badge/github.com/dadi/fsm)](http://gocover.io/github.com/dadi/fsm)

A simple easy to use go Finite State Machine library. FSM is a programmatic map of whitelisted transitions, where each transition triggers a synchronous operation within the context of a State.

## Install

```
go get github.com/dadi/fsm
```

## Usage

```go
// Create new instance of FSM with a context
ctx := context.Background()
f := fsm.New().WithContext(ctx)

// Wildcard states are useful for transitioning to an Error or Shutdown state.
f.NewState().From("*").To("ERROR").OnEnter(func(*fsm.State) {
	// Do something
})

// Parallel transitions are none blocking.
f.NewState().From("INITIALIZING").To("READY").OnEnter(func(*fsm.State) {
	// Do something here
}).Parallel(true)

// Each state has a context that is closed before the state changes. You can use this with methods called within the state OnEnter method.
f.NewState().From("FETCHING_DATA").To("STARTING_SERVER").OnEnter(func(s *fsm.State) {
	doSomething(st.Context())
})

// Handle async transitions.
go func() {
  for transition := range f.Transitions() {
		transition.Do()
	}
}()
```