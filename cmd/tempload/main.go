package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wolfogre/tempload/internal/pkg/bashupload"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(os.Args[0], "filename")
		return
	}
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	progress := bashupload.NewClient().Upload(filename, content)
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
