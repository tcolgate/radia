package graphalg

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _assets_index_html_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x53\xcd\x6e\xdb\x3c\x10\x3c\xc7\x4f\xb1\x1f\x63\xc0\x32\xbe\x9a\x8c\xeb\x1a\x28\x1c\xd9\x87\xfe\x1e\x5b\x24\x2e\x8a\x1e\x29\x6a\x2d\xd1\xa1\x48\x55\x5c\x46\x31\x8a\xbc\x7b\x29\xd9\x49\xe4\x1e\x7a\xa8\x2e\x24\xc8\x99\xd9\xd9\x1d\x2a\xfd\xef\xc3\x97\xf7\xdb\x1f\x5f\x3f\x42\x49\x95\xd9\x8c\xd2\xe3\x72\x91\x96\x28\xf3\xb8\x5e\xa4\x15\x92\x04\x55\xca\xc6\x23\xad\xd9\xb7\xed\xa7\xd9\x5b\x06\xa2\xbf\x22\x4d\x06\x37\xb7\xb2\xaa\x0d\x82\xdb\x41\x8b\x99\x77\xea\x0e\x09\x5a\x4d\x25\x14\xce\x48\x5b\xa4\xe2\x08\xeb\x08\x5e\x35\xba\x26\xf0\x8d\x5a\xb3\x92\xa8\xf6\x2b\x21\x54\x6e\xf7\x9e\x2b\xe3\x42\xbe\x33\xb2\x41\xae\x5c\x25\xe4\x5e\x3e\x08\xa3\x33\x2f\xf6\x3f\x03\x36\x07\xb1\xe0\x57\xfc\x6a\x96\x45\x2f\xf3\xd3\x11\xaf\xb4\xe5\x7b\xcf\x36\xa9\x38\xca\xfe\x6b\x85\x7c\x11\xd5\x97\x7c\xbe\x8c\xbb\xbf\x89\x6e\x46\xd0\x7f\xe3\x64\x17\xac\x22\xed\x6c\x32\x85\x5f\xa7\x43\x80\x7b\xd9\x40\xeb\x61\x0d\x16\x5b\xf8\x8e\xd9\x6d\x3f\x89\x84\xb5\x9d\x05\xe3\x94\x34\xa5\xf3\xb4\x9a\xbf\x5e\xbc\x59\x8a\x50\xe7\x92\xd0\xb3\xe9\xf5\x19\x7f\x1c\x4c\x14\x18\x27\x93\xcb\xca\x17\x33\xa3\x3d\x4d\x06\x88\xd6\x73\x67\x2b\xf4\x5e\x16\x18\x61\xcf\x2e\x70\x68\x03\x40\x39\xeb\x9d\x41\x6e\x5c\x91\xb0\xcf\x8e\x56\xc0\xe0\x7f\xc0\x7b\xb4\xc4\x63\x55\x39\x50\xec\x9a\x99\xa4\x46\x6f\x26\x53\x4e\xf8\x40\xc9\x00\xc5\x65\x5d\xa3\xcd\xb7\x2e\x89\xae\x06\x9c\xc7\xeb\xd1\x68\xc0\xbe\x6c\x82\x7d\x47\x36\x0a\x28\xa3\xd5\xdd\x60\x36\x43\x4f\x63\x5e\xc4\x59\x00\x13\x11\xcd\x5e\xbd\x58\x87\xae\x14\x9c\xfb\xff\xa3\x83\x1b\x8c\x71\x7b\xc2\x1c\x3a\xee\x99\xf9\xc7\xa1\xad\xe7\x7d\xb7\x8b\xa1\xbd\x04\x98\x8a\xd3\x63\x4e\x33\x97\x1f\x9e\x62\x4c\xe3\xac\x75\xbe\x66\x4f\x93\xee\x32\x0f\xdd\xd3\x8f\x5c\x6d\xeb\x40\x40\x87\x1a\xd7\x2c\x0b\x44\xce\xb2\x1e\x7b\xec\x95\xc5\xac\x4c\x88\x57\x37\xd1\x50\x64\xf5\xe8\xbe\xce\x51\x3f\xd6\xeb\x7f\xa2\xdf\x01\x00\x00\xff\xff\x94\x64\xc9\xca\x5c\x03\x00\x00")

func assets_index_html_tmpl_bytes() ([]byte, error) {
	return bindata_read(
		_assets_index_html_tmpl,
		"assets/index.html.tmpl",
	)
}

func assets_index_html_tmpl() (*asset, error) {
	bytes, err := assets_index_html_tmpl_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "assets/index.html.tmpl", size: 860, mode: os.FileMode(420), modTime: time.Unix(1455391058, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/index.html.tmpl": assets_index_html_tmpl,
}

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
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"assets": &_bintree_t{nil, map[string]*_bintree_t{
		"index.html.tmpl": &_bintree_t{assets_index_html_tmpl, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

