package downloadserver

import (
	"io"
	"os"
)

type DownloadServer interface {
	GetFileList(path string) ([]os.FileInfo, error)
	DownloadFile(remotePath string, writer io.Writer) error
}
