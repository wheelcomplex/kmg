package kmgTar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

//package a directory into a tar format.
func PackageDirectoryTar(root string, inW io.Writer) (err error) {
	tarW := tar.NewWriter(inW)
	defer tarW.Close()
	err = filepath.Walk(root, func(path string, info os.FileInfo, inErr error) (err error) {
		if inErr != nil {
			return inErr
		}
		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return
		}
		h.Name = filepath.ToSlash(relPath)
		err = tarW.WriteHeader(h)
		if err != nil {
			return
		}
		if info.IsDir() {
			return
		}
		f, err := os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		_, err = io.Copy(tarW, f)
		if err != nil {
			return
		}
		return
	})
	return
}

//package a directory into a tar.gz format.
func PackageDirectoryTarGz(root string, w io.Writer) (err error) {
	gzipW := gzip.NewWriter(w)
	defer gzipW.Close()
	err = PackageDirectoryTar(root, gzipW)
	return
}
