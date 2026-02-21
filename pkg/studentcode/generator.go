package studentcode

import (
	"fmt"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate 8 karakterli benzersiz öğrenci kodu üretir
// Format: YKS + 5 rakam  →  YKS12345
func Generate() string {
	return fmt.Sprintf("YKS%05d", rng.Intn(100000))
}
