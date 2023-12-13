package stack

import "fmt"

var stackSize int

type CommandsStack struct {
	commands []string
	Top int
}

func (s *CommandsStack) Push(value string) {
	if s.top == stackSize -1{
		fmt.Println("Cannot insert more values")
		return
	}

	s.commands = append(s.commands, value)
	s.top++
	fmt.Printf("Inserted value %v\n", value)
}

func (s *CommandsStack) Pop() string {
	if s.top == -1 {
		fmt.Println("Stack Underflow, No values to remove")
		return ""
	}

  	popedElement := s.commands[s.top]
	s.commands = s.commands[:s.top]
	s.top--

	return popedElement
}

func (s *CommandsStack) isEmpty() bool {
	if s.top == -1{
		fmt.Println("Stack is empty")
		return true
	}

	fmt.Println("Stack is not empty")

	return false
}

func (s *CommandsStack) isFull() bool {
	if s.top == stackSize -1{
		fmt.Println("Stack is full")
		return true
	}

	fmt.Println("Stack is not full")
	return false
}


func (s *CommandsStack) Peek() string {
	if s.top == -1 {
		fmt.Println("Stack is empty")
		return ""
	}

	fmt.Printf("Removed value %v\n", s.commands[s.top])
	element := s.commands[s.top]
	return element
}

func (s *CommandsStack) Display() {
	for _, v := range s.commands{
		fmt.Println(v)
	}
}