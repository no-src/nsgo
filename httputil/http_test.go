package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testHttpClient HttpClient

	errInvalidControlCharacterInURL = errors.New("net/url: invalid control character in URL")
)

func TestHttpGet(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/get_hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	mux.HandleFunc("/get_world", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "world")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		path       string
		expectBody string
	}{
		{"/get_hello", "hello"},
		{"/get_world", "world"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpGet(reqUrl)
			if err != nil {
				t.Errorf("HttpGet: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpGet: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpGet: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpGetWithCookie(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	key := "key"
	mux.HandleFunc("/get_with_cookie", func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get(key)
		if len(v) == 0 {
			cookie, err := r.Cookie(key)
			if err == nil {
				v = cookie.Value
			}
		}
		fmt.Fprintf(w, v)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	header := http.Header{}
	header.Add(key, "hello")
	cookie := &http.Cookie{
		Name:  key,
		Value: "world",
	}
	testCases := []struct {
		name       string
		path       string
		header     http.Header
		cookie     *http.Cookie
		expectBody string
		expectErr  error
	}{
		{"return data from header", "/get_with_cookie", header, nil, "hello", nil},
		{"return data from cookie", "/get_with_cookie", http.Header{}, cookie, "world", nil},
		{"invalid url", "/get_with_cookie\t", http.Header{}, cookie, "", errInvalidControlCharacterInURL},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpGetWithCookie(reqUrl, tc.header, tc.cookie)
			if tc.expectErr == nil && err != nil {
				t.Errorf("HttpGetWithCookie: request error, url=%s err=%v", reqUrl, err)
				return
			}
			if tc.expectErr != nil {
				if err == nil {
					t.Errorf("HttpGetWithCookie: request error, expect get an error but get nil, url=%s", reqUrl)
				} else if !strings.Contains(err.Error(), tc.expectErr.Error()) {
					t.Errorf("HttpGetWithCookie: request error, expect get an error =>%v but get error =>%v, url=%s", tc.expectErr, err, reqUrl)
				}
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpGetWithCookie: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpGetWithCookie: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpPost(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	key := "key"
	mux.HandleFunc("/post_data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.FormValue(key))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name       string
		path       string
		expectBody string
	}{
		{"return data hello", "/post_data", "hello"},
		{"return data world", "/post_data", "world"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			reqData := url.Values{}
			reqData.Add(key, tc.expectBody)
			resp, err := testHttpClient.HttpPost(reqUrl, reqData)
			if err != nil {
				t.Errorf("HttpPost: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpPost: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpPost: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpPostWithCookie(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	key := "key"
	mux.HandleFunc("/post_data_with_cookie", func(w http.ResponseWriter, r *http.Request) {
		v := r.FormValue(key)
		cookie, err := r.Cookie(key)
		if err == nil {
			v = cookie.Value
		}
		fmt.Fprintf(w, v)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	reqData := url.Values{}
	reqData.Add(key, "hello")
	cookie := &http.Cookie{
		Name:  key,
		Value: "world",
	}
	testCases := []struct {
		name       string
		path       string
		cookie     *http.Cookie
		expectBody string
		expectErr  error
	}{
		{"return data from form", "/post_data_with_cookie", nil, "hello", nil},
		{"return data from cookie", "/post_data_with_cookie", cookie, "world", nil},
		{"invalid url", "/post_data_with_cookie\t", cookie, "", errInvalidControlCharacterInURL},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPostWithCookie(reqUrl, reqData, tc.cookie)
			if tc.expectErr == nil && err != nil {
				t.Errorf("HttpPostWithCookie: request error, url=%s err=%v", reqUrl, err)
				return
			}
			if tc.expectErr != nil {
				if err == nil {
					t.Errorf("HttpPostWithCookie: request error, expect get an error but get nil, url=%s", reqUrl)
				} else if !strings.Contains(err.Error(), tc.expectErr.Error()) {
					t.Errorf("HttpPostWithCookie: request error, expect get an error =>%v but get error =>%v, url=%s", tc.expectErr, err, reqUrl)
				}
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpPostWithCookie: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpPostWithCookie: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpPostFileChunkWithCookie(t *testing.T) {
	initDefaultClient()
	fieldName := "up_file"
	fileName := "hello.txt"
	chunk := []byte("some test contents")

	mux := http.NewServeMux()
	key := "key"
	mux.HandleFunc("/post_file_chunk_with_cookie", func(w http.ResponseWriter, r *http.Request) {
		v := r.FormValue(key)
		cookie, err := r.Cookie(key)
		if err == nil {
			v = cookie.Value
		}
		file, fh, err := r.FormFile(fieldName)
		if err == nil && fh.Filename == fileName {
			if vv, err := io.ReadAll(file); err == nil && len(vv) > 0 {
				v = string(vv)
			}
		}
		fmt.Fprintf(w, v)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	reqData := url.Values{}
	reqData.Add(key, "hello")
	cookie := &http.Cookie{
		Name:  key,
		Value: "world",
	}
	testCases := []struct {
		name       string
		path       string
		chunk      []byte
		cookie     *http.Cookie
		expectBody string
		expectErr  error
	}{
		{"return data from form", "/post_file_chunk_with_cookie", nil, nil, "hello", nil},
		{"return data from cookie", "/post_file_chunk_with_cookie", nil, cookie, "world", nil},
		{"return data from file", "/post_file_chunk_with_cookie", chunk, cookie, string(chunk), nil},
		{"invalid url", "/post_file_chunk_with_cookie\t", chunk, cookie, "", errInvalidControlCharacterInURL},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPostFileChunkWithCookie(reqUrl, fieldName, fileName, reqData, tc.chunk, tc.cookie)
			if tc.expectErr == nil && err != nil {
				t.Errorf("HttpPostFileChunkWithCookie: request error, url=%s err=%v", reqUrl, err)
				return
			}
			if tc.expectErr != nil {
				if err == nil {
					t.Errorf("HttpPostFileChunkWithCookie: request error, expect get an error but get nil, url=%s", reqUrl)
				} else if !strings.Contains(err.Error(), tc.expectErr.Error()) {
					t.Errorf("HttpPostFileChunkWithCookie: request error, expect get an error =>%v but get error =>%v, url=%s", tc.expectErr, err, reqUrl)
				}
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpPostFileChunkWithCookie: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpPostFileChunkWithCookie: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpPostWithoutRedirect(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	key := "key"
	mux.HandleFunc("/post_data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.FormValue(key))
	})
	mux.HandleFunc("/post_data_redirect_301", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/post_data", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/post_data_redirect_302", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/post_data", http.StatusFound)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	reqData := url.Values{}
	reqData.Add(key, "hello")

	testCases := []struct {
		name           string
		path           string
		expectCode     int
		expectLocation string
	}{
		{"return 301 code", "/post_data_redirect_301", http.StatusMovedPermanently, "/post_data"},
		{"return 302 code", "/post_data_redirect_302", http.StatusFound, "/post_data"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPostWithoutRedirect(reqUrl, reqData)
			if err != nil {
				t.Errorf("HttpPostWithoutRedirect: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.expectCode {
				t.Errorf("HttpPostWithoutRedirect: expect status code => %d, but actual status code => %d url=%s", tc.expectCode, resp.StatusCode, reqUrl)
				return
			}
			actual := resp.Header.Get("Location")
			if tc.expectLocation != actual {
				t.Errorf("HttpPostWithoutRedirect: expect Location => %s, but actual Location => %s url=%s", tc.expectLocation, actual, reqUrl)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	remotePath := fmt.Sprintf("/hello_%d.txt", time.Now().UnixMilli())
	localPath := "." + remotePath
	content := "hello world"
	mux.HandleFunc(remotePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, content)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	defer os.Remove(localPath)

	reqUrl := server.URL + remotePath

	testCases := []struct {
		name           string
		reqUrl         string
		alwaysDownload bool
	}{
		{"0.download file with empty url", "", false},
		{"1.download file", reqUrl, false},
		{"2.download file again", reqUrl, false},
		{"3.download file again force", reqUrl, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testHttpClient.Download(localPath, tc.reqUrl, tc.alwaysDownload)
			if len(tc.reqUrl) == 0 {
				if !errors.Is(err, errEmptyUrl) {
					t.Errorf("Download: expect get error %v but actual get %v, url=%s", errEmptyUrl, err, reqUrl)
				}
				return
			}
			if err != nil {
				t.Errorf("Download: request error, url=%s err=%v", tc.reqUrl, err)
				return
			}

			data, err := os.ReadFile(localPath)
			if err != nil {
				t.Errorf("Download: read file error, url=%s err=%v", tc.reqUrl, err)
				return
			}
			actual := string(data)
			if content != actual {
				t.Errorf("Download: expect content => %s, but actual content => %s url=%s", content, actual, tc.reqUrl)
			}
		})
	}
}

func TestNewHttpClient(t *testing.T) {
	testCases := []struct {
		name               string
		insecureSkipVerify bool
		certFile           string
		enableHTTP3        bool
		expectErr          bool
	}{
		{"disable verify and no cert file", true, "", false, false},
		{"enable verify and no cert file", false, "", false, true}, // return not exist error
		{"disable verify and use cert file", true, "./testdata/cert.pem", false, false},
		{"enable verify and use cert file", false, "./testdata/cert.pem", false, false},
		{"enable verify and use invalid cert file", false, "./testdata/key.pem", false, true}, // return errAppendCertsFromPemFailed error

		{"disable verify and no cert file with HTTP3", true, "", true, false},
		{"enable verify and no cert file with HTTP3", false, "", true, true}, // return not exist error
		{"disable verify and use cert file with HTTP3", true, "./testdata/cert.pem", true, false},
		{"enable verify and use cert file with HTTP3", false, "./testdata/cert.pem", true, false},
		{"enable verify and use invalid cert file with HTTP3", false, "./testdata/key.pem", true, true}, // return errAppendCertsFromPemFailed error
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewHttpClient(tc.insecureSkipVerify, tc.certFile, tc.enableHTTP3)
			if !tc.expectErr && err != nil {
				t.Errorf("Init: init http client error, err => %v", err)
				return
			}
			if tc.expectErr && (!os.IsNotExist(err) && !errors.Is(err, errAppendCertsFromPemFailed)) {
				t.Errorf("Init: init http client error, not get expect error, current get %v", err)
			}
		})
	}
}

func TestHttpPostData(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/post_data", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, string(data))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name       string
		path       string
		expectBody string
	}{
		{"return data hello", "/post_data", "hello"},
		{"return data world", "/post_data", "world"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPostData(reqUrl, []byte(tc.expectBody))
			if err != nil {
				t.Errorf("HttpPostData: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpPostData: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpPostData: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpPut(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/put_data", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, string(data))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name       string
		path       string
		expectBody string
	}{
		{"return data hello", "/put_data", "hello"},
		{"return data world", "/put_data", "world"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPut(reqUrl, []byte(tc.expectBody))
			if err != nil {
				t.Errorf("HttpPut: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpPut: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpPut: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func TestHttpDelete(t *testing.T) {
	initDefaultClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/delete_data", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, string(data))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name       string
		path       string
		expectBody string
	}{
		{"return data hello", "/delete_data", "hello"},
		{"return data world", "/delete_data", "world"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqUrl := server.URL + tc.path
			resp, err := testHttpClient.HttpPut(reqUrl, []byte(tc.expectBody))
			if err != nil {
				t.Errorf("HttpDelete: request error, url=%s err=%v", reqUrl, err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HttpDelete: read response body error, url=%s err=%v", reqUrl, err)
				return
			}
			actual := string(data)
			if tc.expectBody != actual {
				t.Errorf("HttpDelete: expect body => %s, but actual body => %s url=%s", tc.expectBody, actual, reqUrl)
			}
		})
	}
}

func initDefaultClient() {
	testHttpClient, _ = NewHttpClient(true, "", false)
}
