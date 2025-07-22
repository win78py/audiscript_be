package util

import (
    "net/http"
    "time"
)

var DefaultHTTPClient = &http.Client{
    Timeout: 60 * time.Second,
}