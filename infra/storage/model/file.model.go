package storage

import "io"

type File struct {
	Name    string
	Content io.Reader
	Bucket  string
}
