package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// GenerateBlob return reader
func GenerateBlob(index int, size int64) (*bytes.Buffer, error) {
	var gzbuf = &bytes.Buffer{}
	var gw = gzip.NewWriter(gzbuf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// add each file as needed into the current tar archive
	return gzbuf, addFile(tw, size)
}

func addFile(tw *tar.Writer, size int64) (err error) {
	// now lets create the header as needed for this file within the tarball
	header := new(tar.Header)
	header.Name = "kofj.random.bytes"
	header.Size = size
	var now = time.Now().UTC()
	header.AccessTime = now
	header.ChangeTime = now
	header.Mode = 0655

	var blobBuf = make([]byte, size)
	// write the header to the tarball archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	n, err := rand.Read(blobBuf)
	if err != nil {
		return
	}
	logrus.WithField("bytes", n).Info("Generate bytes")
	// copy the file data to the tarball
	var bfw = bytes.NewBuffer(blobBuf)
	if _, err := io.Copy(tw, bfw); err != nil {
		return err
	}

	return nil
}
