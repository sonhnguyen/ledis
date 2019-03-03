package main

// Set type for data
type Set struct {
	Set map[string]bool
}

// SAdd add values to set stored at key
func (s *Set) SAdd(values []string) int {
	count := 0
	for _, v := range values {
		if !s.Set[v] {
			count++
			s.Set[v] = true
		}
	}
	return count
}

// SCard return the number of elements of the set stored at key
func (s *Set) SCard() int {
	count := 0
	for _, v := range s.Set {
		if v {
			count++
		}
	}
	return count
}

// SMembers return array of all members of set
func (s *Set) SMembers() []string {
	results := []string{}
	for k, v := range s.Set {
		if v {
			results = append(results, k)
		}
	}
	return results
}

// SRem remove values from set
func (s *Set) SRem(values []string) int {
	count := 0
	for _, v := range values {
		if s.Set[v] {
			count++
			s.Set[v] = false
		}
	}
	return count
}
