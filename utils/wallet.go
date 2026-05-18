package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}
