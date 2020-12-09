/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"

	"github.com/kofj/registry-debug-cli/pkg/docker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "inspect tag's manifest",
	Long: `inspect tag's manifest

Example: 
1. special tag: registry-debug-cli docker manifest -t 202012
2. default latest tag: registry-debug-cli docker manifest
	`,
	Run: dockerManifestListHandler,
}

func init() {
	dockerCmd.AddCommand(manifestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dockerManifestListHandler(cmd *cobra.Command, args []string) {
	var endpoint = viper.GetString("endpoint")
	var tls = viper.GetBool("tls")
	var insecure = viper.GetBool("insecure")
	var repository = viper.GetString("repository")
	var reference = viper.GetString("tag")
	var username = viper.GetString("username")
	var password = viper.GetString("password")

	logrus.WithField("endpoint", endpoint).
		WithField("repository", repository).
		WithField("tag", tagName).
		WithField("size", blobSize).
		WithField("username", username).
		Warn("Info")
	hub, err := docker.New(endpoint, username, password, tls, insecure)
	if err != nil {
		logrus.WithError(err).Error("hub failed")
		return
	}

	if reference == "" {
		reference = "latest"
		logrus.Warningln("Use latest tag")
	}

	manifest, err := hub.ManifestV2(repository, reference)
	if err != nil {
		logrus.WithError(err).Errorln("get ManifestV2 failed")
		return
	}
	logrus.WithField("DeserializedManifest", manifest).Infoln("DeserializedManifest")
	bytes, err := json.MarshalIndent(manifest, "  ", "")
	logrus.WithError(err).Infoln(string(bytes))

	hasConfig, err := hub.HasBlob(repository, manifest.Config.Digest)
	logrus.WithError(err).WithField("hasConfig", hasConfig).Infoln("hasConfig")
}
