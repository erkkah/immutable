// +build !go1.14 cheaphash

package immutable

func hashFunc(bytes []byte) uint32 {
	var hash uint32

	for _, byte := range bytes {
		hash = hash*31 + uint32(byte)
	}

	return hash
}
