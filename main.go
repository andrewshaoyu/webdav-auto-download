package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/andrewshaoyu/webdav-manager/pkg/webdavserver"
)

var DownloadDir string
var DavServer string
var User string
var Password string
var webdav *webdavserver.WebDAVServer

func main() {
	flag.StringVar(&DownloadDir, "download-dir", "", "where file put")
	flag.StringVar(&DavServer, "server", "", "webdav server")
	flag.StringVar(&User, "user", "", "webdav user name")
	flag.StringVar(&Password, "password", "", "webdav password")
	flag.Parse()

	if DownloadDir == "" {
		DownloadDir = os.Getenv("DOWNLOAD_DIR")
	}
	if DavServer == "" {
		DavServer = os.Getenv("SERVER")
	}
	if User == "" {
		User = os.Getenv("USER")
	}
	if Password == "" {
		Password = os.Getenv("PASSWORD")
	}

	ticker := time.NewTicker(time.Hour)
	factory := webdavserver.ServerFactory{}
	webdav = factory.CreateWebDAVServer(webdavserver.WebDAVConfig{
		URL:            DavServer,
		Username:       User,
		Password:       Password,
		MaxConcurrency: 1,
	})
	semaphoreChan := make(chan struct{}, 1)
	downloadfiles()
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			semaphoreChan <- struct{}{}
			if time.Since(start) > time.Minute*10 {
				<-semaphoreChan
				continue
			}
			downloadfiles()
			<-semaphoreChan
		}
	}
}

func downloadfiles() {
	files, err := webdav.GetFileList("/")
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Println(file.Name())
		fileexists := true

		existsFile, err := os.Stat(DownloadDir + "/" + file.Name())
		if err != nil {
			if !errors.Is(err, fs.ErrInvalid) && !errors.Is(err, fs.ErrNotExist) {
				fmt.Println(err)
				continue
			}
			fileexists = false
		}

		// do not need to download
		if fileexists && existsFile.Size() == file.Size() {
			continue
		}

		// remove first then download
		if fileexists {
			err = os.Remove(DownloadDir + "/" + file.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		dlFile, err := os.Create(DownloadDir + "/" + file.Name())
		if err != nil {
			return
		}
		wg.Add(1)
		go func(file os.FileInfo, writer io.Writer) {
			err := webdav.DownloadFile("/"+file.Name(), dlFile)
			if err != nil {
				fmt.Println(err)
			}
			dlFile.Close()
			wg.Done()
		}(file, dlFile)
	}
	wg.Wait()
}
