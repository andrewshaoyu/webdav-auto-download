package webdavserver

import "github.com/studio-b12/gowebdav"

type ServerFactory struct {
}

func (f *ServerFactory) CreateWebDAVServer(config WebDAVConfig) *WebDAVServer {
	cli := gowebdav.NewClient(config.URL, config.Username, config.Password)
	return &WebDAVServer{
		client:        cli,
		semaphoreChan: make(chan struct{}, config.MaxConcurrency),
	}
}
