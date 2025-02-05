package model

import (
	"net"
	"time"
)

type Address struct {
	IP                 net.IP    `json:"ip"`
	LastSuccessfulPing time.Time `json:"last_successful_ping"`
	ResponseTime       uint      `json:"response_time"`
}

type GetAddressListResponse struct {
	Addresses []Address `json:"addresses"`
	Last      int       `json:"Last"`
	Page      int       `json:"Page`
}

type GetAddressListRequest struct {
	Page  int `json: "Page"`
	Ofset int `json: "Ofset", omitempty"`
	Limit int `json: "Limit", omitempty`
}

// type SetArticleResponse struct {
// 	Ofset string `json:"Ofset"`
// 	Limit string `json:"Limit"`
// }

type SetAddressListRequest struct {
	Addresses []Address `json:"addresses"`
}
