package openaiproxy

import (
	"bytes"
	"io"
	"net/http"
)

// TODO: check if you can process the streaming output as well
// for tool calling, modifying response etc. also check if you can use the official openapi client
// also there's some response headers that are set, including cookies, overwrite those
// https://openai.com/policies/cookie-policy/
type BodyProcessor interface {
	Process(body []byte) ([]byte, error)
}

// ProcessRequestBody is a helper function to process the request body
func ProcessRequestBody(req *http.Request, processor BodyProcessor) error {
	if req.Body == nil {
		return nil
	}

	// Read the original body
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// Close the original body
	req.Body.Close()

	// Process the body
	processedBody, err := processor.Process(bodyBytes)
	if err != nil {
		return err
	}

	// Recreate the body reader
	req.Body = io.NopCloser(bytes.NewBuffer(processedBody))

	return nil
}
