package uuid7

import (
	"crypto/rand"
	"fmt"
	"time"
)

type UUID7 [16]byte

// New menghasilkan UUID v7 baru
func New() string {
	var u UUID7

	// Timestamp 48-bit (milliseconds since Unix epoch)
	ts := uint64(time.Now().UnixMilli())
	u[0] = byte(ts >> 40)
	u[1] = byte(ts >> 32)
	u[2] = byte(ts >> 24)
	u[3] = byte(ts >> 16)
	u[4] = byte(ts >> 8)
	u[5] = byte(ts)

	// random 10 bytes
	_, err := rand.Read(u[6:])
	if err != nil {
		panic(err) // mirip MustGenerateUUID7
	}

	// set version (v7)
	u[6] = (u[6] & 0x0F) | 0x70

	// set variant (RFC 4122)
	u[8] = (u[8] & 0x3F) | 0x80

	return u.String()
}

// String mengembalikan UUID v7 dalam format standar
func (u UUID7) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:16],
	)
}
