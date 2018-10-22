package cntr

func DistinctStrings(elems ...string) []string {
	m := map[string]bool{}
	for _, e := range elems {
		m[e] = true
	}
	var retElems []string
	for e, _ := range m {
		retElems = append(retElems, e)
	}
	return retElems
}
