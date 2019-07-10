package bashupload

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/wolfogre/tempload/internal/pkg/filekeeper"
)

func NewClient() *_BashUpload {
	return &_BashUpload{}
}

const (
	rootPath = "https://bashupload.com/"
)

type _BashUpload struct {
}

func (*_BashUpload) Name() string {
	return "bashupload"
}

func (*_BashUpload) Ping() error {
	resp, err := http.Head(rootPath)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}
	return nil
}

func (u *_BashUpload) Upload(name string, content []byte) <-chan filekeeper.UploadProgress {
	ret := make(chan filekeeper.UploadProgress)
	go u.upload(name, content, ret)
	return ret
}

func (*_BashUpload) upload(name string, content []byte, progress chan filekeeper.UploadProgress) {
	progress <- filekeeper.UploadProgress{
		Total:   len(content),
		Current: 0,
		Done:    false,
		Result:  "",
		Err:     nil,
	}
	progressReader := filekeeper.NewProgressReader(content, func(n int) {
		progress <- filekeeper.UploadProgress{
			Total:   len(content),
			Current: n,
			Done:    false,
			Result:  "",
			Err:     nil,
		}
	})
	req, err := http.NewRequest(http.MethodPost, rootPath+name, progressReader)
	if err != nil {
		progress <- filekeeper.UploadProgress{
			Total:   len(content),
			Current: 0,
			Done:    true,
			Result:  "",
			Err:     err,
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		progress <- filekeeper.UploadProgress{
			Total:   len(content),
			Current: 0,
			Done:    true,
			Result:  "",
			Err:     err,
		}
		return
	}
	if resp.StatusCode != http.StatusOK {
		progress <- filekeeper.UploadProgress{
			Total:   len(content),
			Current: 0,
			Done:    true,
			Result:  "",
			Err:     errors.New(http.StatusText(resp.StatusCode)),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		progress <- filekeeper.UploadProgress{
			Total:   len(content),
			Current: 0,
			Done:    true,
			Result:  "",
			Err:     err,
		}
		return
	}

	progress <- filekeeper.UploadProgress{
		Total:   len(content),
		Current: len(content),
		Done:    true,
		Result:  string(body),
		Err:     nil,
	}
}
