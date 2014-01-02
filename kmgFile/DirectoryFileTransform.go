package kmgFile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryFileTransform struct {
	InputExt     string                 //input file ext (case sensitive)
	OuputExt     string                 //output file ext
	IsIgnoreFile func(path string) bool //default IsDotFile
	Transform    func(r io.Reader, w io.Writer) error
}

func (transform *DirectoryFileTransform) Run(inputPath string, outputPath string) error {
	if transform.IsIgnoreFile == nil {
		transform.IsIgnoreFile = IsDotFile
	}
	fi, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return transform.transformOneFile(inputPath, outputPath)
	}
	//check output path is dir
	fi, err = os.Stat(outputPath)
	if err == nil {
		if !fi.IsDir() {
			return errors.New("input file path is a directory,but output file path is not a directory.")
		}
	} else {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		isIgnore := transform.IsIgnoreFile(path)
		if info.IsDir() {
			if isIgnore {
				return filepath.SkipDir
			}
			return nil
		}
		if isIgnore {
			return nil
		}
		if filepath.Ext(path) != "."+transform.InputExt {
			return nil
		}
		relPath, err := filepath.Rel(inputPath, path)
		if err != nil {
			return err
		}
		oFilePath := filepath.Join(outputPath, relPath)
		oFilePath = strings.TrimSuffix(oFilePath, "."+transform.InputExt)
		oFilePath = oFilePath + "." + transform.OuputExt
		err = transform.transformOneFile(path, oFilePath)
		if err != nil {
			return fmt.Errorf("[%s] %s", relPath, err.Error())
		}
		return nil
	})
}

func (transform *DirectoryFileTransform) transformOneFile(inputPath string, outputPath string) error {
	iFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer iFile.Close()
	oFileDir := filepath.Dir(outputPath)
	err = os.MkdirAll(oFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	oFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer oFile.Close()
	//in case it will remain some content
	err = oFile.Truncate(0)
	if err != nil {
		return err
	}
	return transform.Transform(iFile, oFile)
}
