package golangcontext

import (
	"context"
	"fmt"
	"testing"
)

// Context
// context is represent as data
// this data can add value, cancel signal, timeout, and deadline signal
// context have concept called parent & child
// child context will have the same feature or functionallity like the parent context
// context is IMMUTABLE
// if we add value/signal cancel or more -> in the background it creates new context with that value/functionallity
// context have func -> Deadline(), Done(), Err(), and Value() (can see and click inside the context.Background()

func TestContext(t *testing.T) {

	// create new context (commonly use)
	background := context.Background()
	fmt.Println("Background:", background)

	// create new context like background (use if we dont know yet what type of context we wanna create, rarely use)
	todo := context.TODO()
	fmt.Println("TODO:", todo)
	
}

