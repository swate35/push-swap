package main

// stacks holds both stacks (index 0 is the top) and records every
// instruction that is executed on them.
type stacks struct {
	a, b         []int
	instructions []string
}

func swap(s []int) {
	if len(s) >= 2 {
		s[0], s[1] = s[1], s[0]
	}
}

func rotate(s []int) {
	if len(s) >= 2 {
		first := s[0]
		copy(s, s[1:])
		s[len(s)-1] = first
	}
}

func reverseRotate(s []int) {
	if len(s) >= 2 {
		last := s[len(s)-1]
		copy(s[1:], s)
		s[0] = last
	}
}

// push moves the top of from onto to and returns both stacks.
func push(from, to []int) ([]int, []int) {
	if len(from) == 0 {
		return from, to
	}
	to = append([]int{from[0]}, to...)
	return from[1:], to
}

func (s *stacks) record(name string) {
	s.instructions = append(s.instructions, name)
}

func (s *stacks) sa() { swap(s.a); s.record("sa") }
func (s *stacks) sb() { swap(s.b); s.record("sb") }
func (s *stacks) ss() { swap(s.a); swap(s.b); s.record("ss") }

func (s *stacks) pa() { s.b, s.a = push(s.b, s.a); s.record("pa") }
func (s *stacks) pb() { s.a, s.b = push(s.a, s.b); s.record("pb") }

func (s *stacks) ra() { rotate(s.a); s.record("ra") }
func (s *stacks) rb() { rotate(s.b); s.record("rb") }
func (s *stacks) rr() { rotate(s.a); rotate(s.b); s.record("rr") }

func (s *stacks) rra() { reverseRotate(s.a); s.record("rra") }
func (s *stacks) rrb() { reverseRotate(s.b); s.record("rrb") }
func (s *stacks) rrr() { reverseRotate(s.a); reverseRotate(s.b); s.record("rrr") }
