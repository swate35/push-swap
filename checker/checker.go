package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

type Stack struct {
	Top    *Node
	Bottom *Node
	Size   int
}

func (s *Stack) PushTop(value int) {
	node := &Node{Value: value}

	if s.Size == 0 {
		s.Top = node
		s.Bottom = node
	} else {
		node.Next = s.Top
		s.Top.Prev = node
		s.Top = node
	}

	s.Size++
}

func (s *Stack) PopTop() (int, bool) {
	if s.Size == 0 {
		return 0, false
	}

	value := s.Top.Value

	if s.Size == 1 {
		s.Top = nil
		s.Bottom = nil
	} else {
		s.Top = s.Top.Next
		s.Top.Prev = nil
	}

	s.Size--

	return value, true
}

func (s *Stack) Swap() {
	if s.Size < 2 {
		return
	}

	first := s.Top
	second := first.Next

	first.Next = second.Next
	second.Prev = nil
	second.Next = first

	if first.Next != nil {
		first.Next.Prev = first
	} else {
		s.Bottom = first
	}

	first.Prev = second
	s.Top = second
}

func (s *Stack) Rotate() {
	if s.Size < 2 {
		return
	}

	first := s.Top

	s.Top = first.Next
	s.Top.Prev = nil

	first.Next = nil
	first.Prev = s.Bottom

	s.Bottom.Next = first
	s.Bottom = first
}

func (s *Stack) ReverseRotate() {
	if s.Size < 2 {
		return
	}

	last := s.Bottom

	s.Bottom = last.Prev
	s.Bottom.Next = nil

	last.Prev = nil
	last.Next = s.Top

	s.Top.Prev = last
	s.Top = last
}

func push(dest, src *Stack) {
	value, ok := src.PopTop()
	if !ok {
		return
	}
	dest.PushTop(value)
}

func isSorted(stack *Stack) bool {
	if stack.Size < 2 {
		return true
	}

	current := stack.Top

	for current.Next != nil {
		if current.Value > current.Next.Value {
			return false
		}
		current = current.Next
	}

	return true
}

func parseArguments(args []string) (*Stack, error) {
	stack := &Stack{}
	seen := make(map[int]struct{})

	for _, arg := range args {
		tokens := strings.Fields(arg)

		if len(tokens) == 0 {
			return nil, fmt.Errorf("empty argument")
		}

		for _, token := range tokens {
			value, err := strconv.ParseInt(token, 10, 32)
			if err != nil {
				return nil, err
			}

			num := int(value)

			if _, exists := seen[num]; exists {
				return nil, fmt.Errorf("duplicate")
			}

			seen[num] = struct{}{}
		}
	}

	// Conservation de l'ordre original
	var values []int

	for _, arg := range args {
		for _, token := range strings.Fields(arg) {
			num, _ := strconv.Atoi(token)
			values = append(values, num)
		}
	}

	for i := len(values) - 1; i >= 0; i-- {
		stack.PushTop(values[i])
	}

	return stack, nil
}

func applyInstruction(a, b *Stack, instruction string) bool {
	switch instruction {

	case "pa":
		push(a, b)

	case "pb":
		push(b, a)

	case "sa":
		a.Swap()

	case "sb":
		b.Swap()

	case "ss":
		a.Swap()
		b.Swap()

	case "ra":
		a.Rotate()

	case "rb":
		b.Rotate()

	case "rr":
		a.Rotate()
		b.Rotate()

	case "rra":
		a.ReverseRotate()

	case "rrb":
		b.ReverseRotate()

	case "rrr":
		a.ReverseRotate()
		b.ReverseRotate()

	default:
		return false
	}

	return true
}

func printError() {
	fmt.Fprintln(os.Stderr, "Error")
}

func main() {
	if len(os.Args) == 1 {
		return
	}

	stackA, err := parseArguments(os.Args[1:])
	if err != nil {
		printError()
		return
	}

	stackB := &Stack{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if !applyInstruction(stackA, stackB, scanner.Text()) {
			printError()
			return
		}
	}

	if scanner.Err() != nil {
		printError()
		return
	}

	if stackB.Size == 0 && isSorted(stackA) {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
