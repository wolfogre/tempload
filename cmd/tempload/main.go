package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wolfogre/tempload/internal/pkg/bashupload"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(os.Args[0], "filename")
		return
	}
	for _, v := range os.Args[1:] {
		upload(v)
	}
}

func upload(name string) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	progress := bashupload.NewClient().Upload(filepath.Base(name), content)
	bar := pb.Full.Start(len(content))
	for p := range progress {
		if p.Done {
			bar.Finish()
			if p.Err != nil {
				fmt.Println(p.Err)
				os.Exit(1)
			}
			fmt.Println(p.Result)
			return

		}
		bar.SetCurrent(int64(p.Current))
	}
}
