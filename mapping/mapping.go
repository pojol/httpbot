package mapping

// Mapping mapping
type Mapping struct {
	state map[string]interface{}
}

// NewMapping new mapping
func NewMapping() *Mapping {
	return &Mapping{
		state: make(map[string]interface{}),
	}
}

// GetAll get mapping data
func (mp *Mapping) GetAll() map[string]interface{} {
	return mp.state
}

// Get get mapping data by key
func (mp *Mapping) Get(key string) interface{} {
	return mp.state[key]
}

// Set set data 2 mapping struct
func (mp *Mapping) Set(key string, val interface{}) {
	mp.state[key] = val
}
