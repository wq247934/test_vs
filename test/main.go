package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir := "/root/work/release/etc/certs/"

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

					if subFile.Name() == "config.yaml" {
						leafBytes, leaf := parseLeaf(filepath.Join(subDir, "fullchain.pem"))
						os.WriteFile(filepath.Join(dnsDir, "leaf.pem"), leafBytes, 0755)
						conf := viper.New()
						conf.SetConfigFile(filepath.Join(dnsDir, "config.yaml"))
						conf.Set("domain", leaf.Subject.CommonName)
						conf.Set("algo", "algo")
						conf.Set("method", "dns01")
						conf.WriteConfig()
					} else {
						err := copyFile(oldPath, newPath)
						if err != nil {
							log.Fatal(err)
						}
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

func parseLeaf(path string) ([]byte, *x509.Certificate) {
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
			if !cert.IsCA {
				return pem.EncodeToMemory(
					&pem.Block{Type: "CERTIFICATE", Bytes: block.Bytes}), cert
			}
		}

		data = rest
	}
	return nil, nil
}
