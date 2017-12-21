package main

func Unique(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	r := ss[:1]
	for _, s := range ss[1:] {
		if r[len(r)-1] != s {
			r = append(r, s)
		}
	}
	return r
}
