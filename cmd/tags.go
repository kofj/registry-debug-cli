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
	"github.com/kofj/registry-debug-cli/pkg/docker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "list tags",
	Long:  ``,
	Run:   dockerTagsListHandler,
}

var tagsListRepo string

func init() {
	dockerCmd.AddCommand(tagsCmd)

	rootCmd.MarkPersistentFlagRequired("repository")

}

func dockerTagsListHandler(cmd *cobra.Command, args []string) {
	var endpoint = viper.GetString("endpoint")
	var tls = viper.GetBool("tls")
	var insecure = viper.GetBool("insecure")
	var repository = viper.GetString("repository")
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

	tags, err := hub.Tags(repository)
	if err != nil {
		logrus.WithField("repository", repository).WithError(err).Error("list tags failed")
		return
	}

	var log = logrus.WithField("endpoint", endpoint).WithField("repository", repository)
	log.WithField("total", len(tags)).Infoln("Listed tags")
	for k, tag := range tags {
		log.WithField("index", k).WithField("tag", tag).Infof("%s/%s:%s", endpoint, repository, tag)
	}
}
