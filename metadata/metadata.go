package metadata

import "strings"

const (
	UserId      = "userid"
	ParentName  = "parentname"
	ParentId    = "parentid"
	CreatedAt   = "createdat"
	FileType    = "filetype"
	DisplayName = "displayname"
	Checksum    = "checksum"
	Filename    = "filename"
)

type Metadata map[string]string
type Mapper func(data Metadata) Metadata

func (m Metadata) Get(key string) string {
	lk := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lk {
			return v
		}
	}
	return ""
}

func (m Metadata) Set(key, value string) {
	m[strings.ToLower(key)] = value
}

func (m Metadata) Map(from, to string) {
	m.Set(to, m.Get(from))
}
