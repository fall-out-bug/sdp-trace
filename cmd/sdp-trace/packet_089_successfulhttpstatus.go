package main

func successfulHTTPStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
