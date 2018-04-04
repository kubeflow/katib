// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
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

package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dconfig = ".dlkctlconfig"
)

//PersistentFlags is used for pass persistent flag value to other command
type PersistentFlags struct {
	cfgFile  string
	endpoint string
	registry string
	docker   string
	username string
}

var cfgFile string

// NewRootCommand represents the base command when called without any subcommands
func NewRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "dlkctl",
		Short: "dlk cli program",
		Long:  `this is dlk control cli program using cobra framework`,
	}

	//initialize config
	initFlag(cmd)

	//add command
	//POST
	cmd.AddCommand(NewCommandRun())
	//GET
	cmd.AddCommand(NewCommandGet())
	cmd.AddCommand(NewCommandLogs())
	//PUT

	//DELETE
	cmd.AddCommand(NewCommandDel())

	//MISC
	cmd.AddCommand(NewCommandVersion())

	//Generate bash completion file
	cmd.AddCommand(NewCommandBashCmp())

	return cmd
}

//CheckPersistentFlags check values is not empty and retun values
func CheckPersistentFlags() (*PersistentFlags, error) {
	var err error
	//get home dir for error message
	i, _ := user.Current()
	conf := i.HomeDir + "/" + dconfig

	//dlkmanager REST API endpoint
	e := viper.GetString("endpoint")
	if e == "" {
		err = errors.New("dlkmanager REST API endpoint is not specified,use --endpoint or edit " + conf + " to provide value")
		return nil, err
	}

	//registry endpoint
	r := viper.GetString("registry")
	if r == "" {
		err = errors.New("registry endpoint is not specified,use --registry or edit " + conf + " to provide value")
		return nil, err
	}

	//docker daemon API listen ip
	d := viper.GetString("docker")
	if d == "" {
		err = errors.New("docker daeom API endpoint is not specified,use --docker or edit " + conf + " to provide value")
		return nil, err
	}

	//username
	u := viper.GetString("user")
	if u == "" {
		err = errors.New("username is not specified,use --user or edit " + conf + " to provide value")
		return nil, err
	}
	rtn := PersistentFlags{
		endpoint: e,
		registry: r,
		docker:   d,
		username: u,
	}

	return &rtn, err
}

//initFlag manage persistent flags
func initFlag(cmd *cobra.Command) {
	cobra.OnInitialize(initConfig)
	// add Pesistent flags
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dlkctlconfig)")
	cmd.PersistentFlags().String("endpoint", "localhost:1323", "dlkmanager API endpoint")
	cmd.PersistentFlags().String("registry", "localhost:5000", "docker registry endpoint")
	cmd.PersistentFlags().String("docker", "localhost:2375", "docker daemon API listen address")
	cmd.PersistentFlags().String("user", "user", "username (default= $HOME/"+dconfig)

	//bind viper
	viper.BindPFlag("endpoint", cmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("registry", cmd.PersistentFlags().Lookup("registry"))
	viper.BindPFlag("docker", cmd.PersistentFlags().Lookup("docker"))
	viper.BindPFlag("user", cmd.PersistentFlags().Lookup("user"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var c string
	if cfgFile != "" {
		// Use config file from the flag.
		c = cfgFile
		// check existance
		_, err := os.Stat(c)
		if err != nil {
			fmt.Printf("ERROR file not found: %s\n", c)
			os.Exit(1)
		}

	} else {
		// Find home directory.
		u, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c = u.HomeDir + "/" + dconfig

		// check existance
		_, err = os.Stat(c)
		if err != nil {
			fmt.Println("config file not found")
			fmt.Printf("Generate Config File : %s\n", c)
			createConfig(c)
			fmt.Println("Edit the config and re-excute command")
			os.Exit(0)
		}

	}

	viper.SetConfigFile(c)
	viper.SetConfigType("toml")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}

}

//createConfig generate config file
func createConfig(c string) {
	//create new conig and write out info
	w, err := os.Create(c)
	if err != nil {
		fmt.Printf("Cannot Generate configfile : %s", dconfig)
		os.Exit(1)
	}
	user, err := user.Current()

	fmt.Fprintf(w, "#dlkmanager api server endpoint\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "endpoint", "localhost:1323")
	fmt.Fprintf(w, "#docker registry address\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "registry", "localhost:5000")
	fmt.Fprintf(w, "#docker daemon API listen address\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "docker", "localhost:2375")
	fmt.Fprintf(w, "#username\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "user", user.Username)

}
