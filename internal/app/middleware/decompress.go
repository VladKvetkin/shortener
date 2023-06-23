package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type deCompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newDecompressReader(r io.ReadCloser) (*deCompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &deCompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c deCompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c deCompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}

func DecompressBodyReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		contentEncoding := req.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			decompressReader, err := newDecompressReader(req.Body)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			req.Body = decompressReader
			defer decompressReader.Close()
		}

		next.ServeHTTP(resp, req)
	})
}
