package imgui

func findMaxLen(lines []string) int {
	if len(lines) == 0 {
		return 0
	}
	if len(lines) == 1 {
		return internalLen(lines[0])
	}
	cur := internalLen(lines[0])
	for _, s := range lines {
		if internalLen(s) > cur {
			cur = internalLen(s)
		}
	}
	return cur
}
