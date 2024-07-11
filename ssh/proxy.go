package ssh

import "math/rand/v2"

func randomPort() uint16 {
	return uint16(rand.Uint32())
}
