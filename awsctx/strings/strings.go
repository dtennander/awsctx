package strings

func UnionOf(as []string, bs []string) []string {
	unionMap := map[string]bool{}
	for _, a := range as {
		unionMap[a] = true
	}
	for _, b := range bs {
		unionMap[b] = true
	}
	var union []string
	for k := range unionMap {
		union = append(union, k)
	}
	return union
}

func Contains(ls []string, s string) bool {
	m := map[string]bool{}
	for _, l := range ls {
		m[l] = true
	}
	return m[s]
}
