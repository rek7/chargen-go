package chargen

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	// Seed randomness for packet generation
	rand.Seed(time.Now().UnixNano())
}

func genData(num int) []byte {
	if num == 0 {
		num = rand.Intn(512-1) + 1
	}
	b := new(bytes.Buffer)
	for i := 0; num > i; i++ {
		b.Write([]byte(fmt.Sprintf("%c", rand.Intn(126-33)+3)))
	}
	return b.Bytes()
}
