package swn

import (
	"bufio"
	"context"
	"encoding/base64"

	"github.com/go-errors/errors"
)

func WriteB64(rw *bufio.ReadWriter, req []byte) error {
	encoded := base64.StdEncoding.EncodeToString(req)

	if _, err := rw.WriteString(encoded + "\n"); err != nil {
		return err
	}

	if err := rw.Flush(); err != nil {
		return err
	}

	return nil
}

func ReadB64(rw *bufio.ReadWriter) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AUTH_TIMEOUT)
	defer cancel()

	resultCh := make(chan []byte)
	errorCh := make(chan error, 1)

	go func() {
		// Read the encoded message from the stream.
		encodedReq, err := rw.ReadString('\n')
		if err != nil {
			errorCh <- err
			return
		}
		if len(encodedReq) == 0 {
			errorCh <- errors.Errorf("empty buffer")
			return
		}

		req, err := base64.StdEncoding.DecodeString(encodedReq)
		if err != nil {
			errorCh <- err
			return
		}
		if len(req) == 0 {
			errorCh <- errors.Errorf("empty buffer")
			return
		}

		resultCh <- req
	}()

	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errorCh:
		return nil, err
	case <-ctx.Done():
		return nil, errors.Errorf("auth timeout: %v sec passed", AUTH_TIMEOUT)
	}
}
