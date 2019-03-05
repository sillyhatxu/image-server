package uuid

import (
	"github.com/satori/go.uuid"
)

func UUID() string {
	u := uuid.Must(uuid.NewV4())
	return u.String()
}
