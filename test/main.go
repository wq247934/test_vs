package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir := "/Users/wangqian/Documents/go-project/test_vs/test/certs"

	// 读取目录内容
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历目录内容
	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			fmt.Println("Subdirectory:", subDir)
			dnsDir := filepath.Join(subDir, "dns01")
			err := os.MkdirAll(dnsDir, 0755)
			subFiles, err := os.ReadDir(subDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, subFile := range subFiles {
				if !subFile.IsDir() {
					oldPath := filepath.Join(subDir, subFile.Name())
					newPath := filepath.Join(dnsDir, subFile.Name())
					fmt.Printf("Moving %s to %s\n", oldPath, newPath)
					if subFile.Name() == "fullchain.pem" {
						leafBytes := parseLeaf(filepath.Join(subDir, subFile.Name()))
						os.WriteFile(filepath.Join(dnsDir, "leaf.pem"), leafBytes, 0755)
					}
					err := copyFile(oldPath, newPath)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

func parseLeaf(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// 解码 PEM 块
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}

		// 判断是否为证书
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}
			if len(cert.DNSNames) != 0 {
				return pem.EncodeToMemory(
					&pem.Block{Type: "CERTIFICATE", Bytes: block.Bytes})
			}
		}

		data = rest
	}
	return nil
}
