package main

import "net/url"

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	path := u.Path

	if path == "/" {
		path = ""
	} else if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return u.Host + path, nil
}
