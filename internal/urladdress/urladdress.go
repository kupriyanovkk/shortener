package urladdress

import (
	"errors"
	"net"
	"net/url"
)

type UrlAddress struct {
	Host string
	Port string
}

func (n UrlAddress) String() string {
	return n.Host + ":" + n.Port
}

func (n *UrlAddress) Set(s string) error {
	url, parseErr := url.ParseRequestURI(s)
	if parseErr != nil {
		return errors.New(parseErr.Error())
	}

	host, port, _ := net.SplitHostPort(url.Host)

	n.Host = host
	n.Port = port
	return nil
}

func New() UrlAddress {
	return UrlAddress{}
}
