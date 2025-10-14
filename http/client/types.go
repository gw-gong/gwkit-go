package client

const (
	HTTPStatusUnknown = -1
)

const (
	HTTPMethodGet    = "GET"
	HTTPMethodPost   = "POST"
	HTTPMethodPut    = "PUT"
	HTTPMethodDelete = "DELETE"
	HTTPMethodPatch  = "PATCH"
)

type HeaderItem struct {
	Key   string
	Value string
}

// TODO: add socks5 config
// type SOCKS5Config struct {
// 	Address  string
// 	User     string
// 	Password string
// }
