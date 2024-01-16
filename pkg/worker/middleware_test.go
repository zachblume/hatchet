package worker

import (
	"context"
	"errors"
	"testing"
)

type testHatchetContext struct {
	context.Context
}

func (c *testHatchetContext) SetContext(ctx context.Context) {
	c.Context = ctx
}

func (c *testHatchetContext) GetContext() context.Context {
	return c.Context
}

func (c *testHatchetContext) StepOutput(step string, target interface{}) error {
	return nil
}

func (c *testHatchetContext) TriggeredByEvent() bool {
	return false
}

func (c *testHatchetContext) WorkflowInput(target interface{}) error {
	return nil
}

func TestAddMiddleware(t *testing.T) {
	m := middlewares{}
	middlewareFunc := func(ctx HatchetContext, next func(HatchetContext) error) error {
		return nil
	}
	m.add(middlewareFunc)

	if len(m.middlewares) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(m.middlewares))
	}
}

func TestRunAllWithNoMiddleware(t *testing.T) {
	m := middlewares{}
	err := m.runAll(&testHatchetContext{context.Background()}, func(ctx HatchetContext) error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRunAllWithMiddleware(t *testing.T) {
	m := middlewares{}
	called := false
	middlewareFunc := func(ctx HatchetContext, next func(HatchetContext) error) error {
		called = true
		return next(ctx)
	}
	m.add(middlewareFunc)

	err := m.runAll(&testHatchetContext{context.Background()}, func(ctx HatchetContext) error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Errorf("Expected middleware to be called")
	}
}

func TestRunAllWithPropagatedContext(t *testing.T) {
	m := middlewares{}
	key := "key"
	value := "value"

	// Middleware that sets a value in the context
	middlewareFunc := func(ctx HatchetContext, next func(HatchetContext) error) error {
		ctx.SetContext(context.WithValue(ctx, key, value))

		return next(ctx)
	}
	m.add(middlewareFunc)

	// Next function that checks for the value in the context
	err := m.runAll(&testHatchetContext{context.Background()}, func(ctx HatchetContext) error {
		if ctx.Value(key) != value {
			t.Errorf("Expected value %v in context, got %v", value, ctx.Value(key))
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRunAllWithErrorInMiddleware(t *testing.T) {
	m := middlewares{}
	expectedErr := errors.New("middleware error")
	middlewareFunc := func(ctx HatchetContext, next func(HatchetContext) error) error {
		return expectedErr
	}
	m.add(middlewareFunc)

	err := m.runAll(&testHatchetContext{context.Background()}, func(ctx HatchetContext) error {
		return nil
	})

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestRunAllWithErrorInNext(t *testing.T) {
	m := middlewares{}
	expectedErr := errors.New("next error")
	middlewareFunc := func(ctx HatchetContext, next func(HatchetContext) error) error {
		return next(ctx)
	}
	m.add(middlewareFunc)

	err := m.runAll(&testHatchetContext{context.Background()}, func(ctx HatchetContext) error {
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}
