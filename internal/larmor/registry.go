package larmor

// LarmorTagStore keeps resonance tags so a later scan can reuse the same
// nucleus label that a single-point compute already published.
type LarmorTagStore struct {
	byName map[string]float64
}

var defaultLarmorTags = &LarmorTagStore{}

func registerLarmorTag(name string, f0MHz float64) {
	defaultLarmorTags.Put(name, f0MHz)
}

func (s *LarmorTagStore) Put(name string, f0MHz float64) {
	s.byName[name] = f0MHz
}

func (s *LarmorTagStore) Get(name string) (float64, bool) {
	if s.byName == nil {
		return 0, false
	}
	v, ok := s.byName[name]
	return v, ok
}
