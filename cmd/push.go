package cmd

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	digest "github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kofj/registry-debug-cli/pkg/docker"
)

const (
	// B 1 Byte
	B uint = 1
	// KB 1 KBytes
	KB uint = 1 << (10 * iota)
	// MB 1 MBytes
	MB
	// GB 1 GBytes
	GB
)

var blobSizeString string
var blobSize uint
var randomTag bool
var tagName string
var blobBuf []byte

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "generate and push docker image",
	Long:  ``,
	Run:   dockerPushHandler,
}

func init() {
	dockerCmd.AddCommand(pushCmd)

	dockerCmd.PersistentFlags().StringVarP(&blobSizeString, "size", "s", "10MB", "Blob's size")
	dockerCmd.PersistentFlags().BoolVarP(&randomTag, "random-tag", "R", true, "Random tag name")

	viper.BindPFlag("size", dockerCmd.PersistentFlags().Lookup("size"))
	viper.BindPFlag("random-tag", dockerCmd.PersistentFlags().Lookup("random-tag"))

	cobra.MarkFlagRequired(dockerCmd.Flags(), "size")
}

func dockerPushHandler(cmd *cobra.Command, args []string) {
	blobSize = viper.GetSizeInBytes("size")
	if blobSize > 5*GB || blobSize <= KB {
		logrus.WithField("current", blobSize).WithField("string", blobSizeString).Error("Blob' size must be between 1KB and 5GB")

		dockerCmd.Help()
		return
	}

	if viper.GetString("tag") == "" && !randomTag {
		logrus.WithField("tag", tag).Error("Not set tag's name and not allow random")
		dockerCmd.Help()
		return
	}

	if randomTag {
		rand.Seed(time.Now().UnixNano())
		tagName = fmt.Sprint(time.Now().Format("20060102.150405."), rand.Intn(9999))
	} else {
		tagName = viper.GetString("tag")
	}

	// make blob buf
	blobBuf = make([]byte, blobSize)

	logrus.WithField("endpoint", viper.GetString("endpoint")).
		WithField("repository", viper.GetString("repository")).
		WithField("tag", tagName).
		Info("generate image")
	fmt.Println("push called")

	// genarete blobs
	gzipbuf, err := docker.GenerateBlob(1, int64(blobSize))
	if err != nil {
		logrus.WithError(err).Errorln("generate blob file")
		return
	}

	// temporary buffer
	tempbuf := bytes.NewBuffer(make([]byte, 0))

	// calc blob sha256
	var h = sha256.New()
	r := io.TeeReader(gzipbuf, tempbuf)
	n, err := io.Copy(h, r)
	if err != nil {
		logrus.WithField("bytes", n).WithError(err).Error("calc sha256 error")
		return
	}
	var bs = h.Sum(nil)
	var sha = fmt.Sprintf("%x", bs)
	logrus.WithField("bytes", n).WithField("digest", sha).Info("blob digest")

	var endpoint = viper.GetString("endpoint")
	var repository = viper.GetString("repository")
	var username = viper.GetString("username")
	var password = viper.GetString("password")
	var url = fmt.Sprintf("https://%s", endpoint)
	var blobSize = int64(tempbuf.Len())
	var blobDigest = digest.NewDigestFromHex("sha256", sha)

	logrus.WithField("endpoint", endpoint).
		WithField("repository", repository).
		WithField("tag", tagName).
		WithField("digest", blobDigest).
		WithField("size", blobSize).
		WithField("username", username).
		Warn("Info")

	hub, err := registry.NewInsecure(url, username, password)
	if err != nil {
		logrus.WithError(err).Error("hub failed")
		return
	}
	hub.Logf = func(format string, args ...interface{}) {
		logrus.Infof(format, args...)
	}

	pushBlob(repository, "normal blob", blobDigest, blobSize, tempbuf, hub)

	dockerConfig, err := docker.BuildConfigBytes(blobDigest)
	if err != nil {
		logrus.WithError(err).Error("build digest failed")
		return
	}

	logrus.WithField("config", string(dockerConfig)).Warn("docker image config")

	// calc config sha256
	var configDigest = digest.NewDigestFromHex("sha256", fmt.Sprintf("%x", sha256.Sum256(dockerConfig)))
	var configBuf = bytes.NewBuffer(dockerConfig)
	var configSize = int64(len(dockerConfig))
	pushBlob(repository, "docker config blob", configDigest, configSize, configBuf, hub)

	var blobsDescriptors = distribution.Descriptor{
		MediaType: schema2.MediaTypeLayer,
		Size:      blobSize,
		Digest:    blobDigest,
	}
	var dockerImageManifest = docker.BuildManifest(configSize, configDigest, blobsDescriptors)
	time.Sleep(3e9)
	err = hub.PutManifest(repository, tagName, dockerImageManifest)
	var dockerImageDigest string
	if err != nil {
		logrus.WithField("image digest", dockerImageDigest).WithError(err).Error("Push docker image manifest failed")
		return
	}

	logrus.WithField("image digest", dockerImageDigest).
		WithField("image name", fmt.Sprintf("%s/%s:%s", endpoint, repository, tagName)).Info("Push docker image success")

}

func pushBlob(repository, comments string, digest digest.Digest, size int64, buf *bytes.Buffer, client *registry.Registry) {

	// check

	exists, err := client.HasBlob(repository, digest)
	if err != nil {
		logrus.WithError(err).Error("Check blob exist failed")
		return
	}

	// push blob
	if !exists {
		client.UploadBlob(repository, digest, buf)
		if err != nil {
			logrus.WithField("digest", digest).WithError(err).Error("Upload blob failed")
			return
		}
		recheck, err := client.HasBlob(repository, digest)
		if err != nil {
			logrus.WithField("digest", digest).WithError(err).Error("Recheck uploaded blob failed")
		}
		if recheck {
			logrus.WithField("digest", digest).Warnln("Recheck uploaded blob success")
		} else {
			logrus.WithField("digest", digest).Error("Recheck uploaded blob not exist")
		}
	}

	logrus.WithField("blob", digest).WithField("comments", comments).Info("push blob success")

}
