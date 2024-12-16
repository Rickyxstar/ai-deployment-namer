package common

import "net/http"

type MockTransport struct {
	ResponseFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.ResponseFunc(req)
}

var _ http.RoundTripper = &MockTransport{}
