package utils

import (
	"fmt"
	"hash/crc32"
)

func GenerateName(prefix, from string) string {
	cs := crc32.Checksum([]byte(from), crc32.IEEETable)

	return fmt.Sprintf("%s%08x", prefix, cs)
}
