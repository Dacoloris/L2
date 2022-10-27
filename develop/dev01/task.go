package dev01

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"time"
)

func main() {
	ntpTime, err := ntp.Time("time.nist.gov")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	ntpTimeFormatted := ntpTime.Format(time.UnixDate)

	fmt.Printf("NTP Time: %v\n", ntpTime)
	fmt.Printf("NTP Unix Date Time: %v\n", ntpTimeFormatted)
}
