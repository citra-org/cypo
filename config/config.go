package config

import (
	"sync"
	"time"
)

const (
    Port          = ":3000"
    OTP           = "123455"
    OTPExpiration = 5 * time.Minute
)

var (
    OtpIssued     time.Time
    Authenticated bool
    Mu            sync.Mutex
)
