package domain

import (
	"net"
	"time"
)

type Address struct {
	Id                 int
	IP                 net.IP
	LastSuccessfulPing time.Time
	ResponseTime       uint
}

// type AddressList struct {
// 	Addresses []Address
// 	Last      int
// 	Page      int
// }
