package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateSimpleID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}
