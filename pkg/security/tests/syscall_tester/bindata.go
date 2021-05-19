// Code generated by go-bindata. DO NOT EDIT.
// sources:
// bin/syscall_x86_tester
// +build functionaltests,amd64

package syscall_tester

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataSyscallx86tester = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\x0d\x70\x54\xd7\x75\xfe\xde\xdb\x27\x78\x86\x65\x59\x84\x02\x32\xd0\x76\x8b\x17\x07\x27\xb0\xfa\x41\xc6\xb8\xc6\x53\x04\x92\x80\x56\xc2\x32\x46\xb8\xa9\x7f\x9e\x57\xbb\x4f\xda\xb5\xf7\xaf\xbb\x4f\x36\x72\xb1\x03\xde\x92\x78\xbb\xa1\x5d\x5c\xea\x9f\x71\xd3\x21\xa4\xc9\x78\xda\x78\xc6\x33\x25\x0e\x53\x7b\x8a\x5c\x11\x64\xb7\xee\x84\x38\xe9\x8f\xeb\xce\x94\x34\xa6\x5d\x26\xc2\x43\x6a\x8c\x65\x49\x70\x3b\xe7\xde\xfb\x76\xdf\xae\xa5\xa8\x6d\x3c\xe3\x99\x66\x8f\xe6\xee\x7b\xdf\xbd\xe7\x9c\x7b\xce\xb9\xf7\xed\x3b\x77\x75\xef\x17\x3b\xbb\xbb\x14\x45\x81\x4d\x2a\x5c\x20\x74\xb6\xa0\xe9\x6d\x00\x22\xb7\x8a\xfa\x36\xf8\xb0\x10\x6b\xf1\x4b\x58\x85\x79\x1c\x03\x6d\x07\x34\x9d\xca\x03\x0a\x40\x45\x83\x28\x2e\x00\x47\x15\xe0\xe8\x41\x4d\xa7\xb2\x14\xc0\x52\xd9\xa6\xc8\xc2\xe9\x80\xa6\x53\xf1\x6a\x00\x15\x6a\x87\x57\xb6\x7b\x01\x14\x34\x9d\x8a\x75\x1d\x40\xa5\xce\xd9\xee\x03\x70\x4c\xd3\xa9\x1c\x75\x01\x54\x9c\xf2\xee\x26\xc0\x7d\x4a\xd3\xa9\xac\x53\x80\x36\x05\xdc\x6e\x6a\x57\x01\x34\x34\x01\x0d\xa7\x34\x9d\x4a\x11\x00\x95\x79\x28\xfb\xf0\x82\x02\xbc\x70\x50\xd3\xa9\x74\x00\xe8\x70\xb4\xf5\x9e\xb7\xc2\xfb\x57\x03\xfb\x8f\x6b\x3a\x95\xcd\x00\x36\x3b\xda\xef\x3c\x6f\x85\x31\x03\xcd\x93\x6e\xed\x3e\x6f\x85\x9d\xf6\x5d\x06\x70\xd9\x11\x9f\xa6\x58\xb4\xbf\x29\x16\x5e\x1f\x8b\x26\x86\xf6\x05\x32\xc9\x40\xab\x68\x6b\x90\xb1\xdd\xbe\xab\x0f\xf7\x3e\xfa\x60\x28\x7b\x7a\xd1\xe7\x5b\xde\xec\x39\xbe\xe6\xda\xc6\xad\xb1\x7b\xde\x9d\xaf\x49\xfd\x8a\xe4\x81\xe4\x57\x4b\x63\x0b\x2c\x94\xed\x3c\x96\x3e\x75\x35\x64\xdd\xa6\xff\xca\xaf\x7a\xf1\x37\x7f\x3c\x52\x6d\xf3\x9d\x8e\xfb\x7a\x00\x9f\xad\xc2\xbe\x2a\xfc\xb8\x03\x53\xdb\x6f\x55\xb5\xb7\x57\xe1\x15\x55\x78\x67\x15\xbe\xb9\x0a\x6f\xae\xc2\x01\xb2\x7d\x44\xd3\xc9\xf7\x25\x68\xe4\xbe\x68\xc7\x6c\xbc\x18\x88\x45\xfb\x43\x14\xc3\x8d\x30\x76\xde\x61\x64\xac\x70\x34\x61\x0c\x65\xcc\x30\x06\x92\x29\x33\x81\x94\x99\x4e\x27\xd3\x18\x08\xc5\x92\x19\x13\x19\x2b\x6c\xa6\xd3\x18\x88\xc6\xcc\x44\x12\x41\x2b\x19\xc5\x40\x2a\x1d\x4d\x58\x03\xc8\x0c\x67\x42\xc1\x58\x0c\x19\x2b\x1d\x8a\xa7\x60\x18\xa4\xd8\xc8\x58\xc1\xb4\x65\xc4\x83\xd1\x04\xb6\x77\xef\xdc\xba\xcd\x68\x0d\xb4\x94\xee\x9a\x61\x18\x83\xf1\x64\x42\x72\x19\x90\xcf\x97\xca\xaf\xe2\x4e\xfc\x29\x50\xf8\xa7\x18\x3f\xa2\x25\xd1\xe8\x22\x1a\xbd\x07\x65\x9d\x97\x63\x15\x96\x6c\x9f\x3e\xa5\xe9\xf3\x34\xe1\x7b\xdd\x42\xc0\x3d\xa2\xe9\xf3\x15\xc0\x4b\x57\x15\x68\xa0\xab\x0b\x68\xa4\x6b\x1d\xb0\x8a\xae\xf3\x00\x1f\x5d\xe7\x03\x7e\xba\xea\xc0\x5a\xba\x5e\x07\xac\xa3\xeb\x82\x99\x66\x6d\x8d\x6a\x54\xa3\x1a\xd5\xa8\x46\x35\xaa\x51\x8d\x6a\x54\xa3\xff\x2f\xf4\xfe\xe2\x5f\x9e\xba\x2b\x3b\xae\x17\x97\x2a\xc0\xc1\xd1\xf7\x9b\x80\x7c\x76\x9a\x31\x76\x68\xc4\x52\xd9\xd9\xec\x69\xfd\x9e\x51\x27\x3f\xbb\x59\x1b\xd1\x74\xb6\x46\x1f\xd1\x74\x8e\xd7\xd0\xca\x23\x42\xb7\x17\xce\x31\xc6\xd8\x1a\x5a\x81\x44\xa8\xed\xc2\x59\x8e\x69\x25\x12\xa1\x25\xcc\x85\x11\x8e\x69\x45\x12\x69\x24\xfc\x12\xc7\xb4\x32\x89\xd0\x52\xf1\xc2\x31\x8e\x69\x85\x12\x59\x4b\xb8\xc0\x31\xad\x54\x22\xcd\x84\x0f\x70\x4c\x2b\x96\xc8\x26\xc2\x29\x8e\x69\xe5\x12\xd9\x42\xf8\x01\xc6\x18\xf9\xd3\x72\xf1\xfe\xdc\x8f\xb2\xef\x5e\xea\xdd\xb3\xbb\x78\x03\xc8\xad\x55\x4d\xc0\xe1\x6c\xef\xdb\x8c\xf5\x1e\xce\x9e\xfb\x67\xc6\x7a\xef\xdc\x3b\x36\x32\xf2\xb4\xa6\xf7\x16\x1f\x64\x8c\x5d\xce\xaf\xf2\x8f\x0e\x14\x06\x0a\x05\x92\xe7\x77\xe2\xaf\x60\x37\xd8\x7f\x27\x69\xc1\x75\x3b\x7d\x58\xfe\x93\xe4\xf5\xa1\x11\x6b\x65\x5f\xee\x7c\x76\xbc\x21\x42\xb5\x14\x32\xef\x1b\xa3\x87\x4f\xdc\x48\x8d\x03\x05\xfb\x4e\x5e\x84\xfc\x7a\xfa\xc8\xfd\xcd\x6b\xc5\x5f\x79\x6d\x52\x55\xce\x7e\x7f\xc2\xf2\xfd\x15\xd7\xf5\x96\xb5\x9c\xeb\xf2\xf6\x0a\x65\x6f\x09\x65\xd6\x8d\xb0\xf5\x90\x7d\x07\x6e\xdf\x4c\xd1\x1f\x12\xfd\xea\xc5\x08\x63\xec\x4c\x1d\xd5\x29\xa5\x9e\x9d\xfc\x3f\xf9\xfd\xb2\xfd\x5c\x64\x6d\xbe\xd3\x9d\xef\xd1\xb3\x8f\xe9\xda\xe2\xdf\xa3\x08\x7d\x8d\x1b\xa4\xf9\x0f\xd7\xe9\xc7\x34\x3d\xd7\xe1\xd7\x8a\xa1\x6b\x8c\x8d\x75\x4e\xd3\x9a\xf4\x42\x0f\x00\x12\xd9\x42\x3c\xc5\xdc\x35\xc6\x72\x9d\x13\xbc\xc2\x4d\x15\xfb\x78\xc5\x65\x5e\xa1\xe5\x7b\x26\xf2\x7d\x97\xc7\x34\xff\xcb\x00\x48\x53\xae\xdb\xaf\xe7\xf6\xf8\xdd\xc5\x16\xa1\x11\x5c\xd9\x74\xf6\xf4\xda\xfb\x46\x0b\x73\xda\x73\x87\x6d\xcf\x0f\xae\x96\xed\xf9\x86\x6d\x8f\xc6\x79\x9e\xb1\x79\x4e\x5d\xe5\x96\x65\x1f\x9b\xc0\xe2\x43\x2b\x28\xea\x75\xcf\x1d\xe3\x8a\x8a\x96\x43\x7c\x80\x8b\x4f\x50\xf5\x8b\x57\xcb\xa6\x73\xbe\x6f\xf1\x8a\x4b\x65\xe7\x9e\xe7\x15\xe3\xf9\xce\xcb\xf9\x9e\x4b\xf9\xbe\xf1\x31\xcd\x6f\x54\x7b\x76\xf7\x55\xc6\xf2\x3d\x13\x39\xb7\x3f\xd7\x59\x2c\x7e\x34\x3d\x83\x9f\x36\xf5\xe5\xce\xef\xcd\x8e\xb7\xcd\xee\xf0\xcb\xb6\x33\xff\xce\xd5\x4c\x70\x8b\xfb\x2b\x06\x40\x77\x18\xcd\x6d\xfc\x70\x9a\xb1\x96\x37\x6c\xb3\xb5\x7c\xdf\xe5\xfc\xd0\xa5\x31\xcd\xff\x39\x05\x18\xeb\xf0\xf3\xdf\xc5\x72\x1d\xc2\xd8\x9c\xe5\xf7\x8e\x75\xf8\x1b\xa0\x00\xb9\x9e\xf1\xe2\xf3\xa2\x1b\x61\xed\x44\xf6\x74\xdb\xfd\x4e\x73\xe7\x1c\x1f\x7c\x4d\x9a\xbb\x61\x7a\xb6\xf9\xf2\x85\xe9\xaa\xf9\xb2\x63\x7a\xa6\xf9\xe2\xad\x8e\xaa\x32\x53\x1c\xe7\xb2\x67\x87\x6d\xcf\xd7\xa7\x66\xb3\xe7\x8d\xa9\x2a\x7b\x4e\x4e\xcd\x64\xcf\x0f\xab\xed\xf9\x9d\xa9\xff\x83\x3d\x47\x6c\x7b\x96\x4e\xcd\x3d\x7f\xd9\xe4\xac\xf3\xf7\xc4\xe4\x8c\xf3\xf7\xe2\x64\xd5\xfc\x1d\x9f\xac\x9a\xbf\xef\x4c\x7e\x6c\xfe\x7e\xaf\xda\xb3\x67\x26\x9d\xf3\xf7\xd6\xc9\x9f\x3d\x7f\x7f\xa6\xbf\xef\xd8\xfe\xde\x30\x39\x5b\xfc\x77\x4e\x56\xc5\x7f\xd3\xe4\x4c\xf1\x3f\x53\x6d\xe5\x4f\x3f\x9a\x2d\xfe\x8d\xc2\x1e\xd9\x98\x7d\x4c\x57\x16\xff\x61\xa5\x59\x6b\x8e\x4b\xb3\xbe\xf2\x51\xd9\xac\xb7\x94\xd2\x30\x90\x39\x13\xb9\x1f\xe5\x94\xb1\x76\xed\xf6\xe3\x9a\x5e\xf4\x7f\xc4\x58\x96\x86\xc1\xc7\xbb\xd3\xb3\x45\x25\xdf\xe3\xce\xbe\x46\x23\x96\xeb\xf6\x6b\xc5\xbf\xe4\x46\x4f\x5f\x38\xac\x88\xc1\x10\xb2\x9b\x49\xb6\x38\x31\x87\xec\xfc\x29\x21\x7b\x67\x85\xec\x36\x92\xfd\xd6\x5c\xb2\xaf\x48\xd9\xcf\x54\xc8\x76\x93\x6c\x6a\x2e\xd9\xa6\x69\x21\xfb\x36\x9c\xb2\x7b\x48\xb6\x6d\x2e\xd9\xac\x94\xfd\x6a\x85\xec\x5d\x24\x3b\xf1\xe1\x1c\xb2\xff\x24\x65\x83\x15\xb2\xf7\x92\xec\xab\x73\xc9\xe6\xaf\x0a\xd9\x35\xf6\x78\xd2\x44\xa5\x21\x35\xe5\x90\xd2\xfc\x28\x7e\xfb\xc3\xd2\xb0\xf2\xb9\xd1\x78\x9f\x7c\x6b\xd3\xfb\xaf\xaf\xf8\x10\x65\x01\xdf\x5d\xf6\x59\xe0\xee\xbd\x77\x65\xc7\xdd\xb9\x9f\xe4\xf7\xfb\xd7\x16\xd9\x15\xc6\x0e\x7f\xd5\xcb\x18\x3b\x7c\xc8\xcd\x18\xbb\x69\xf4\xb5\x29\xd5\xba\xa9\xe5\x8a\xfd\xd2\x7e\x18\xd9\x71\xed\x6e\x66\xf9\xd7\xf1\x72\xf4\xdb\xc4\x95\x3d\xa3\x64\x4f\x7b\x6f\x7d\x7f\xe8\xc7\xd9\xd3\xee\x7b\xee\x37\xee\x2b\xbd\xe3\xad\x1b\xf9\xfb\x76\x34\xbf\xce\x3f\xea\xcc\xad\xbe\xf3\x21\x63\x07\x47\x9f\xb9\x11\xf8\x58\x3e\x55\xa3\x1a\xd5\xa8\x46\x35\xaa\xd1\xa7\x41\x2e\xfe\xff\x60\x15\xbd\x31\x33\x98\x31\x7d\xa9\x60\x26\xe3\x0b\xfa\x06\xa2\x31\xba\xb7\x22\xeb\x7c\x61\x33\x63\x45\x13\x41\x2b\x9a\x4c\xf8\x86\xa2\x61\x5f\x30\x11\xae\xa8\x1b\x8c\x86\x7d\x56\xd2\x17\x8a\x24\x1f\x49\x2c\xf8\xb9\xd5\x0c\x48\x3d\x69\x74\x05\xa3\x31\x93\xd7\x25\x53\x66\xc2\x97\x4a\x27\x1f\x8e\x86\xcd\x30\xd7\xf9\x09\x75\x13\xb4\x7e\x7e\x83\x63\x9f\x90\xe3\x5c\xcd\x86\xd6\x4f\x2a\x82\x9f\x84\xa6\xd8\x2c\x9a\x42\xc9\x78\x3c\x98\x08\x2f\x80\xe8\x0a\x76\x30\x11\x73\xe2\x0d\xad\xb0\x15\xa0\x2f\xf1\x50\x22\xf9\x48\xc2\x96\xf4\x3d\xb0\x26\xf3\xc0\x02\x40\x59\xe9\xba\x6d\x93\xdc\x17\x71\xe2\x22\x63\xfb\x01\xec\x79\x8f\xb1\x3d\x00\x8e\xbe\xc7\x58\x04\x40\xf8\x0a\x63\xc7\x00\xfc\xf0\x0a\x63\xe3\x00\xde\xbe\xc2\x18\x2d\x5d\x1b\xe4\xfc\x55\x1e\xdd\x0d\x65\xbf\xae\xac\x74\x6b\xda\x93\x8a\xf8\x5f\xf9\x2a\x00\x13\x17\x19\xbb\x8d\x18\x3a\xe6\xeb\x54\xd7\x0c\xc0\xff\x1e\x63\x75\x52\xce\x27\xf7\x74\x34\x5f\x14\xfa\xe1\xd1\xbb\x3c\xee\xdf\x58\xbc\xd0\xd2\xf6\xe1\xd7\x57\xdc\xf6\xb9\x56\xff\xea\x1d\x00\xc8\x86\x57\x3f\x60\xcc\x24\x9e\x4e\x8f\x7e\x48\xed\xf6\xb8\xbf\xec\x6a\xf7\x78\xbf\xa4\xb5\x7b\x1a\xb2\x75\xdb\x3c\xbe\xb8\xc7\xdf\xee\x59\xdb\xe1\x59\xd7\xe1\x69\xee\xf1\xf8\xb6\x7b\x1a\xda\x47\x3d\xde\xf6\x33\x1e\x77\xfb\x98\x47\x6f\xff\xae\x47\x13\x76\x9d\x00\x70\xee\x83\xb2\x0d\x54\xf7\x3a\x80\xef\x7f\xc0\x98\xf6\x69\x3c\xfc\x35\xaa\x51\x8d\x6a\x54\xa3\x1a\xd5\xa8\x46\x35\xaa\xd1\x2f\x30\x9d\x3b\xa2\xe9\x2f\x1d\xd1\x74\x7b\x8f\xba\x1b\x62\x0f\xfa\x22\x00\xf7\x3e\xa7\xe9\xd7\x43\xec\xd9\x5e\x29\xf7\x6a\xaf\xa0\x1c\xfe\x94\xa6\xaf\x92\xf8\x83\x6b\x2c\x39\x7e\x50\xd3\x29\xbf\x3f\xf7\x84\xa6\xf3\xfd\xde\x4f\x68\xfa\x02\x00\x87\xe5\x3e\x6b\xca\xf9\x3f\x23\xfb\xa3\xb5\x37\x46\x34\x5d\x05\xb0\x43\xee\xef\x5e\x02\x60\x39\xad\x0b\xb2\x9a\x4e\xf7\x27\xb3\x9a\x5e\x2f\xe5\x96\x02\xd0\x01\x5c\x63\x2c\xf9\x64\x56\xd3\x19\x63\x49\xb2\xf3\x12\x63\xc9\x44\x56\xfc\x3f\xfe\x7f\x4b\x0d\xa7\xca\x72\x5d\x05\x4d\xdf\x5b\xd0\xf4\x81\x82\xa6\x3f\x5c\xd0\xf4\x2f\x15\x34\xfd\xe9\x82\xa6\x7f\xb3\xa0\xe9\x2f\x17\x34\xfd\x4c\xa1\xcc\xbb\x7d\xdb\xb6\x5f\xf3\xad\xed\xeb\x1f\x4a\x58\x43\xbe\x5b\x03\x1b\x02\xcd\xeb\x5b\x6e\x19\xe2\xb0\xe5\xf1\xd6\xe6\x40\x73\xdb\x4d\xa2\x1a\x92\x27\x14\x0b\x26\x06\x7d\x0f\x9b\xe9\x0c\x2d\x31\x5b\x5a\x02\xcd\x81\xe6\xf5\xad\x8f\x0b\x11\x2e\x10\x68\x99\xc5\xc6\xa3\x07\x45\xbf\x2e\x79\x72\xe0\x85\x12\x16\x3b\xdb\xdf\x2c\x61\x1e\x51\x8c\x97\xb0\x58\x55\xd1\x18\x08\x2c\x56\x5e\xe7\x4a\x98\x9f\x08\x80\x1d\x3b\x17\xe6\x73\xfc\x64\x09\x0b\x77\x4f\x96\xf0\x75\x1c\xbf\x5e\xc2\x72\xc3\x74\xc1\xc6\x0b\x39\x6c\x2e\x61\x3e\x83\xf8\x59\x0a\x81\x17\x71\x4c\x73\x49\x60\x8f\x90\x3f\x66\xe3\xc5\x1c\xee\x3f\x6e\x63\xb1\x1b\xfc\x64\x09\x2f\xe1\xd8\x7d\xca\xc6\xf5\xc2\xbf\x12\x5e\x5a\x31\xa6\x2e\xb9\x4a\x9e\x2e\x61\x39\xf3\x46\x6c\xbc\x4c\xd8\x5b\xc2\xcb\x39\xde\x54\xc2\x8d\x15\xe3\xe0\xc2\xf5\xe5\xb3\x1b\x7c\xce\xff\x94\x91\x87\xcd\x47\x04\xbf\x8a\x45\xdc\xa3\x94\x03\xff\x2a\x80\x97\x1c\xf8\x16\x00\x7c\xa7\x02\x7f\xc6\x1a\xd1\xe5\xb0\x5f\xc1\x52\xc4\xe5\x73\x68\xf3\x0f\x3b\xfc\x55\x50\x8f\x67\xab\xfa\xaf\xb6\xe7\x45\x1a\x8e\xaf\xdb\xfc\x4b\x2a\xec\xa7\xf6\x57\x1c\xfd\x01\xf5\x78\xd3\x11\x2f\x05\x0d\xf8\x07\x47\x7f\xd4\x5e\x74\x8c\x07\x3d\x85\x53\x8e\xf8\x29\x58\x86\x25\x0a\xd0\xfb\x9c\x78\xe6\xeb\xc9\x5f\x87\x41\x34\x3a\x1b\x14\x60\xcb\x1f\x6b\xfa\x7f\xc8\xf6\x0e\x05\xf0\x1d\x11\xe7\x13\xea\xd5\x45\xf8\x9e\x52\x8e\xbf\x0f\xcb\xf1\xdb\x4a\xe5\x79\x86\x98\x02\xf4\xc9\xf9\x42\xfc\x4f\x28\xe5\xf1\xf1\x62\x39\x72\x55\xfd\xad\x54\xca\xf3\xab\x5e\xf5\xe0\x79\x05\x38\xf0\x54\xb9\xff\x3f\xa7\xfe\x9e\xd6\xf4\x3f\x90\xf8\x95\x2a\xf9\xbf\x77\xd8\x43\xfa\xff\xc5\xd1\xce\x37\xe9\x28\x40\x9b\x6c\x5f\xa2\x2e\xc7\x15\x05\x48\xfd\x91\xa6\x7f\x43\xea\x9b\xaf\x56\x9e\xbf\x58\xa6\x56\xc5\xa3\x0a\x77\xa9\xc0\xb9\x67\x35\xdd\x94\xf2\x7b\xab\xda\xc7\x00\x6c\x29\xd9\xd3\x88\x41\x15\xf0\x1e\x91\xf1\x56\x17\xe1\x8b\x55\xfc\x3f\x50\xc4\xf3\x76\x9b\xd4\xf7\x54\x55\xfb\x9f\xa9\xe2\x79\xb3\xed\xfb\x0b\xd5\x19\xcf\x46\x7c\x47\x05\x46\x9e\xd6\xf4\x15\xaa\x90\x7f\x55\x05\x7c\x47\xcb\xf1\x6a\x53\x80\x97\x1c\xf8\xef\xaa\xf4\xbf\xe3\xd0\x47\xf1\xf9\x13\x05\xb8\x74\xa4\xcc\xdf\xab\x96\xbf\x2f\xea\xd5\x85\x40\x28\x6d\x65\xac\xa1\x81\x81\x40\x08\x61\x33\x6d\x0e\x46\x33\x96\x99\x36\xac\xb8\x11\x8a\x25\x13\x66\x06\x86\x11\x4e\x1a\x83\xb1\x64\x7f\x30\x66\x84\xad\x64\x3a\x63\x04\x87\xf6\x21\x94\x8c\xa7\x62\xa6\x65\x86\x03\xb7\x6c\x6c\x6d\x9d\x99\xc9\x18\x88\x26\xa2\x46\x30\x9d\x0e\x0e\x1b\x66\xc2\x4a\x0f\x63\x20\x1d\x8c\x9b\x46\x78\x28\x1e\x1f\x86\x61\x38\x90\x11\x4d\x44\xad\x0a\x56\x79\xdc\xc5\xd8\xb7\x69\xa3\x61\x99\x64\x53\x20\x04\xc3\xe8\xda\xdd\xde\xd3\x69\x74\xee\xea\x30\x0c\x18\x95\x52\x61\x18\x1d\x5f\xd8\xd5\xde\xb3\x73\x5b\x65\x0b\x3f\xfc\x02\xc3\xd8\xbe\xab\xcf\xe8\xdc\x21\x35\xec\xe8\xd8\x0d\x63\x7b\xf7\x1d\x5b\xdb\xbb\x8d\x3b\xba\xba\xee\xea\xdc\x63\xec\x69\xdf\xda\xdd\x69\xd8\x07\x6b\x42\x99\x21\x6e\xbf\x3c\x6f\xb3\x65\x4b\xf9\x50\x8d\xfd\xfb\xa2\x61\x1f\xc9\x31\xc8\xca\xc0\xa0\x69\x19\xa9\x90\x61\x45\x86\x12\x0f\x05\xfa\xf7\xc9\xb3\x3c\x4e\xc1\x99\xf8\x52\x30\xcc\x70\xd0\x0a\xca\x13\x40\x65\xf6\x16\xd9\x4f\xa9\x97\x58\x75\xb7\xe2\xe8\x50\x65\x07\xa4\xaa\xe4\x71\xc5\xd9\x1f\xc3\x08\x67\x92\x46\x24\x98\x08\xc7\xcc\xd2\x4f\xa6\x65\x17\x2a\xcf\x26\x7d\xec\x74\x51\x85\xff\xe2\x44\x52\x65\xc7\xa5\xa0\x51\xe4\xc5\xd9\x26\xa7\x2f\x46\x38\x66\xa4\xcd\x58\x32\x14\xb4\x4c\x52\x6b\x45\x43\x46\x2a\x6a\xda\xc3\x5c\xa1\x9e\x1f\x7e\xaa\xd0\x3e\x90\x32\x22\x8f\xc0\x30\xfa\x33\x19\xe9\x1c\x3f\xf1\x14\xab\x0c\x50\xd0\x4a\x46\x2b\x8d\xda\xd3\xb3\x4d\xce\x15\x04\x32\xc3\x71\x2b\xd8\x8f\x40\xc6\x4a\x8b\x6b\xc4\xbe\x8b\x26\x2c\x33\x9d\x42\x20\x91\xb4\xcc\xc0\x60\x62\x28\xd0\x3f\x14\x8d\x85\xd7\x47\xc3\xb2\xaa\x7d\xeb\xce\xf5\x56\x70\x10\xbc\x2d\x12\xcc\x44\x10\x08\x0f\x27\x32\xc3\x71\x71\xb5\xd2\xa2\x45\xa6\x16\x15\xc0\x48\x23\x90\x36\x63\xc4\x27\x6e\x52\x31\x8b\x3a\x8c\x5a\x08\x58\xe6\x3e\x0b\x01\x3e\xc7\x02\xe9\x24\x9f\x03\x01\x33\x22\x9f\x8a\x48\x38\x5d\x46\x42\x42\x4c\x67\x21\x61\xdf\x87\x87\x13\xc1\x78\x34\x84\xc0\x60\xd2\xe2\x1f\xa2\x03\xa1\xac\x3f\x93\x41\x20\x94\x8c\xc7\xcd\x84\x85\xff\x29\xad\x94\xb9\xa7\x2a\xf3\x9e\xa3\x0a\xe4\x5b\x5d\x90\xfd\x75\x7c\x03\xc0\xb3\x15\x55\xe6\x43\x2f\x28\x80\xdf\xc1\x67\xff\x9e\xbc\xd1\xc1\x47\x79\xd2\x9b\x4a\xe5\x99\x3d\x9b\xaf\x03\xc0\x15\xc6\x92\xaa\xcc\x9f\xc6\xa5\xbe\x3a\x07\x1f\x95\x5d\x32\x8f\x55\x65\x5e\xe5\x55\x81\xb3\xf2\x37\x7c\x45\xf2\x50\xd6\xb2\xd7\x71\xf6\x90\xf2\xad\x73\xaa\xc8\x81\xab\xfd\xb8\x1f\x00\x93\xfd\x52\x1e\x96\x70\x89\xbc\xda\xee\x57\x95\xe5\x21\x99\xf7\xaa\x32\x3f\x7b\xd2\x25\x7e\xd3\x77\xf6\x4b\xf4\x28\xc0\xb3\x34\x55\xe6\x6d\x27\x5d\xc2\x1e\xa7\x1f\xf4\xe5\x9b\x95\x7c\x5b\x65\x3e\xf7\xba\x4b\xe4\xe0\xc4\xb7\xcc\xc1\xf7\x15\xa9\x9f\x67\x8a\xf4\xbd\xed\x9d\x39\xce\x5f\x76\xf0\x51\xfe\xd7\xec\x05\x8e\x39\xf8\xbc\x92\xf7\x29\x07\x1f\xbd\xa7\xce\x7a\x81\x9c\x5e\xc9\x47\xf4\xac\x83\x8f\xde\xe7\xf7\x5e\x8f\x8a\x6c\xcc\xee\xf7\x4f\x1d\xf3\x85\xe7\x91\x3e\xe0\x77\xd5\x8f\xf3\x7d\xd3\xc1\xc7\xcf\xad\xae\xae\x3c\x4b\x69\xf3\x9d\x70\xf0\x51\xde\x79\x72\x35\xf0\xaf\x33\xf0\xfd\x35\xc0\x73\x3d\x97\xcc\x97\xdc\x4d\xe5\x36\xe7\x7c\xf9\x5b\x80\x67\xb5\x2e\x99\x77\x79\x67\xe1\xfb\x47\xe9\xab\x4b\xe6\x63\x0d\x4d\xe5\x33\xb9\xce\x71\xfb\x37\x69\x9f\x4b\xe6\xb5\xd3\xb3\xe8\x7b\xd7\xc1\xc7\xf3\xb5\x66\x31\x5f\xaa\xf9\x2e\x3a\xf8\x28\xef\x69\x6e\x06\xf4\x19\xfc\x7d\x5f\xf6\xef\x92\xf9\xf1\xa6\x2a\x3e\xfb\x7e\x42\xea\xb3\xfb\x22\xbe\xbd\x0e\x3e\xc5\x51\x1c\xc3\x84\x42\x33\x70\x49\x13\xcf\xff\xe7\x1d\xcf\xd1\x75\xb6\x0f\x92\x0e\xdc\x0c\xfc\xa7\x5a\xa9\x0f\x72\xdd\xe8\xe4\x1b\xdc\x04\x38\x8e\x75\x97\xf8\xfe\x3b\x00\x00\xff\xff\xb5\x30\xb9\x39\xf0\x3d\x00\x00")

func bindataSyscallx86testerBytes() ([]byte, error) {
	return bindataRead(
		_bindataSyscallx86tester,
		"/syscall_x86_tester",
	)
}

func bindataSyscallx86tester() (*asset, error) {
	bytes, err := bindataSyscallx86testerBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name:        "/syscall_x86_tester",
		size:        15856,
		md5checksum: "",
		mode:        os.FileMode(509),
		modTime:     time.Unix(1, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"/syscall_x86_tester": bindataSyscallx86tester,
}

//
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op:   "open",
					Path: name,
					Err:  os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"": {Func: nil, Children: map[string]*bintree{
		"syscall_x86_tester": {Func: bindataSyscallx86tester, Children: map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
