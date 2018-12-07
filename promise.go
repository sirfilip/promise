package promise

import (
	"log"
	"sync"
)

type state int

const (
	pending state = iota
	resolved
	rejected
)

type ResolveFunc func(interface{}) interface{}
type RejectFunc func(error) interface{}

type Promise struct {
	sync.WaitGroup
	state state
	data  interface{}
	err   error
}

func (p *Promise) Resolved() bool {
	return p.state == resolved
}

func (p *Promise) Rejected() bool {
	return p.state == rejected
}

func (p *Promise) Data() interface{} {
	return p.data
}

func (p *Promise) Error() error {
	return p.err
}

func (p *Promise) Pending() bool {
	return p.state == pending
}

func (p *Promise) Resolve(data interface{}) interface{} {
	defer p.Done()
	if p.state != pending {
		log.Fatalf("Resolving promise in non pending state")
	}
	p.data = data
	p.state = resolved
	return nil
}

func (p *Promise) Reject(err error) interface{} {
	defer p.Done()
	if p.state != pending {
		log.Fatalf("Resolving promise in non pending state")
	}
	p.err = err
	p.state = rejected
	return nil
}

func (p *Promise) Then(resolver ResolveFunc) *Promise {
	p.Wait()
	if p.state != resolved {
		return p
	}
	switch result := resolver(p.data); response := result.(type) {
	case *Promise:
		return response
	case error:
		return &Promise{
			data:  nil,
			state: rejected,
			err:   response,
		}
	default:
		return &Promise{
			data:  response,
			state: resolved,
			err:   nil,
		}
	}
}

func (p *Promise) Catch(rejector RejectFunc) *Promise {
	p.Wait()
	if p.state == rejected {
		rejector(p.err)
	}
	return p
}

func NewPromise(callback func(ResolveFunc, RejectFunc)) *Promise {
	p := &Promise{state: pending}
	p.Add(1)
	go callback(p.Resolve, p.Reject)
	return p
}

func All(promises ...*Promise) []interface{} {
	result := make([]interface{}, 0, len(promises))
	wg := &sync.WaitGroup{}
	wg.Add(len(promises))
	for _, promise := range promises {
		go func(p *Promise) {
			defer wg.Done()
			p.Wait()
		}(promise)
	}
	wg.Wait()
	for _, p := range promises {
		if p.Resolved() {
			result = append(result, p.Data())
		} else {
			result = append(result, p.Error())
		}
	}
	return result
}
