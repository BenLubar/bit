package bitnum

type NumberMap struct {
	Keys   Numbers
	Values []interface{}
}

func (m *NumberMap) Get(k *Number) (v interface{}, ok bool) {
	i := m.Keys.Search(k)
	if i < len(m.Keys) && m.Keys[i].Equal(k) {
		return m.Values[i], true
	}
	return nil, false
}

func (m *NumberMap) Set(k *Number, v interface{}) {
	i := m.Keys.Search(k)
	if i >= len(m.Keys) || !m.Keys[i].Equal(k) {
		m.Keys = append(m.Keys, nil)
		m.Values = append(m.Values, nil)
		copy(m.Keys[i+1:], m.Keys[i:])
		copy(m.Values[i+1:], m.Values[i:])
		m.Keys[i] = k
	}
	m.Values[i] = v
}
