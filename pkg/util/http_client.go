package util

import (
    "net/http"
    "time"
)

var DefaultHTTPClient = &http.Client{
    Timeout: 5 * time.Minute,
}