package i

import "io"

// HttpRequester defines an interface for HTTP requests.
type HttpRequester interface {
	Post(uri string, body io.Reader) (io.Reader, error)
	Get(uri string) (io.Reader, error)
}
