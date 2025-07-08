package uploader

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ZipDist(srcDir, zipName string) (string, error) {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	return zipName, err
}

func UploadToAliyun(zipFilePath, filePath, version, fileName, versionDesc, username, password string) error {
	url := fmt.Sprintf(
		"https://packages.aliyun.com/api/protocol/63800a90126bcb821717cd3e/generic/iot-saas-web/files/%s?version=%s",
		filePath, version,
	)
	if fileName != "" {
		url += fmt.Sprintf("&fileName=%s", fileName)
	}
	if versionDesc != "" {
		url += fmt.Sprintf("&versionDescription=%s", versionDesc)
	}

	file, err := os.Open(zipFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	pr, pw := io.Pipe()
	formWriter := multipart.NewWriter(pw)

	go func() {
		part, _ := formWriter.CreateFormFile("file", filepath.Base(zipFilePath))
		io.Copy(part, file)
		formWriter.Close()
		pw.Close()
	}()

	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", formWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上传失败：状态码 %d，响应: %s", resp.StatusCode, string(body))
	}

	return nil
}
