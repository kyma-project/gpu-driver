package tmp

import (
	"context"
	"fmt"
)

type Action func(ctx context.Context) (context.Context, error)

type BaseState struct {
	Name string
}

type StateOne struct {
	*BaseState
	OneProp string
}

type StateTwo struct {
	*StateOne
	TwoProp string
}

func SetStateToCtx(ctx context.Context, state any) context.Context {
	return context.WithValue(ctx, "state", state)
}

func GetStateFromCtx[T any](ctx context.Context) T {
	x := ctx.Value("state")
	return x.(T)
}

func ActionOneOne(ctx context.Context) (context.Context, error) {
	fmt.Println("ActionOneOne")
	state := GetStateFromCtx[*StateOne](ctx)
	fmt.Printf("  Name: %s\n", state.Name)
	state.OneProp = "one"
	return ctx, nil
}

func ActionOneTwo(ctx context.Context) (context.Context, error) {
	fmt.Println("ActionOneTwo")
	state := GetStateFromCtx[*StateOne](ctx)
	fmt.Printf("  Name: %s\n", state.Name)
	fmt.Printf("  OneProp: %s\n", state.OneProp)
	return ctx, nil
}

func ActionTwoOne(ctx context.Context) (context.Context, error) {
	fmt.Println("ActionTwoOne")
	state := &StateTwo{
		StateOne: GetStateFromCtx[*StateOne](ctx),
		TwoProp:  "two",
	}
	ctx = SetStateToCtx(ctx, state)
	return ctx, nil
}

func ActionTwoTwo(ctx context.Context) (context.Context, error) {
	fmt.Println("ActionTwoTwo")
	state := GetStateFromCtx[*StateTwo](ctx)
	fmt.Printf("  Name: %s\n", state.Name)
	fmt.Printf("  OneProp: %s\n", state.OneProp)
	fmt.Printf("  TwoProp: %s\n", state.TwoProp)
	return ctx, nil
}
