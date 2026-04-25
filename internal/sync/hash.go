package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashEvent(event Event) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s",
		event.Title,
		event.Start.Format("2006-01-02T15:04:05Z"),
		event.End.Format("2006-01-02T15:04:05Z"),
		event.Description,
		event.Location,
	)
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
