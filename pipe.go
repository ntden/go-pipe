package pipe

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Pipe contains the functions that need to be executed in order, where one's outputs are another's inputs (think of unix pipes).
type Pipe struct {
	funcs []interface{}
	mux   sync.Mutex
}

// New instantiates a new Pipe with initial functions in it.
func New(funcs ...interface{}) (*Pipe, error) {
	p := &Pipe{}
	for _, f := range funcs {
		if reflect.TypeOf(f).Kind() != reflect.Func {
			return nil, errors.New("argument is not a function")
		}
		p.funcs = append(p.funcs, f)
	}
	return p, nil
}

// Add can be used to insert an additional function to the end of the execution stack.
func (p *Pipe) Add(f interface{}) error {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return errors.New("argument is not a function")
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	p.funcs = append(p.funcs, f)
	return nil
}

// Execute loops through all internal functions and executes them in the order they were added.
//
// It executes functions one after the other, passing the outputs of one function as arguments
// to the next (the first function's arguments are the args passed to the Execute function).
//
// When a function returns an error and that error is not nil, it will be returned from Execute.
//
// Make sure the next function's signature is compatible with the current executing function.
// The following rules apply:
//
//  1. If a function has less arguments than the next one, an error is returned.
//  2. If a function has more arguments than the next one, only the first arguments thats
//     match the function's signature will be used.
//
// The last function's output will also be returned from the Execute function.
func (p *Pipe) Execute(args ...interface{}) ([]interface{}, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	var inputs []interface{} = args
	var outputs []interface{}
	var err error

	for _, fn := range p.funcs {
		// Determine the expected number of inputs.
		fnType := reflect.TypeOf(fn)
		numIn := fnType.NumIn()

		in := make([]reflect.Value, numIn)
		for j := 0; j < numIn; j++ {
			if j < len(inputs) {
				in[j] = reflect.ValueOf(inputs[j])
			} else {
				in[j] = reflect.Zero(reflect.TypeOf(fn).In(j))
			}
		}

		if len(inputs) < numIn {
			return nil, fmt.Errorf("not enough arguments for function %v", fnType)
		} else if len(inputs) > numIn {
			// Loop through the inputs to determine which ones match the expected types.
			newIn := make([]reflect.Value, numIn)
			var j int
			for i := 0; i < numIn; i++ {
				if inputs[i] != nil && reflect.TypeOf(inputs[i]).AssignableTo(reflect.TypeOf(in[i].Interface())) {
					newIn[i] = reflect.ValueOf(inputs[i])
					j++
				}
			}
			if j != numIn {
				return nil, fmt.Errorf("invalid arguments function %v", fn)
			}
			in = newIn
		}

		// Call the function with the determined arguments.
		out := reflect.ValueOf(fn).Call(in)

		// Store the outputs.
		for _, o := range out {
			if o.IsValid() {
				outputs = append(outputs, o.Interface())
				if o.Type().Name() == "error" && !o.IsNil() {
					err = o.Interface().(error)
					return nil, err
				}
			}
		}

		// Set the inputs for the next function.
		inputs = outputs
		outputs = []interface{}{}
	}

	return inputs, err
}
