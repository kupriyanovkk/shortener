package urladdress

import (
	"errors"
	"net"
	"net/url"
)

type URLAddress struct {
	Host string
	Port string
}

func (n URLAddress) String() string {
	return n.Host + ":" + n.Port
}

func (n *URLAddress) Set(s string) error {
	url, parseErr := url.ParseRequestURI(s)
	if parseErr != nil {
		return errors.New(parseErr.Error())
	}

	host, port, _ := net.SplitHostPort(url.Host)

	n.Host = host
	n.Port = port
	return nil
}

func New() URLAddress {
	return URLAddress{}
}
