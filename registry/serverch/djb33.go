package registry

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
)

var djb33Seed uint32

func init() {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("random seed error %v", err)
	}
	djb33Seed = uint32(rnd.Uint64())
}

func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
	)
	if l > 0 {
		d = djb33Loop(d, k)
	}
	return d ^ (d >> 16)
}

func djb33Loop(d uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		i = uint32(0)
	)
loop:
	if i >= l-1 {
		goto exit
	}
	d = (d * 33) ^ uint32(k[i])
	i++
	goto loop
exit:
	return d
}
