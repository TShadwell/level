//+build !purego

package legacy

import (
	lvg "github.com/TShadwell/level/levigo"
)

func init() {
	lvl = lvg.Level()
}
