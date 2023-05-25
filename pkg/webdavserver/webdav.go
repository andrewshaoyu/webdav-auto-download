package webdavserver

import (
	"io"
	"os"

	"github.com/studio-b12/gowebdav"
)

type WebDAVConfig struct {
	URL            string
	Username       string
	Password       string
	MaxConcurrency int
}

type WebDAVServer struct {
	semaphoreChan chan struct{}
	client        *gowebdav.Client
}

func (s *WebDAVServer) GetFileList(path string) ([]os.FileInfo, error) {
	return s.client.ReadDir(path)
}

func (s *WebDAVServer) DownloadFile(remotePath string, writer io.Writer) error {
	s.semaphoreChan <- struct{}{}
	defer func() {
		<-s.semaphoreChan
	}()

	reader, err := s.client.ReadStream(remotePath)
	defer reader.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}
	return nil
}
