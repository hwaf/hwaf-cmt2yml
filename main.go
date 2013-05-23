package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func handle_err(err error) {
	if err != nil {
		panic(fmt.Errorf("hwaf-cmt2yml: %v", err))
	}
}

func main() {
	fmt.Printf("::: hwaf-cmt2yml\n")

	dir := "."
	switch len(os.Args) {
	case 1:
		dir = "."
	case 2:
		dir = os.Args[1]
	default:
		panic(fmt.Errorf("cmt2yml takes at most 1 argument (got %d)", len(os.Args)))
	}

	var err error
	//dir, err = filepath.Abs(dir)
	handle_err(err)

	fnames := []string{}
	fmt.Printf(">>> dir=%q\n", dir)
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		//fmt.Printf("::> [%s]...\n", path)
		if filepath.Base(path) != "requirements" {
			return nil
		} else {
			fnames = append(fnames, path)
			fmt.Printf("::> [%s]...\n", path)
		}
		return err
	})
	handle_err(err)

	reqfiles := make([]*ReqFile, len(fnames))
	for i, fname := range fnames {
		reqfiles[i], err = parse_file(fname)
		handle_err(err)
	}

	for _, req := range reqfiles {
		err = render_yaml(req)
		handle_err(err)
	}
}

// EOF
