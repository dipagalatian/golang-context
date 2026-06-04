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

func TestContextWithParentChildTimeout(t *testing.T) {
	// create new context
	rootCtx := context.Background()

	// create child ctx with timeout
	childCtx, cancel := context.WithTimeout(rootCtx, 2*time.Second)
	defer cancel() // good practice to call cancel in defer to avoid memory leak eventhough the context already have timeout

	// now we can pass the childCtx to another function that need context, and we can check if the context is canceled or not with func -> Err() or Done()
	fmt.Println("Child Context Err:", childCtx.Err())
	fmt.Println("Child Context Done:", childCtx.Done())
	fmt.Println("Child context created successfully", childCtx)
}

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

// Context with value pass to another func layers

// Main func (controller layer func)

// private type for key userId
type contextKey string
const userIdKey contextKey = "userId"

func TestCtxWithValuePassToAnotherFunc(t *testing.T) {
	// create context
	rootCtx := context.Background()

	// create value to the context
	ctxWithUser := context.WithValue(rootCtx, userIdKey, "user_dipa_1234")

	// pass the context to ProcessBooking func
	ProcessBooking(ctxWithUser, "court_1")
}

// Service layer func (business logic func)
func ProcessBooking(ctx context.Context, courtId string) {
	fmt.Printf("Process booking for court %s...\n", courtId)

	// pass the context to SaveToDatabase func
	SaveToDatabase(ctx, courtId)
}

// Repository layer func (data access func)
func SaveToDatabase(ctx context.Context, courtId string) {
	// get the userId from context and validate it exist or not
	if userId, ok := ctx.Value(userIdKey).(string); ok {
		fmt.Printf("Saving booking for user %s and court %s to database...\n", userId, courtId)
	} else {
		fmt.Println("User ID not found in context, cannot save booking to database.")
	}
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

// example context with cancel to stop the long process in anothe goroutine
// main func that create context with cancel
func TestCtxWithCancelToStopLongProcess(t *testing.T) {
	// create context with cancel
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// call the DBWorker func
	go DBWorker(childCtx)
	
	// let the DBworker func run for 2 seconds
	time.Sleep(2 * time.Second)

	// cancel the context manually to stop the DBWorker func
	fmt.Printf("MAIN: Canceling the context to stop the DBWorker...\n")
	cancel()

	// check if the context is canceled or not, this need for seeing the logs on case Done in select DBWorker
	time.Sleep(500 * time.Millisecond)
}

// example of background worker goroutine that will cancel the process if the context is canceled
func DBWorker(ctx context.Context) {
	// loop to simulate long process and wait the cancel signal from context
	for {
		select {
		case <- ctx.Done():
			fmt.Println("DBWorker: Received cancel signal. Stopping the DB operation...")
			return
		default: 
			fmt.Println("DBWorker: Running DB operation...")
			time.Sleep(500 * time.Millisecond)
		}
	}

}

// Context with timeout
// This example use case is for handling long process payment take before it canceled because of timeout
func TestContextWithTimeoutToStopLongProcess(t *testing.T) {
	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	fmt.Println("Starting payment verification with timeout...")
	err := PaymentVerification(ctxTimeout)
	if err != nil {
		fmt.Println("Payment verification failed:", err)
	} else {
		fmt.Println("Payment verification succeeded.")
	}
}

// payment verification func that will cancel the process if the context is canceled
func PaymentVerification(ctx context.Context) error {
	ch := make(chan string, 1)

	go func() {
		// change this to 3 seconds to simulate timeouts
		// change this to 1 second to simulate success
		time.Sleep(5 * time.Second)
		ch <- "PAYMENT_VERIFIED"
	}()

	select {
	case resuslt := <- ch:
		fmt.Println("Payment verification result:", resuslt)
		return nil
	case <- ctx.Done():
		fmt.Println("Payment verification canceled due to timeout.")
		return ctx.Err()
	}
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
func TestContextWithDeadline(t *testing.T) {

	// This context with deadline is good for handling long process that we want to limit the time to run
	// for example, we call the external API that sometimes take too long to respond
	// or we call database query that sometimes take too long to respond
	// the difference between timeout and deadline is
	// timeout is relative time from now
	// deadline is absolute time (example: 1st Jan 2025 10:00 AM)

	fmt.Println("Total goroutines start:", runtime.NumGoroutine())
	
	parentCtx := context.Background()

	// add deadline signal to context
	// after 5 seconds, the context will be canceled automatically eventhough the CreateCounter function is not finished yet
	ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5 * time.Second))
	defer cancel() // good practice to call cancel in defer to avoid memory leak eventhough the context already have timeout

	destination := CreateCounter(ctx)
	fmt.Println("Total goroutines running:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter:", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutines end:", runtime.NumGoroutine())
}


