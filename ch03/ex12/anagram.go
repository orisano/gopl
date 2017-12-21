package main

func IsAnagram(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	acnt := map[int32]int{}
	for _, c := range a {
		acnt[c]++
	}
	bcnt := map[int32]int{}
	for _, c := range b {
		bcnt[c]++
	}
	if len(acnt) != len(bcnt) {
		return false
	}

	for k, v := range acnt {
		bv, ok := bcnt[k]
		if !ok {
			return false
		}
		if v != bv {
			return false
		}
	}
	return true
}
