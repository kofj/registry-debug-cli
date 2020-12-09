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
	"github.com/kofj/registry-debug-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	countTags bool
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: dockerRepoListHandler,
}

func init() {
	reposCmd.Flags().BoolVarP(&countTags, "count-tags", "c", false, "count tags")

	dockerCmd.AddCommand(reposCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reposCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reposCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dockerRepoListHandler(cmd *cobra.Command, args []string) {
	var endpoint = viper.GetString("endpoint")
	var tls = viper.GetBool("tls")
	var insecure = viper.GetBool("insecure")
	var repository = viper.GetString("repository")
	var username = viper.GetString("username")
	var password = viper.GetString("password")

	logrus.WithField("endpoint", endpoint).
		WithField("repository", repository).
		WithField("countTags", countTags).
		WithField("username", username).
		Warn("Info")
	hub, err := docker.New(endpoint, username, password, tls, insecure)
	if err != nil {
		logrus.WithError(err).Error("hub failed")
		return
	}

	repos, err := hub.Repositories()
	if err != nil {
		logrus.WithError(err).Error("List Repositories failed")
		return
	}
	logrus.WithField("total", len(repos)).WithField("repos", repos).Infoln("Listed repos")

	if countTags {
		var runner = utils.NewLimitedConcurrentRunner(10)
		var tags []string
		var tagsCounter int
		for i := range repos {
			index := i
			runner.AddTask(func() error {
				t, err := hub.Tags(repos[index])
				if err != nil {
					return err
				}
				tags = append(tags, t...)
				tagsCounter += len(t)
				return nil
			})
		}
		err = runner.Wait()
		if err != nil {
			logrus.WithError(err).Error("Count tags of all repository failed")
			return
		}

		logrus.
			WithField("total", tagsCounter).
			WithField("len", len(tags)).
			// WithField("tags", tags).
			Infoln("Count tags")
	}
}
