// hello package is simply to try out the testify assertion based testing library.
package hello

import (
	"rsc.io/quote"
)

type Service interface {
	Go() string
}

type defaultService struct {
	// normally this might have stuff on it (connection props etc)
}

func (ds *defaultService) Go() string {
	return quote.Go()
}

// NewDefaultService returns a default service
func NewDefaultService() Service {
	return &defaultService{}
}

type Hello struct {
	Service
}

// DoHelloViaInterface returns go's motto via the service interface
func (h *Hello) DoHelloViaInterface() string {
	return h.Service.Go()
}

// DoHello returns go's motto.
func DoHello() string {
	return quote.Go()
}
