package imgui

import "fmt"

type Stack struct {
	items []string
}

func (st *Stack) Push(s string) {
	st.items = append(st.items, fmt.Sprintf("%p", &s))
}

func (st *Stack) Pop() {
	if !st.IsEmpty() {
		st.items = st.items[:len(st.items)-1]
	}
}

func (st *Stack) Top() string {
	if !st.IsEmpty() {
		return st.items[len(st.items)-1]
	}
	return ""
}

func (st *Stack) IsEmpty() bool {
	return len(st.items) == 0
}
