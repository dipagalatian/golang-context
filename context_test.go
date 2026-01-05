package golangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
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

// Context with Cancel
// context can have cancel signal, we can cancel the context with func -> cancel()
// we can check if the context is canceled or not with func -> Err()
// we can get the cancel signal with func -> Done()

func TestContextWithCancelAndErr(t *testing.T) {

	// create new context
	background := context.Background()

	// add cancel signal to context
	context, cancel := context.WithCancel(background)

	// cancel the context
	cancel()

	// check if the context is canceled or not
	fmt.Println("Context:", context.Err())
	fmt.Println("Context Done:", context.Done())

}

func CreateCounter(ctx context.Context) chan int {

	destination := make(chan int)

	go func ()  {
		defer close(destination)

		counter := 1
		for {
			select {
			case <- ctx.Done():
				return
			default:
				destination <- counter
				counter++

				time.Sleep(2 * time.Second) // simulate long process
			}
			
		}
		
	}()

	return destination	
}

func TestContextWithCancel(t *testing.T) {

	fmt.Println("Total goroutines start:", runtime.NumGoroutine())

	parentCtx := context.Background()
	// add cancel signal to context
	ctx, cancel := context.WithCancel(parentCtx)

	destination := CreateCounter(ctx)

	fmt.Println("Total goroutines when CreateCounter running:", runtime.NumGoroutine())
	
	for n := range destination {
		fmt.Println("Counter:", n)
		if n == 10 {
			break
		}

	}
	cancel()

	time.Sleep(2 * time.Second) // to make sure goroutine have time to close
	
	fmt.Println("Total goroutines end:", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {

	// This context with timeout is good for handling long process that we want to limit the time to run
	// for example, we call the external API that sometimes take too long to respond
	// or we call database query that sometimes take too long to respond

	fmt.Println("Total goroutines start:", runtime.NumGoroutine())
	
	parentCtx := context.Background()

	// add timeout signal to context
	// after 5 seconds, the context will be canceled automatically eventhough the CreateCounter functin is not finished yet
	ctx, cancel := context.WithTimeout(parentCtx, 5 * time.Second)
	defer cancel() // good practice to call cancel in defer to avoid memory leak eventhough the context already have timeout

	destination := CreateCounter(ctx)
	fmt.Println("Total goroutines running:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter:", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutines end:", runtime.NumGoroutine())
}


