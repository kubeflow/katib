package configs

import (
	"fmt"
	"os"
	"os/user"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//PersistentFlags is used for pass persistent flag value to other command
type PersistentFlags struct {
	cfgFile      string
	Addr         string
	Mt           int
	Pbt          int
	Nbt          int
	Pcs          int
	FakeMachines bool
	Nm           int
	Scheduler    string
	Ns           string
	Logdir       string
	Loglvl       int
	Username     string
}

const (
	defaultConfig = ".dlkmanagerconfig"
)

var (
	cfgFile string
	Pflg    PersistentFlags
)

//SetFlags init cobra flags
func SetFlags(cmd *cobra.Command) {
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+defaultConfig)
	cmd.PersistentFlags().String("addr", "localhost:6443", "k8s api endpoint")
	cmd.PersistentFlags().Int("mt", 1000, "maximum number of tasks per machine")
	cmd.PersistentFlags().Int("pbt", 2, "pods batch timeout in seconds")
	cmd.PersistentFlags().Int("nbt", 2, "node batch timeout in seconds")
	cmd.PersistentFlags().Int("pcs", 5000, "pod channel size in client's pod informer")
	// Fake the machine topology, only useful for testing when the API server has no nodes
	cmd.PersistentFlags().Bool("fakeMachines", false, "fake the machine topology if running without a real cluster")
	cmd.PersistentFlags().Int("nm", 2, "number of machines, only needed if faking the resource topology")
	cmd.PersistentFlags().String("scheduler", "default-scheduler", "a name of this scheduler")
	cmd.PersistentFlags().String("ns", "default", "k8s default namespace")
	cmd.PersistentFlags().String("logdir", "/tmp/", "log directory")
	cmd.PersistentFlags().Int("loglvl", 4, "log level (default is info level)")
	cmd.PersistentFlags().String("user", "user", "username (default = $HOME/"+defaultConfig)

	//bind viper
	viper.BindPFlag("addr", cmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("mt", cmd.PersistentFlags().Lookup("mt"))
	viper.BindPFlag("pbt", cmd.PersistentFlags().Lookup("pbt"))
	viper.BindPFlag("nbt", cmd.PersistentFlags().Lookup("nbt"))
	viper.BindPFlag("pcs", cmd.PersistentFlags().Lookup("pcs"))
	viper.BindPFlag("fakeMachines", cmd.PersistentFlags().Lookup("fakeMachines"))
	viper.BindPFlag("nm", cmd.PersistentFlags().Lookup("nm"))
	viper.BindPFlag("scheduler", cmd.PersistentFlags().Lookup("scheduler"))
	viper.BindPFlag("ns", cmd.PersistentFlags().Lookup("ns"))
	viper.BindPFlag("logdir", cmd.PersistentFlags().Lookup("logdir"))
	viper.BindPFlag("loglvl", cmd.PersistentFlags().Lookup("loglvl"))
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
			log.WithFields(log.Fields{
				"File Name": c,
				"Error":     err,
			}).Error("file not found")
			os.Exit(1)
		}
	} else {
		// Find user home directory.
		u, err := user.Current()
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("current user not aquired")
			os.Exit(1)
		}

		// Set default path and config file
		c = u.HomeDir + "/" + defaultConfig

		// check existance
		_, err = os.Stat(c)
		if err != nil {
			log.Info("config file not found")
			log.Info(fmt.Sprintf("Generate Config File : %s", c))
			createConfig(c)
			log.Info("Edit the config and re-excute command")
			//os.Exit(0)
		}

	}

	viper.SetConfigFile(c)
	viper.SetConfigType("toml")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	} else {
		log.WithFields(log.Fields{
			"File Name": c,
			"Error":     err,
		}).Error("config file not found")
		panic(err)
	}

	Pflg.Addr = viper.GetString("addr")
	Pflg.Mt = viper.GetInt("mt")
	Pflg.Pbt = viper.GetInt("pbt")
	Pflg.Nbt = viper.GetInt("nbt")
	Pflg.Pcs = viper.GetInt("pcs")
	Pflg.FakeMachines = viper.GetBool("fakeMachines")
	Pflg.Nm = viper.GetInt("nm")
	Pflg.Scheduler = viper.GetString("scheduler")
	Pflg.Ns = viper.GetString("ns")
	Pflg.Logdir = viper.GetString("logdir")
	Pflg.Loglvl = viper.GetInt("loglvl")
	Pflg.Username = viper.GetString("user")
}

//createConfig generate config file
func createConfig(c string) {
	//create new config and write out info
	w, err := os.Create(c)
	if err != nil {
		log.WithFields(log.Fields{
			"File Name": c,
			"Error":     err,
		}).Error("Cannot Generate configfile")
		os.Exit(1)
	}
	user, err := user.Current()

	fmt.Fprintf(w, "#k8s api endpoint\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "addr", "localhost:6443")
	fmt.Fprintf(w, "#maximum number of tasks per machine\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "mt", "1000")
	fmt.Fprintf(w, "#pods batch timeout in seconds\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "pbt", "2")
	fmt.Fprintf(w, "#node batch timeout in seconds\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "nbt", "2")
	fmt.Fprintf(w, "#pod channel size in client's pod informer\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "pcs", "5000")
	fmt.Fprintf(w, "#fake the machine topology if running without a real cluster\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "fakeMachines", "false")
	fmt.Fprintf(w, "#number of machines, only needed if faking the resource topology\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "nm", "2")
	fmt.Fprintf(w, "#a name of this scheduler\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "scheduler", "default-scheduler")
	fmt.Fprintf(w, "#k8s default namespace\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "ns", "default")
	fmt.Fprintf(w, "#log directory\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "logdir", "/tmp/")
	fmt.Fprintf(w, "#log level\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "loglvl", "4")
	fmt.Fprintf(w, "#username\n")
	fmt.Fprintf(w, "%s = '%s'\n\n", "user", user.Username)
}
