package main

import (
	"fmt"
	"os"

	"github.com/kubeflow/hp-tuning/dlk/dlkmanager/configs"
	"github.com/labstack/echo"       // Web server framework for REST API
	lgr "github.com/sirupsen/logrus" // logging framework
	"github.com/spf13/cobra"
)

var (
	log = lgr.New() // Creates a new logger
)

/*
// log file inital setting
func logSetting() (*os.File, error) {
	log.Formatter = new(lgr.TextFormatter) // log format setting
	log.Level = lgr.DebugLevel // log level setting

	logfile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")

	log.Out = logfile // log output to file

	return logfile, err
}
*/

func apiMain() {
	log.Formatter = new(lgr.TextFormatter) // log format setting
	log.Level = lgr.DebugLevel             // log level setting
	log.Out = os.Stdout                    // log output to console

	/*
		// log file initial setting
		lfile, err := logSetting()
		if err != nil {
			log.Error(err)
		}
	*/

	// Echo instance
	e := echo.New()

	// Routes
	e.POST("/learningTask", runLearningTask)
	e.GET("/learningTasks/:namespace", getLearningTasks)
	e.GET("/learningTask/:namespace/:lt", getLearningTask)
	e.GET("/learningTasks/logs/:namespace/:lt/:role", getLearningTaskLogs)
	e.PUT("/learningTask", updateLearningTask)
	e.DELETE("/learningTasks/:namespace/:lt", deleteLearningTask)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	// lfile.Close() // log file closed
}

func main() {

	//init command
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewRootCommand represents the base command when called without any subcommands
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dlkmanager",
		Short: "dlkmanager",
		Long:  `dlkmanager`,
		Run:   run,
	}

	//initialize config
	configs.SetFlags(cmd)

	//add command
	//Generate bash completion file
	cmd.AddCommand(NewCommandBashCmp())

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	apiMain()
	//schedMain()
}
