package metadata

import "strings"

const (
	// The ID of the user the file belongs to.
	UserId = "userid"

	// The name of the container the file belongs to.
	ParentName = "parentname"

	// The ID of the container the file belongs to.
	ParentId = "parentid"

	// The RFC3339 compliant timestamp at which the file was created.
	CreatedAt = "createdat"

	// The mime type of the file.
	FileType = "filetype"

	// The name given to the file by the user.
	DisplayName = "displayname"

	// The checksum of the file.
	Checksum = "checksum"

	// The name of the file on the file system.
	Filename = "filename"
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
