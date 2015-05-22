// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import "testing"

// Used to view level while tweaking algorithm.
func TestDungeonGenerate(t *testing.T) {
	d := &dungeon{}
	d.Generate(80, 40)
	w, h := d.Size()
	if w != 81 || h != 41 {
		t.Error("Could not create dungeon")
	}
	// d.dump() // view level.
}
