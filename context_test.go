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

// Context With Value
// context can have value, we can add value to context
// we can add value to context with func -> WithValue(parent, key, value)
// we can get value from context with func -> Value(context, key)

func TestContextWithValue(t *testing.T) {

	// create new context
	contextA := context.Background()

	// add value to context from parent contextA
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	// add value to context from parent contextB
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	// add value to context from parent contextC
	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println("contextA:", contextA)
	fmt.Println("contextB:", contextB)
	fmt.Println("contextC:", contextC)
	fmt.Println("contextD:", contextD)
	fmt.Println("contextE:", contextE)
	fmt.Println("contextF:", contextF)

	// access the value of context
	// if the value is not exist in the context, it will return nil
	fmt.Println("value contextB:", contextB.Value("b")) // exist
	fmt.Println("value contextC:", contextC.Value("b")) // nil, there is no key "b" in contextC or parents of contextC"
	fmt.Println("value contextB from child contextD:", contextD.Value("b")) // exist, get value key "b" from child contextD
	fmt.Println("value contextD from parent contextB:", contextB.Value("d")) // nil, cannot get value contextD from parent contextB

	

	// // create new context
	// background := context.Background()

	// // add value to context
	// contextA := context.WithValue(background, "a", "A")
	// contextB := context.WithValue(background, "b", "B")
	// contextC := context.WithValue(contextA, "c", "C")
	// contextD := context.WithValue(contextB, "d", "D")
	// contextE := context.WithValue(contextC, "e", "E")

	// // get value from context
	// fmt.Println("Context A:", contextA.Value("a"))
	// fmt.Println("Context B:", contextB.Value("b"))
	// fmt.Println("Context C:", contextC.Value("c"))
	// fmt.Println("Context D:", contextD.Value("d"))
	// fmt.Println("Context E:", contextE.Value("e"))

}


