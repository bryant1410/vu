// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// aid.go defines asset identfiers. See eid.go for entity identifiers.

import (
	"math"
)

// asset describes any data asset that can uniquely identify itself.
type asset interface {
	aid() aid      // Data type and name combined.
	label() string // Identifier unique with data type.
}

// ============================================================================

// aid is a unique asset identifier.
// Asset identifiers are hashes generated from an asset name and type.
type aid uint64

// dataType returns the type of asset data for this aid.
func (a aid) dataType() uint32 { return uint32(a & math.MaxUint32) }

// Asset data types. See aid.dataType and assetID.
const (
	fnt = iota // font
	shd        // shader
	mat        // material
	msh        // mesh
	tex        // texture
	snd        // sound
	anm        // animation
)

// =============================================================================
// asset utility methods.

// assetID produces a unique asset identifier using for the given
// asset type t, and asset name.
func assetID(t int, name string) aid { return aid(t) + aid(stringHash(name))<<32 }

// stringHash turns a string into a number.
// Algorithm based on java String.hashCode().
//     s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]
func stringHash(s string) uint32 {
	bytes := []byte(s)
	n := len(bytes)
	hash := uint32(0)
	for index, b := range bytes {
		hash += uint32(b) * uint32(math.Pow(31, float64(n-index)))
	}
	return hash
}
