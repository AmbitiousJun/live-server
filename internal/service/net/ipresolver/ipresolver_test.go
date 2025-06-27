package ipresolver_test

import (
	"fmt"
	"testing"

	"github.com/AmbitiousJun/live-server/internal/service/net/ipresolver"
)

func TestIpwcn(t *testing.T) {
	r := ipresolver.V6()
	fmt.Println(r.Resolve("2606:4700:8d77:b110:d58e:5cf3:ff93:9644"))
}
