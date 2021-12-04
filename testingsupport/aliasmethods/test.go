package aliasmethods

type (
	Code int

	Properties map[string]interface{}

	StringList []string
)

func (c Code) AsInt() int {
	return int(c)
}

func (ps Properties) Copy() Properties {
	m := make(Properties)
	for key := range ps {
		m[key] = ps[key]
	}
	return m
}

func (sl StringList) Add(s string) StringList {
	return append(sl, s)
}
