// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/shard.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package util

import (
	"crypto/md5"
	"encoding/binary"
	"math"
	"unsafe"
)

const (
	// Sharding strategies.
	ShardingStrategyDefault = "default"
	ShardingStrategyShuffle = "shuffle-sharding"
)

var (
	seedSeparator = []byte{0}
)

// ShuffleShardSeed returns seed for random number generator, computed from provided identifier.
func ShuffleShardSeed(identifier, zone string) int64 {
	// Use the identifier to compute an hash we'll use to seed the random.
	hasher := md5.New()               //nolint:gosec
	hasher.Write(yoloBuf(identifier)) // nolint:errcheck
	if zone != "" {
		hasher.Write(seedSeparator) // nolint:errcheck
		hasher.Write(yoloBuf(zone)) // nolint:errcheck
	}
	checksum := hasher.Sum(nil)

	// Generate the seed based on the first 64 bits of the checksum.
	return int64(binary.BigEndian.Uint64(checksum))
}

// ShuffleShardExpectedInstancesPerZone returns the number of instances that should be selected for each
// zone when zone-aware replication is enabled. The algorithm expects the shard size to be divisible
// by the number of zones, in order to have nodes balanced across zones. If it's not, we do round up.
func ShuffleShardExpectedInstancesPerZone(shardSize, numZones int) int {
	return int(math.Ceil(float64(shardSize) / float64(numZones)))
}

// ShuffleShardExpectedInstances returns the total number of instances that should be selected for a given
// tenant. If zone-aware replication is disabled, the input numZones should be 1.
func ShuffleShardExpectedInstances(shardSize, numZones int) int {
	return ShuffleShardExpectedInstancesPerZone(shardSize, numZones) * numZones
}

func yoloBuf(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s)) // nolint:gosec
}
