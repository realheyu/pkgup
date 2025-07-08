package cmd

import (
	"fmt"
	"log"

	"github.com/realheyu/pkgup/internal/uploader"
	"github.com/spf13/cobra"
)

var (
	distPath string
	zipName  string
	filePath string
	version  string
	fileName string
	desc     string
	username string
	password string
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "压缩并上传目录到阿里云效制品仓库",
	Run: func(cmd *cobra.Command, args []string) {
		zipFile, err := uploader.ZipDist(distPath, zipName)
		if err != nil {
			log.Fatalf("压缩失败: %v", err)
		}

		err = uploader.UploadToAliyun(zipFile, filePath, version, fileName, desc, username, password)
		if err != nil {
			log.Fatalf("上传失败: %v", err)
		}

		fmt.Println("上传成功")
	},
}

func init() {
	uploadCmd.Flags().StringVar(&distPath, "dist", "", "待压缩目录路径")
	uploadCmd.Flags().StringVar(&zipName, "zip-name", "", "生成的 zip 文件名")
	uploadCmd.Flags().StringVar(&filePath, "file-path", "", "阿里云制品仓路径（如 a/b/c）")
	uploadCmd.Flags().StringVar(&version, "version", "", "版本号")
	uploadCmd.Flags().StringVar(&fileName, "file-name", "", "制品名称")
	uploadCmd.Flags().StringVar(&desc, "desc", "", "版本描述")
	uploadCmd.Flags().StringVar(&username, "username", "", "用户名")
	uploadCmd.Flags().StringVar(&password, "password", "", "密码")

	err := uploadCmd.MarkFlagRequired("dist")
	if err != nil {
		log.Fatalf("标记 dist 为必填失败: %v", err)
	}
	err = uploadCmd.MarkFlagRequired("zip-name")
	if err != nil {
		log.Fatalf("标记 zip-name 为必填失败: %v", err)
	}
	err = uploadCmd.MarkFlagRequired("file-path")
	if err != nil {
		log.Fatalf("标记 file-path 为必填失败: %v", err)
	}
	err = uploadCmd.MarkFlagRequired("version")
	if err != nil {
		log.Fatalf("标记 version 为必填失败: %v", err)
	}
	err = uploadCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("标记 username 为必填失败: %v", err)
	}
	err = uploadCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("标记 password 为必填失败: %v", err)
	}

	rootCmd.AddCommand(uploadCmd)
}
