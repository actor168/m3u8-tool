package compose

import "github.com/actor168/m3u8-tool/pkg"

type Composer struct {
	M3U8 *pkg.M3U8
}

func (c *Composer) Compose(string, suffix string) bool {
	return false
}
