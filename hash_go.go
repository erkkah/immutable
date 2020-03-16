// +build go1.14,!cheaphash

package immutable

import (
	"hash/maphash"
)

var hashSeed = maphash.MakeSeed()

func hashFunc(bytes []byte) uint32 {
	var hash maphash.Hash

	hash.SetSeed(hashSeed)
	hash.Write(bytes)

	return uint32(hash.Sum64())
}
