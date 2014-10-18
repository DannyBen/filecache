// Package filecache provides simple file system cache
package filecache

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

// Type Handler holds the cache directory and cache
// life in minutes
type Handler struct {
	Dir  string
	Life float64
}

// Set saves content to the cache
func (h Handler) Set(key string, data []byte) error {
	return ioutil.WriteFile(h.Filename(key), data, 0644)
}

// Get returns content from the cache
func (h Handler) Get(key string) []byte {
	file := h.Filename(key)
	result, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	fi, err := os.Stat(file)
	panicon(err)
	age := time.Since(fi.ModTime()).Minutes()
	if age > h.Life {
		err := os.Remove(file)
		panicon(err)
		return nil
	}
	return result
}

// Filename returns the full path to the cache file
func (h Handler) Filename(key string) string {
	return h.dir() + "/" + hash(key)
}

// dir creates and returns a valid cache directory
// temp dir is returned if Handler.Dir is not set
func (h Handler) dir() string {
	if h.Dir == "" {
		h.Dir = os.TempDir()
		return h.Dir
	}

	h.Dir = resolveHomeDir(h.Dir)
	yes, err := exists(h.Dir)
	panicon(err)
	if !yes {
		err := os.MkdirAll(h.Dir, 0766)
		panicon(err)
	}
	return h.Dir
}

// hash returns the checksum of a string
func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// resolveHomeDir replaces ~ with the users home directory
func resolveHomeDir(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return usr.HomeDir + strings.TrimPrefix(path, "~")
	}
	return path
}

func panicon(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
