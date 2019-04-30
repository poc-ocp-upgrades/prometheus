package testutil

import (
	"net/http"
)

type roundTrip struct {
	theResponse	*http.Response
	theError	error
}

func (rt *roundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rt.theResponse, rt.theError
}

type roundTripCheckRequest struct {
	checkRequest	func(*http.Request)
	roundTrip
}

func (rt *roundTripCheckRequest) RoundTrip(r *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rt.checkRequest(r)
	return rt.theResponse, rt.theError
}
func NewRoundTripCheckRequest(checkRequest func(*http.Request), theResponse *http.Response, theError error) http.RoundTripper {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &roundTripCheckRequest{checkRequest: checkRequest, roundTrip: roundTrip{theResponse: theResponse, theError: theError}}
}
