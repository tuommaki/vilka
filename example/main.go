package main

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/tuommaki/vilka"
)

var (
	// ErrNoParams occurs when function doesn't get enough parameters
	ErrNoParams = errors.New("no params")

	// ErrDivByZero if denominator is zero
	ErrDivByZero = errors.New("division by zero")
)

// Operation type for message
type Operation int

const (
	// Addition operation
	Add Operation = iota

	// Subtract operation
	Sub Operation = iota

	// Multiply operation
	Mul Operation = iota

	// Divide operation
	Div Operation = iota

	// Result operation
	Result Operation = iota

	// Quit program operation
	Quit Operation = iota
)

type Message struct {
	Op     Operation
	Params []float64
	Error  error
}

func Calculator(mc vilka.MsgChan) int {
	defer mc.Close()

	for {
		msg := Message{}
		err := mc.Recv(&msg)
		if err != nil {
			log.Fatal(err)
		}

		var res float64

		switch msg.Op {
		case Add:
			res, err = add(msg.Params)

		case Sub:
			res, err = sub(msg.Params)

		case Mul:
			res, err = mul(msg.Params)

		case Div:
			res, err = mul(msg.Params)

		case Result:
			err = errors.New("invalid format")

		case Quit:
			return 0
		}

		msg.Op = Result
		msg.Params = []float64{res}
		msg.Error = err

		err = mc.Send(&msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	return 1
}

func add(xs []float64) (float64, error) {
	if len(xs) < 2 {
		return math.NaN(), ErrNoParams
	}

	res := xs[0]
	for _, x := range xs[1:] {
		res += x
	}

	return res, nil
}

func sub(xs []float64) (float64, error) {
	if len(xs) < 2 {
		return math.NaN(), ErrNoParams
	}

	res := xs[0]
	for _, x := range xs[1:] {
		res -= x
	}

	return res, nil
}

func mul(xs []float64) (float64, error) {
	if len(xs) < 2 {
		return math.NaN(), ErrNoParams
	}

	res := xs[0]
	for _, x := range xs[1:] {
		res *= x
	}

	return res, nil
}

func div(xs []float64) (float64, error) {
	if len(xs) != 2 {
		return math.NaN(), ErrNoParams
	}

	if xs[1] == 0.0 {
		return math.NaN(), ErrDivByZero
	}

	return (xs[0] / xs[1]), nil
}

func init() {
	vilka.Register("calculator", Calculator)
}

func main() {
	// Execute actual sub-command or continue execution
	// When execution branches to sub-command, it will never return
	vilka.MaybeExec()

	// In main program
	fmt.Printf("Calculator 5000\n")

	// Launch our child program to execute math operations
	mc, err := vilka.Launch("calculator")
	if err != nil {
		log.Fatal(err)
	}

	msg := Message{
		Op:     Add,
		Params: []float64{1.0, 2.0, 3.0},
	}

	err = mc.Send(&msg)
	if err != nil {
		log.Fatal(err)
	}

	err = mc.Recv(&msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("1.0 + 2.0 + 3.0 = %f\n", msg.Params[0])

	msg.Op = Quit
	err = mc.Send(&msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Child is ought to exit now.")
	fmt.Println("Closing the shop...")
}
