package main

import (
	"errors"
	"fmt"
)

const STACK_MAX = 256

type Object struct {
	Mark  int
	Next  *Object
	Value int
}

type VM struct {
	Stack       []*Object
	StackSize   int
	FirstObject *Object
	NumObjects  int
	MaxObjects  int
}

func NewVM() *VM {
	vm := &VM{Stack: make([]*Object, STACK_MAX), MaxObjects: 8}
	return vm
}

func (v *VM) Push(val *Object) error {
	if v.StackSize > STACK_MAX {
		return errors.New("Stack Overflow")
	}
	v.Stack[v.StackSize] = val
	v.StackSize += 1
	return nil
}

func (v *VM) Pop() (*Object, error) {
	if v.StackSize < 0 {
		return nil, errors.New("Stack Underflow")
	}
	obj := v.Stack[v.StackSize]
	v.StackSize -= 1
	return obj, nil
}

func Mark(obj *Object) {
	if obj.Mark == 1 {
		return
	}
	obj.Mark = 1
}

func (v *VM) MarkAll() {
	for i := 0; i < v.StackSize; i++ {
		Mark(v.Stack[i])
	}
}

func (v *VM) Sweep() {
	object := v.FirstObject
	prevObj := &Object{}
	for object != nil {
		if object.Mark != 1 {
			if prevObj.Next == nil {
				v.FirstObject = object.Next
			} else {
				prevObj.Next = object.Next
			}
			v.NumObjects -= 1
		} else {
			object.Mark = 0
			prevObj = object
		}
		object = object.Next
	}
}

func (v *VM) GC() {
	numObjects := v.NumObjects
	v.MarkAll()
	v.Sweep()
	v.MaxObjects = v.NumObjects * 2

	fmt.Printf("Collected %d objects , %d remaining \n", numObjects-v.NumObjects, v.NumObjects)
}

func NewObject(v *VM, val int) *Object {
	if v.NumObjects == v.MaxObjects {
		v.GC()
	}
	obj := &Object{Value: val}
	obj.Next = v.FirstObject
	v.FirstObject = obj
	obj.Mark = 0
	v.NumObjects += 1
	return obj
}

func (v *VM) PushInt(val int) {
	obj := NewObject(v, val)
	v.Push(obj)
}

func Test1() {
	fmt.Println("Test 1: Objects on Stack are preserved")
	vm := NewVM()
	vm.PushInt(1)
	vm.PushInt(2)
	vm.GC()
	if vm.NumObjects == 0 {
		fmt.Println("Should have collected objects")
	}
}

func Test2() {
	fmt.Println("Test 2 : Unreached objects are collected")
	vm := NewVM()
	vm.PushInt(1)
	vm.PushInt(2)
	vm.PushInt(3)
	vm.Pop()
	vm.Pop()
	vm.GC()
	if vm.NumObjects != 1 {
		fmt.Println("Should have collected some objects")
	}
}

func main() {
	Test1()
	Test2()
}
