module github.com/kofj/registry-debug-cli

go 1.14

require (
	github.com/Sirupsen/logrus v1.7.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20200213202729-31a86c4ab209
	github.com/goharbor/harbor/src v0.0.0-20201110082039-ebc3443da94e
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/selinux v1.6.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/vbatts/tar-split v0.11.1 // indirect
)

replace github.com/Sirupsen/logrus v1.7.0 => github.com/sirupsen/logrus v1.7.0
