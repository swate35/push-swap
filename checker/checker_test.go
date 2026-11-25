package main

import (
	"bytes"
	"strings"
	"testing"
)

func stackValues(stack *Stack) []int {
	values := make([]int, 0, stack.Size)
	for node := stack.Top; node != nil; node = node.Next {
		values = append(values, node.Value)
	}
	return values
}

func valuesEqual(got []int, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func newStack(values ...int) *Stack {
	stack := &Stack{}
	for i := len(values) - 1; i >= 0; i-- {
		stack.PushTop(values[i])
	}
	return stack
}

func TestPushTopAndPopTop(t *testing.T) {
	stack := &Stack{}
	stack.PushTop(10)
	stack.PushTop(20)
	stack.PushTop(30)

	if got, want := stackValues(stack), []int{30, 20, 10}; !valuesEqual(got, want) {
		t.Fatalf("PushTop values = %v, want %v", got, want)
	}

	value, ok := stack.PopTop()
	if !ok || value != 30 {
		t.Fatalf("PopTop = (%d, %v), want (30, true)", value, ok)
	}
	if got, want := stackValues(stack), []int{20, 10}; !valuesEqual(got, want) {
		t.Fatalf("values after PopTop = %v, want %v", got, want)
	}

	empty := &Stack{}
	if value, ok := empty.PopTop(); ok || value != 0 {
		t.Fatalf("PopTop on empty = (%d, %v), want (0, false)", value, ok)
	}
}

func TestSwap(t *testing.T) {
	stack := newStack(1, 2, 3)
	stack.Swap()

	if got, want := stackValues(stack), []int{2, 1, 3}; !valuesEqual(got, want) {
		t.Fatalf("Swap values = %v, want %v", got, want)
	}
	if stack.Bottom.Value != 3 || stack.Bottom.Prev.Value != 1 {
		t.Fatalf("Swap left incorrect bottom links")
	}

	empty := &Stack{}
	empty.Swap()
	if empty.Size != 0 {
		t.Fatalf("Swap changed an empty stack")
	}
}

func TestRotate(t *testing.T) {
	stack := newStack(1, 2, 3)
	stack.Rotate()

	if got, want := stackValues(stack), []int{2, 3, 1}; !valuesEqual(got, want) {
		t.Fatalf("Rotate values = %v, want %v", got, want)
	}
	if stack.Bottom.Value != 1 || stack.Bottom.Prev.Value != 3 {
		t.Fatalf("Rotate left incorrect bottom links")
	}

	one := newStack(42)
	one.Rotate()
	if got, want := stackValues(one), []int{42}; !valuesEqual(got, want) {
		t.Fatalf("Rotate on one element = %v, want %v", got, want)
	}
}

func TestReverseRotate(t *testing.T) {
	stack := newStack(1, 2, 3)
	stack.ReverseRotate()

	if got, want := stackValues(stack), []int{3, 1, 2}; !valuesEqual(got, want) {
		t.Fatalf("ReverseRotate values = %v, want %v", got, want)
	}
	if stack.Top.Value != 3 || stack.Top.Next.Value != 1 {
		t.Fatalf("ReverseRotate left incorrect top links")
	}

	one := newStack(42)
	one.ReverseRotate()
	if got, want := stackValues(one), []int{42}; !valuesEqual(got, want) {
		t.Fatalf("ReverseRotate on one element = %v, want %v", got, want)
	}
}

func TestPushMovesTopBetweenStacks(t *testing.T) {
	a := newStack(1, 2)
	b := newStack(9)

	push(b, a)
	if got, want := stackValues(a), []int{2}; !valuesEqual(got, want) {
		t.Fatalf("stack A after push = %v, want %v", got, want)
	}
	if got, want := stackValues(b), []int{1, 9}; !valuesEqual(got, want) {
		t.Fatalf("stack B after push = %v, want %v", got, want)
	}

	push(a, b)
	if got, want := stackValues(a), []int{1, 2}; !valuesEqual(got, want) {
		t.Fatalf("stack A after second push = %v, want %v", got, want)
	}
	if got, want := stackValues(b), []int{9}; !valuesEqual(got, want) {
		t.Fatalf("stack B after second push = %v, want %v", got, want)
	}

	emptySource := &Stack{}
	destination := newStack(7)
	push(destination, emptySource)
	if got, want := stackValues(destination), []int{7}; !valuesEqual(got, want) {
		t.Fatalf("push from empty source changed destination: %v, want %v", got, want)
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   bool
	}{
		{"empty", nil, true},
		{"one element", []int{5}, true},
		{"ascending", []int{-3, 0, 8, 12}, true},
		{"descending", []int{12, 8, 0, -3}, false},
		{"not fully sorted", []int{1, 3, 2}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stack := &Stack{}
			for i := len(test.values) - 1; i >= 0; i-- {
				stack.PushTop(test.values[i])
			}
			if got := isSorted(stack); got != test.want {
				t.Fatalf("isSorted(%v) = %v, want %v", test.values, got, test.want)
			}
		})
	}
}

func TestParseArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []int
	}{
		{"one string", []string{"3 2 1 0"}, []int{3, 2, 1, 0}},
		{"separate arguments", []string{"3", "2", "1", "0"}, []int{3, 2, 1, 0}},
		{"mixed arguments", []string{"3 2", "1", "0"}, []int{3, 2, 1, 0}},
		{"negative values", []string{"-2", "-10", "0"}, []int{-2, -10, 0}},
		{"32 bit limits", []string{"2147483647", "-2147483648"}, []int{2147483647, -2147483648}},
		{"extra whitespace", []string{"  3   2\t1  "}, []int{3, 2, 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stack, err := parseArguments(test.args)
			if err != nil {
				t.Fatalf("parseArguments(%v) returned error: %v", test.args, err)
			}
			if got := stackValues(stack); !valuesEqual(got, test.want) {
				t.Fatalf("parsed values = %v, want %v", got, test.want)
			}
		})
	}
}

func TestParseArgumentsErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"not integer", []string{"3", "one", "1"}},
		{"decimal", []string{"1.5"}},
		{"duplicate in one string", []string{"1 2 1"}},
		{"duplicate separate arguments", []string{"1", "1"}},
		{"too large", []string{"2147483648"}},
		{"too small", []string{"-2147483649"}},
		{"empty argument", []string{""}},
		{"whitespace only", []string{"   "}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stack, err := parseArguments(test.args)
			if err == nil {
				t.Fatalf("parseArguments(%v) = %v, want error", test.args, stack)
			}
		})
	}
}

func TestApplyInstructionSingleOperations(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		a           []int
		b           []int
		wantA       []int
		wantB       []int
	}{
		{"pa", "pa", []int{2, 3}, []int{1, 4}, []int{1, 2, 3}, []int{4}},
		{"pb", "pb", []int{1, 2, 3}, []int{4}, []int{2, 3}, []int{1, 4}},
		{"sa", "sa", []int{1, 2, 3}, nil, []int{2, 1, 3}, nil},
		{"sb", "sb", []int{1}, []int{4, 5, 6}, []int{1}, []int{5, 4, 6}},
		{"ra", "ra", []int{1, 2, 3}, nil, []int{2, 3, 1}, nil},
		{"rb", "rb", nil, []int{4, 5, 6}, nil, []int{5, 6, 4}},
		{"rra", "rra", []int{1, 2, 3}, nil, []int{3, 1, 2}, nil},
		{"rrb", "rrb", nil, []int{4, 5, 6}, nil, []int{6, 4, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := newStack(test.a...)
			b := newStack(test.b...)

			if ok := applyInstruction(a, b, test.instruction); !ok {
				t.Fatalf("applyInstruction rejected %q", test.instruction)
			}
			if got := stackValues(a); !valuesEqual(got, test.wantA) {
				t.Fatalf("A after %s = %v, want %v", test.instruction, got, test.wantA)
			}
			if got := stackValues(b); !valuesEqual(got, test.wantB) {
				t.Fatalf("B after %s = %v, want %v", test.instruction, got, test.wantB)
			}
		})
	}
}

