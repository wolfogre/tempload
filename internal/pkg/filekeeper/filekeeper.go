package filekeeper

type UploadProgress struct {
	Total   int
	Current int
	Done    bool
	Result  string
	Err     error
}

type FileKeeper interface {
	Name() string
	Ping() error
	Upload(name string, content []byte) <-chan UploadProgress
}
