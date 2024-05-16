package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/url"
	"os"
)

const (
	// HeaderContentType the Content-Type http header
	HeaderContentType = "Content-Type"
)

var (
	// errAppendCertsFromPemFailed attempts to parse a series of PEM encoded certificates failed
	errAppendCertsFromPemFailed = errors.New("append certs from pem failed")
	errEmptyUrl                 = errors.New("url is empty")
)

// HttpClient an HTTP client component that supports chunked file uploads using the POST method
type HttpClient interface {
	// HttpGet get http resource
	HttpGet(url string) (resp *http.Response, err error)
	// HttpGetWithCookie get http resource with cookies
	HttpGetWithCookie(url string, header http.Header, cookies ...*http.Cookie) (resp *http.Response, err error)
	// HttpPost send a post request with form data
	HttpPost(url string, data url.Values) (resp *http.Response, err error)
	// HttpPostWithCookie send a post request with form data and cookies
	HttpPostWithCookie(url string, data url.Values, cookies ...*http.Cookie) (resp *http.Response, err error)
	// HttpPostFileChunkWithCookie send a post request with form data, a file chunk and cookies
	HttpPostFileChunkWithCookie(url string, fieldName string, fileName string, data url.Values, chunk []byte, cookies ...*http.Cookie) (resp *http.Response, err error)
	// HttpPostWithoutRedirect send a post request with form data and not auto redirect
	HttpPostWithoutRedirect(url string, data url.Values) (resp *http.Response, err error)
	// Download if the local file does not exist and the alwaysDownload is false, downloads the remote file to local path
	Download(path, url string, alwaysDownload bool) error
	// HttpPostData send a post request with data
	HttpPostData(url string, data []byte) (resp *http.Response, err error)
	// HttpPut send a put request with data
	HttpPut(url string, data []byte) (resp *http.Response, err error)
	// HttpDelete send a delete request with data
	HttpDelete(url string, data []byte) (resp *http.Response, err error)
}

// NewTLSConfig create a tls config
func NewTLSConfig(insecureSkipVerify bool, certFile string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	if !insecureSkipVerify {
		roots := x509.NewCertPool()
		pemCerts, err := os.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
		if !roots.AppendCertsFromPEM(pemCerts) {
			return nil, errAppendCertsFromPemFailed
		}
		tlsConfig.RootCAs = roots
	}
	return tlsConfig, nil
}
