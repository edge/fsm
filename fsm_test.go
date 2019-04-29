package fsm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	st := sm.NewState().From("bar")
	assert.Equal(st.Source, []string{"bar"}, "should set state source list")
}

func TestTo(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	st := sm.NewState().To("bar")
	assert.Equal(st.Destination, "bar", "should set state to value")
}

func TestName(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	sm.NewState().From("bar").To("foo")
	sm.NewState().From("foo").To("bar")

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.Transition("bar")
	name := sm.Name()

	assert.EqualValues(name, "bar", "should return name when state is defined")
}

func TestNoName(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	name := sm.Name()

	assert.EqualValues(name, "", "should be empty when no states are defined")
}

func TestContextReturnNil(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	st1 := sm.NewState().From("bar").To("foo")
	sm.NewState().From("foo").To("bar")

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.Transition("bar")
	sm.Transition("foo")

	assert.Nil(st1.Context(), "state should have no context")
}

func TestContextReturnNotNil(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	sm := New().WithContext(ctx)
	st1 := sm.NewState().From("bar").To("foo")
	sm.NewState().From("foo").To("bar")

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.Transition("bar")
	sm.Transition("foo")

	assert.NotNil(st1.Context(), "state should have context")
}

func TestParallelState(t *testing.T) {
	assert := assert.New(t)
	p := true
	sm := New()
	done := make(chan bool)
	f := func(*State) {
		done <- true
	}

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.NewState().From("bar").To("foo")
	st := sm.NewState().From("foo").To("bar").OnEnter(f).Parallel(p)

	assert.Equal(st.parallel, p, "should set parallel to input bool")
	sm.Transition("bar")
	called := <-done
	assert.True(called, "should call on enter function")
}

func TestSameStateTransition(t *testing.T) {
	assert := assert.New(t)
	sm := New()

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.NewState().From("bar").To("foo")
	sm.Transition("foo")
	err := sm.Transition("foo")

	assert.Nil(err, "should return <nil> when new state is the same as current")
}

func TestMissingTransition(t *testing.T) {
	assert := assert.New(t)

	sm := New()
	err := sm.Transition("foo")

	assert.EqualError(err, "Invalid state: foo", "should return an error when state does not exist")
}

func TestCancelContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	sm := New().WithContext(ctx)

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	st := sm.NewState().From("bar").To("foo")
	sm.NewState().From("foo").To("bar")

	sm.Transition("foo")
	sm.Transition("bar")

	assert.Error(st.ctx.Err(), "should cancel context")
	assert.EqualError(st.ctx.Err(), "context canceled", "should return canceled error")
}

func TestContextCreation(t *testing.T) {
	assert := assert.New(t)
	sm := New()

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	st1 := sm.NewState().From("foo").To("bar")
	st2 := sm.NewState().From("bar").To("foo")

	sm.Transition("foo")

	assert.Nil(st2.ctx, "state context should be <nil>")
	assert.Nil(st2.cancel, "state cancel func should be <nil>")

	// Apply a context to state machine.
	ctx := context.Background()
	sm = sm.WithContext(ctx)

	sm.Transition("bar")

	assert.NotNil(st1.ctx, "state context should have a context")
	assert.NotNil(st1.cancel, "state should have a cancel func")
}

func TestInvalidTransition(t *testing.T) {
	assert := assert.New(t)
	sm := New()

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.NewState().From("foo").To("bar")
	sm.NewState().From("bar").To("foo")
	sm.NewState().From("foo").To("baz")
	sm.NewState().From("baz").To("foo")

	sm.Transition("foo")
	sm.Transition("bar")
	err := sm.Transition("baz")

	assert.EqualError(err, "Invalid state change: bar > baz", "should return an error when state change is not permitted")
}

func TestValidTransition(t *testing.T) {
	assert := assert.New(t)
	sm := New()

	// Execute transitions.
	go func() {
		for transition := range sm.Transitions() {
			transition.Do()
		}
	}()

	sm.NewState().From("foo").To("bar")
	sm.NewState().From("bar").To("foo")

	sm.Transition("foo")
	err := sm.Transition("bar")

	assert.Nil(err, "should not return an error")
}

func TestNewState(t *testing.T) {
	assert := assert.New(t)
	sm := New()
	assert.Len(sm.States, 0, "should start with no states")

	st := sm.NewState()
	assert.Len(sm.States, 1, "should insert new state")
	assert.NotNil(st, "should return new state")
}

func TestWithContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	sm := New().WithContext(ctx)

	// Initial test.
	assert.NotNil(sm.ctx, "should have context")
	assert.NotNil(sm.cancel, "should have cancel func")
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	sm := New()

	// Initial test.
	assert.NotNil(sm.transitions, "new state machine should have transitions channel")
}
