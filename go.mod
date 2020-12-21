module github.com/kofj/registry-debug-cli

go 1.14

require (
	github.com/Microsoft/hcsshim v0.8.11 // indirect
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.0+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/symlink v0.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc92 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/vbatts/tar-split v0.11.1 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.7.0
	github.com/sirupsen/logrus => github.com/Sirupsen/logrus v1.7.0
)