func TestApplyInstructionCombinedOperations(t *testing.T) {
	a := newStack(1, 2, 3)
	b := newStack(4, 5, 6)

	if !applyInstruction(a, b, "ss") {
		t.Fatal("ss should be a valid instruction")
	}
	if got, want := stackValues(a), []int{2, 1, 3}; !valuesEqual(got, want) {
		t.Fatalf("A after ss = %v, want %v", got, want)
	}
	if got, want := stackValues(b), []int{5, 4, 6}; !valuesEqual(got, want) {
		t.Fatalf("B after ss = %v, want %v", got, want)
	}

	a = newStack(1, 2, 3)
	b = newStack(4, 5, 6)
	if !applyInstruction(a, b, "rr") {
		t.Fatal("rr should be a valid instruction")
	}
	if got, want := stackValues(a), []int{2, 3, 1}; !valuesEqual(got, want) {
		t.Fatalf("A after rr = %v, want %v", got, want)
	}
	if got, want := stackValues(b), []int{5, 6, 4}; !valuesEqual(got, want) {
		t.Fatalf("B after rr = %v, want %v", got, want)
	}

	a = newStack(1, 2, 3)
	b = newStack(4, 5, 6)
	if !applyInstruction(a, b, "rrr") {
		t.Fatal("rrr should be a valid instruction")
	}
	if got, want := stackValues(a), []int{3, 1, 2}; !valuesEqual(got, want) {
		t.Fatalf("A after rrr = %v, want %v", got, want)
	}
	if got, want := stackValues(b), []int{6, 4, 5}; !valuesEqual(got, want) {
		t.Fatalf("B after rrr = %v, want %v", got, want)
	}
}

func TestApplyInstructionRejectsUnknown(t *testing.T) {
	a := newStack(1, 2)
	b := &Stack{}

	if applyInstruction(a, b, "unknown") {
		t.Fatal("unknown instruction should be rejected")
	}
	if applyInstruction(a, b, " sa") {
		t.Fatal("instruction with whitespace should be rejected")
	}
	if applyInstruction(a, b, "") {
		t.Fatal("empty instruction should be rejected")
	}
}

func TestApplyInstructionNoEffectWithSmallStacks(t *testing.T) {
	valid := []string{"sa", "sb", "ss", "ra", "rb", "rr", "rra", "rrb", "rrr", "pa"}
	for _, instruction := range valid {
		testA := newStack(1)
		testB := &Stack{}
		if !applyInstruction(testA, testB, instruction) {
			t.Fatalf("%q should be valid", instruction)
		}
		if got, want := stackValues(testA), []int{1}; !valuesEqual(got, want) {
			t.Fatalf("A changed unexpectedly after %s: %v, want %v", instruction, got, want)
		}
		if testB.Size != 0 {
			t.Fatalf("B changed unexpectedly after %s: %v", instruction, stackValues(testB))
		}
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		input  string
		stdout string
		stderr string
		code   int
	}{
		{
			name: "no arguments",
			code: 0,
		},
		{
			name:   "already sorted",
			args:   []string{"1", "2", "3"},
			stdout: "OK\n",
			code:   0,
		},
		{
			name:   "unsorted",
			args:   []string{"3", "2", "1"},
			stdout: "KO\n",
			code:   0,
		},
		{
			name:   "sorted after instructions",
			args:   []string{"3", "2", "1", "0"},
			input:  "rra\npb\nsa\nrra\npa\n",
			stdout: "OK\n",
			code:   0,
		},
		{
			name:   "known instruction leaving B not empty",
			args:   []string{"3", "2", "1"},
			input:  "pb\n",
			stdout: "KO\n",
			code:   0,
		},
		{
			name:   "invalid argument",
			args:   []string{"3", "one", "1"},
			stderr: "Error\n",
			code:   1,
		},
		{
			name:   "duplicate argument",
			args:   []string{"3", "2", "3"},
			stderr: "Error\n",
			code:   1,
		},
		{
			name:   "integer overflow",
			args:   []string{"2147483648"},
			stderr: "Error\n",
			code:   1,
		},
		{
			name:   "unknown instruction",
			args:   []string{"3", "2", "1"},
			input:  "nop\n",
			stderr: "Error\n",
			code:   1,
		},
		{
			name:   "incorrectly formatted instruction",
			args:   []string{"3", "2", "1"},
			input:  " sa \n",
			stderr: "Error\n",
			code:   1,
		},
		{
			name:   "empty instruction line",
			args:   []string{"3", "2", "1"},
			input:  "\n",
			stderr: "Error\n",
			code:   1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := run(test.args, strings.NewReader(test.input), &stdout, &stderr)

			if code != test.code {
				t.Fatalf("exit code = %d, want %d", code, test.code)
			}
			if got := stdout.String(); got != test.stdout {
				t.Fatalf("stdout = %q, want %q", got, test.stdout)
			}
			if got := stderr.String(); got != test.stderr {
				t.Fatalf("stderr = %q, want %q", got, test.stderr)
			}
		})
	}
}
