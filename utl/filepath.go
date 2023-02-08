package utl

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 拷贝文件
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, os.ModePerm)
}

// 存在
func Exist(name string) bool {
	matches, err := filepath.Glob(name)
	if err != nil {
		return false
	}

	return len(matches) > 0
}

// 制作目录
func MakeDir(names ...string) error {
	for _, name := range names {
		if !Exist(name) {
			err := os.MkdirAll(name, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 删除
func Remove(names ...string) error {
	for _, name := range names {
		matches, err := filepath.Glob(name)
		if err != nil {
			return err
		}
		for _, m := range matches {
			err = os.Remove(m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 解压缩文件
func UnzipFile(src string, dst ...string) error {
	var dstDir string
	if len(dst) > 0 {
		dstDir = dst[0]
	}
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, f := range reader.File {
		filePath := filepath.Join(dstDir, f.Name)
		if f.FileInfo().IsDir() {
			err = MakeDir(filePath)
			if err != nil {
				return err
			}
			continue
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()
		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// 压缩文件
func ZipFile(src string, dst ...string) error {
	var dstDir string
	if len(dst) > 0 {
		dstDir = dst[0]
	}
	filePath := filepath.Join(dstDir, fmt.Sprintf("%s.zip", src))
	dstFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	writer := zip.NewWriter(dstFile)
	defer writer.Close()
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileHeader.Name = path
		if !info.IsDir() {
			fileHeader.Method = zip.Deflate
		}
		dstFile, err := writer.CreateHeader(fileHeader)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
