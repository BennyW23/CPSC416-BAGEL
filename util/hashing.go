package util

import (
	"encoding/binary"
	"hash/fnv"
)

func HashId(vertexId uint64) int64 {
	inputBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(inputBytes, vertexId)

	algorithm := fnv.New64a()
	algorithm.Write(inputBytes)
	return int64(algorithm.Sum64())
}

func GetFlooredModulo(a int64, b int64) int64 {
	return ((a % b) + b) % b
}
