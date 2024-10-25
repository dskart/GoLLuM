package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

// https://platform.openai.com/docs/api-reference/chat/create

type RequestBody interface {
	ChatCompletionRequestBody
}

type ResponseObject interface {
	ChatCompletionObject
}

func request[B RequestBody, R ResponseObject](ctx context.Context, httpClient *retryablehttp.Client, url string, apiKey string, reqBody B) (*R, error) {
	rawBody, err := json.Marshal(&reqBody)
	if err != nil {
		return nil, err
	}

	r, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, rawBody)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respData))
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respObject R
	if err := json.Unmarshal(respData, &respObject); err != nil {
		return nil, err
	}

	return &respObject, nil
}
