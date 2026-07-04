package plg_backend_s3_csh

import (
	"io"
	"os"

	. "github.com/mickael-kerjean/filestash/server/common"
	. "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_s3"
)

type S3CSHBackend struct {
	s3 S3Backend
}

func init() {
	Backend.Register("s3_csh", S3CSHBackend{})
}

// Cat implements [common.IBackend].
func (this S3CSHBackend) Cat(path string) (io.ReadCloser, error) {
	return this.s3.Cat(path)
}

// Init implements [common.IBackend].
func (this S3CSHBackend) Init(params map[string]string, app *App) (IBackend, error) {
	return this.s3.Init(params, app)
}

// LoginForm implements [common.IBackend].
func (this S3CSHBackend) LoginForm() Form {
	return this.s3.LoginForm()
}

// Ls implements [common.IBackend].
func (this S3CSHBackend) Ls(path string) ([]os.FileInfo, error) {
	if path == "/" || path == "" {
		return []os.FileInfo{
			File{
				FName: "pubsite",
				FType: "directory",
				FTime: 0,
			},
		}
	}
	return this.s3.Ls(path)
}

// Mkdir implements [common.IBackend].
func (this S3CSHBackend) Mkdir(path string) error {
	return this.s3.Mkdir(path)
}

// Mv implements [common.IBackend].
func (this S3CSHBackend) Mv(from string, to string) error {
	return this.s3.Mv(from, to)
}

// Rm implements [common.IBackend].
func (this S3CSHBackend) Rm(path string) error {
	return this.s3.Rm(path)
}

// Save implements [common.IBackend].
func (this S3CSHBackend) Save(path string, file io.Reader) error {
	return this.s3.Save(path, file)
}

// Stat implements [common.IBackend].
func (this S3CSHBackend) Stat(path string) (os.FileInfo, error) {
	return this.s3.Stat(path)
}

// Touch implements [common.IBackend].
func (this S3CSHBackend) Touch(path string) error {
	return this.s3.Touch(path)
}
