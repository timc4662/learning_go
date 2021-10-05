package hello_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timc4662/learning_go/hello"
)

func TestTestify(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestHello(t *testing.T) {
	assert.Equal(t, "Don't communicate by sharing memory, share memory by communicating.", hello.DoHello())
}

func TestZeroValue(t *testing.T) {
	var i int
	assert.Zero(t, i, "is zero")

}

func TestDoHelloViaInterface(t *testing.T) {
	h := hello.Hello{
		Service: hello.NewDefaultService(),
	}
	result := h.DoHelloViaInterface()
	assert.Equal(t, "Don't communicate by sharing memory, share memory by communicating.", result)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Go() string {
	args := m.Called()
	return args.String(0)
}

func TestDoHelloViaInterfaceAndMockedService(t *testing.T) {
	// this shows how to use the mock framework to replace out the
	// Service with a MockService which returns foo bar and not the go motto.
	// So this would be a better test as it isolates the tested part to the hello rather than
	// the service code (or anything reaching out onto the network)
	ms := MockService{}
	ms.On("Go").Return("Foo Bar")
	h := hello.Hello{
		Service: &ms,
	}
	result := h.DoHelloViaInterface()
	assert.Equal(t, "Foo Bar", result)
}
