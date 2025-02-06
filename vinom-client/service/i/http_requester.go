package i

import "io"

// HttpRequester defines an interface for HTTP requests.
type HttpRequester interface {
	Post(uri string, body io.Reader, authToken string) (io.Reader, error)
	Get(uri, authToken string) (io.Reader, error)
}
