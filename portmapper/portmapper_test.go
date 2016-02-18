package portmapper

import (
	"testing"
)

func TestHttp(t *testing.T) {
	pm := NewPortMapper()
	pm.AddPair(5050, 5050)
	pm.run()
}
