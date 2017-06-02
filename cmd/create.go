// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/jpisaac/alert-watch/config"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("watch create called")

		conf, err := config.New()
		if err != nil {
			log.Fatal(err)
		}

		var b bool
		b, err = cmd.Flags().GetBool("svc")
		if err == nil {
			conf.Resource.Services = b
		} else {
			log.Fatal("svc", err)
		}

		b, err = cmd.Flags().GetBool("dp")
		if err == nil {
			conf.Resource.Deployment = b
		} else {
			log.Fatal("deployments", err)
		}

		b, err = cmd.Flags().GetBool("po")
		if err == nil {
			conf.Resource.Pod = b
		} else {
			log.Fatal("po", err)
		}

		b, err = cmd.Flags().GetBool("rs")
		if err == nil {
			conf.Resource.ReplicaSet = b
		} else {
			log.Fatal("rs", err)
		}

		b, err = cmd.Flags().GetBool("rc")
		if err == nil {
			conf.Resource.ReplicationController = b
		} else {
			log.Fatal("rc", err)
		}

		if err = conf.Write(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	//watchCmd.AddCommand(createWatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().Bool("svc", false, "create watch for services")
	createCmd.Flags().Bool("dp", false, "create watch for deployments")
	createCmd.Flags().Bool("po", false, "ctreate watch for pods")
	createCmd.Flags().Bool("rc", false, "create watch for replication controllers")
	createCmd.Flags().Bool("rs", false, "create watch for replicasets")

}
