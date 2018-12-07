package promise

import (
	"errors"
	"testing"
	"time"
)

func TestRunningPromiseWithAsyncTask(t *testing.T) {
	p := NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
		// do something async
		go func() {
			time.Sleep(time.Second)
			resolve(42)
		}()
	})
	p.Then(func(response interface{}) interface{} {
		result := response.(int)
		if result != 42 {
			t.Errorf("Expected 42 but got: %v", result)
		}
		return nil
	}).Catch(func(err error) interface{} {

		// Catch should never been called since there is no error
		t.Fail()
		return nil
	})
}

func TestRunningPromiseWithSyncTask(t *testing.T) {
	p := NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
		resolve(42)
	})
	p.Then(func(response interface{}) interface{} {
		result := response.(int)
		if result != 42 {
			t.Errorf("Expected 42 but got: %v", result)
		}
		return nil
	}).Catch(func(err error) interface{} {
		t.Fail()
		return nil
	})
}

func TestFailures(t *testing.T) {
	p := NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
		reject(errors.New("Huston we have a problem"))
	})

	p.Then(func(data interface{}) interface{} {
		t.Fail()
		return nil
	}).Catch(func(err error) interface{} {
		if err.Error() != "Huston we have a problem" {
			t.Errorf("Expected 'Huston we have a problem' but got: %v", err)
		}
		return nil
	})
}

func TestSuccessfulChaining(t *testing.T) {
	p := NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
		resolve(42)
	})

	p.Catch(func(err error) interface{} {
		t.Fail()
		return nil
	}).Then(func(data interface{}) interface{} {
		if 42 != data.(int) {
			t.Errorf("Expected 42 but got: %v", data)
		}
		return 43
	}).Then(func(data interface{}) interface{} {
		if 43 != data.(int) {
			t.Errorf("Expected 43 but got: %v", data)
		}
		return nil
	})
}

func TestFailureChaining(t *testing.T) {
	p := NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
		resolve(42)
	})

	p.Then(func(data interface{}) interface{} {
		if 42 != data.(int) {
			t.Errorf("Expected 42 got: %v", data)
		}
		return errors.New("An Error")
	}).Catch(func(err error) interface{} {
		if err.Error() != "An Error" {
			t.Fail()
		}
		return nil
	})
}
