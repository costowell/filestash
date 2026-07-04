package plg_backend_s3_csh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/mickael-kerjean/filestash/server/common"
	. "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_s3"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type S3CSHBackend struct {
	s3     IBackend
	params map[string]string
}

func init() {
	Backend.Register("s3csh", S3CSHBackend{})
}

// Cat implements [common.IBackend].
func (this S3CSHBackend) Cat(path string) (io.ReadCloser, error) {
	return this.s3.Cat(path)
}

// Init implements [common.IBackend].
func (this S3CSHBackend) Init(params map[string]string, app *App) (IBackend, error) {
	backend, err := S3Backend{}.Init(params, app)
	if err != nil {
		return nil, err
	}
	return S3CSHBackend{s3: backend, params: params}, nil
}

// LoginForm implements [common.IBackend].
// Homie gets called before init
func (this S3CSHBackend) LoginForm() Form {
	form := S3Backend{}.LoginForm()
	form.Elmnts[0].Value = "s3csh"
	return form
}

// Ls implements [common.IBackend].
func (this S3CSHBackend) Ls(path string) ([]os.FileInfo, error) {
	if path == "/" || path == "" {
		buckets, err := this.listBuckets()
		if err != nil {
			Log.Error("plg_backend_s3_csh::ls admin_bucket err=%s", err.Error())
			return nil, err
		}
		files := make([]os.FileInfo, 0, len(buckets))
		for _, name := range buckets {
			files = append(files, File{
				FName: name,
				FType: "directory",
				FTime: 0,
			})
		}
		return files, nil
	}
	return this.s3.Ls(path)
}

// listBuckets gets every bucket via the admin endpoint
func (this S3CSHBackend) listBuckets() ([]string, error) {
	region := this.params["region"]
	if region == "" {
		region = "us-east-1"
		if strings.HasSuffix(this.params["endpoint"], ".cloudflarestorage.com") {
			region = "auto"
		}
	}
	endpoint := strings.TrimSuffix(this.params["endpoint"], "/")
	req, err := http.NewRequest("GET", endpoint+"/admin/bucket?format=json", nil)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewStaticCredentials(
		this.params["access_key_id"],
		this.params["secret_access_key"],
		this.params["session_token"],
	)
	if _, err := v4.NewSigner(creds).Sign(req, bytes.NewReader([]byte{}), "s3", region, time.Now()); err != nil {
		return nil, err
	}

	res, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, NewError(fmt.Sprintf("radosgw admin bucket list failed status=%d body=%s", res.StatusCode, string(body)), res.StatusCode)
	}

	var buckets []string
	if err := json.Unmarshal(body, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
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
