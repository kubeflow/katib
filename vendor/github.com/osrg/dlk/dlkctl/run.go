// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/osrg/dlk/dlkctl/utils"
	"github.com/osrg/dlk/dlkmanager/api"

	"github.com/spf13/cobra"
)

const (
	dockerImage    = "tensorflow/tensorflow"
	dockerImageGpu = "tensorflow/tensorflow:latest-gpu"
	tempDir        = "/tmp/dlkctl"
)

type runConfig struct {
	script *os.File
	params utils.Params
	pf     *PersistentFlags
}

// Image Names
type ImageNames struct {
	psImage     string
	workerImage string
}

//NewCommandRun generate run cmd
func NewCommandRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <workload.py> | --image <imagename> | --psRawImage <ps raw image name> --workerRawImage <worker raw image name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "create new-learningTask",
		Long:  `create new-learningTask`,
		Run:   runWorkflow,
	}

	//set local flag
	utils.AddImageFlag(cmd)
	utils.AddNameSpaceFlag(cmd)
	utils.AddNameFlag(cmd)
	utils.AddSchedulerFlag(cmd)
	utils.AddNrPsFlag(cmd)
	utils.AddNrWorkerFlag(cmd)
	utils.AddGpuFlag(cmd)
	utils.AddPsRawImageFlag(cmd)
	utils.AddWorkerRawImageFlag(cmd)
	utils.AddDryRunFlag(cmd)
	utils.AddBaseImageFlag(cmd)
	utils.AddEntryPointFlag(cmd)
	utils.AddParametersFlag(cmd)
	utils.AddTimeoutFlag(cmd)
	utils.AddPvcFlag(cmd)
	utils.AddMountPathFlag(cmd)
	utils.AddPriorityFlag(cmd)
	//add subcommand
	return cmd
}

//Main Proceduer of run command
func runWorkflow(cmd *cobra.Command, args []string) {

	//parameter check,init parameters
	//if neither workload file nor docker imagename is specified,then show help
	fmt.Println("*** CHECK PARAMS ***")
	rc := runConfig{}
	err := rc.checkParams(cmd, args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cli := utils.NewDockerClient(rc.pf.docker, rc.pf.registry)

	exist := false
	if rc.params.Image != "" {
		//search the docker image on private registry
		exist, err = cli.IsImageExistOnRegistry(rc.params.Image, rc.pf.username)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Completed")
	}

	rc.displayParams()

	//image names which are used for pulling images
	//from private registry or docker hub
	var imgNames ImageNames
	//if image name or raw image name is specified and workload file is not specified
	if rc.script == nil {
		// raw image is specified
		if rc.params.PsRawImage != "" || rc.params.WorkerRawImage != "" {
			// if raw image name is not specifed, default image is used
			if rc.params.PsRawImage != "" {
				imgNames.psImage = rc.params.PsRawImage
			} else {
				imgNames.psImage = dockerImage
			}
			if rc.params.WorkerRawImage != "" {
				imgNames.workerImage = rc.params.WorkerRawImage
			} else {
				if rc.params.Gpu > 0 {
					imgNames.workerImage = dockerImageGpu
				} else {
					imgNames.workerImage = dockerImage
				}
			}
			// if image name is specified and image exist on registry
		} else if exist {
			imgNames.psImage = fmt.Sprintf("%s/%s/%s", cli.RegistryInterface, rc.pf.username, rc.params.Image)
			imgNames.workerImage = imgNames.psImage
		} else {
			fmt.Printf("image: %s is not exists on the registry(%s)\n", rc.params.Image, cli.RegistryInterface)
			os.Exit(1)
		}

	} else { //if workload file is specified
		//create docker images
		fmt.Println("*** CREATE Docker Image ***")
		imgNames, err = rc.createDockerImage(cli)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%-30s : %s\n", "REPOSITORY NAME for NON-GPU", imgNames.psImage)
		if rc.params.Gpu > 0 {
			fmt.Printf("%-30s : %s\n", "REPOSITORY NAME for GPU", imgNames.workerImage)
		}
		fmt.Println("Completed")
		//push images from local to private registry
		fmt.Println("*** Push Docker Image To Registry ***")
		if !rc.params.DryRun {
			// push non-gpu image
			err = cli.PushImage(imgNames.psImage)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			// push gpu image
			if rc.params.Gpu > 0 {
				err = cli.PushImage(imgNames.workerImage)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
		fmt.Println("Completed")

	}

	//Send LearningTask Request to API server
	fmt.Println("*** Send LearningTask Request to API Server ***")
	fmt.Printf("LearningTask Name: %s\n", rc.params.Name)
	err = rc.sendLearningTaskRequest(imgNames)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Completed. \ndone")
}

//checkParams check args and flag vailidity and return runConfig struct
func (rc *runConfig) checkParams(cmd *cobra.Command, args []string) error {
	var err error

	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err = CheckPersistentFlags()
	if err != nil {
		return err
	}

	//check Flags using common parameter checker
	var params utils.Params
	params, err = utils.CheckFlags(cmd)
	if err != nil {
		return err
	}

	// check if neither workload.py, image nor raw images is specified
	if len(args) == 0 && params.Image == "" &&
		params.PsRawImage == "" && params.WorkerRawImage == "" {
		err = errors.New("either workload.py, image name, or raw image names is required to execute request")
		return err
	}

	var script *os.File

	//get time in order to use in name auto-generation prcess
	now := time.Now()

	//if workload file is specfied,then open it
	if len(args) == 1 {
		script, err = os.Open(args[0])
		if err != nil {
			return err
		}
	}

	//check image name
	//if image name is not specified, automatically generate it
	//<scriptname>-yy-mm-dd-hh-MM-ss
	if params.Image == "" && params.PsRawImage == "" && params.WorkerRawImage == "" {

		var s string
		sname := filepath.Base(script.Name())
		i := strings.LastIndex(sname, ".")
		if i != -1 {
			s = sname[0:i]
		} else {
			s = sname
		}
		params.Image = fmt.Sprintf("%s-%d-%d-%d-%d-%d-%d",
			s, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}

	// if learningTask name is not specified,then generate
	if params.Name == "" {
		params.Name = fmt.Sprintf("%s-%d-%d-%d-%d-%d-%d", "learningtask", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}

	rc.script = script
	rc.params = params
	rc.pf = pf

	return err
}

//createDockerImage build Docker images using docker REST API
//return: image names(ImageNames)
func (rc *runConfig) createDockerImage(cli *utils.DockerClient) (imgnm ImageNames, err error) {

	//push image names for ps and worker
	var imageNames ImageNames
	imageNames.psImage = fmt.Sprintf("%s/%s/%s", cli.RegistryInterface, rc.pf.username, rc.params.Image)
	if rc.params.Gpu > 0 {
		imageNames.workerImage = imageNames.psImage + ":latest-gpu"
	} else {
		imageNames.workerImage = imageNames.psImage
	}

	//dry run
	if rc.params.DryRun {
		return imageNames, nil
	}

	//TO build using Docker API,api require context file witch is tar file contains dockerfile
	//generate dockerfile and compress it. this function return *File of generated tar
	//build non-gpu image
	gpuf := false
	cctx, err := rc.generateDockerContext(gpuf)
	if err != nil {
		return ImageNames{}, err
	}

	ccfile, err := os.Open(cctx)
	if err != nil {
		return ImageNames{}, err
	}

	err = cli.BuildNewImage(imageNames.psImage, ccfile)
	if err != nil {
		return ImageNames{}, err
	}

	//build gpu image
	if rc.params.Gpu > 0 {
		gpuf = true
		// file descriptor is moved back to top of
		// workload file
		rc.script.Seek(0, os.SEEK_SET)
		gctx, err := rc.generateDockerContext(gpuf)
		if err != nil {
			return ImageNames{}, err
		}

		gcfile, err := os.Open(gctx)
		if err != nil {
			return ImageNames{}, err
		}

		err = cli.BuildNewImage(imageNames.workerImage, gcfile)
		if err != nil {
			return ImageNames{}, err
		}
	}

	return imageNames, err
}

//generateDockerContext create dockerfile and compress files required by build REST api
func (rc *runConfig) generateDockerContext(gpuf bool) (string, error) {

	//set temporary output dir,and generate it
	dir := tempDir + "/" + rc.params.Image
	if gpuf {
		dir += "-g"
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	//create dockerfile
	//set base image
	//high priority: --baseImage flag
	var base string
	if rc.params.BaseImage != "" {
		base = rc.params.BaseImage

	} else if rc.params.Gpu > 0 { //if #of Gpu is specified,
		if gpuf { // for gpu, use tensorflow-with-gpu image
			base = dockerImageGpu
		} else { // for non-gpu, use tensorflow-non-gpu image
			base = dockerImage
		}
	} else { // use tensorflow-non-gpu image
		base = dockerImage
	}
	//create dockerfile
	df := []string{}
	df = append(df, "FROM "+base)
	df = append(df, "MAINTAINER dlkctl")
	df = append(df, "RUN mkdir /script")
	df = append(df, "ADD "+filepath.Base(rc.script.Name())+" /script/")
	df = append(df, "WORKDIR /script")

	file := dir + "/" + "Dockerfile"
	ofile, err := os.Create(file)
	if err != nil {
		return "", err
	}

	w := bufio.NewWriter(ofile)
	for _, str := range df {
		fmt.Fprintln(w, str)
	}
	w.Flush()

	//copy workload script to context dir
	dst, err := os.Create(dir + "/" + filepath.Base(rc.script.Name()))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dst, rc.script)
	dst.Close()
	//compress it
	path, err := makeTar(dir)

	return path, err
}

//makeTar function compress a directory specified by passed argument into tar file and produce output within the pash
func makeTar(dirPath string) (string, error) {
	//create output tar file
	dstStr := dirPath + "/context.tar"
	dst, err := os.Create(dstStr)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	//get all items within src dir
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	//add each item into tar file
	var fileWriter io.WriteCloser = dst
	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {

		//if items is dir then skip;otherwise add to tar file
		if fileInfo.IsDir() || fileInfo.Name() == "context.tar" {
			continue
		}
		file, err := os.Open(dirPath + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			return "", err
		}

		defer file.Close()

		// set tar header
		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		if err != nil {
			return "", err
		}

		//copy file
		_, err = io.Copy(tarfileWriter, file)
		if err != nil {
			return "", err
		}

	}

	return dstStr, err
}

//sendLearningTaskRequest send REST API Request using json
func (rc *runConfig) sendLearningTaskRequest(imageName ImageNames) error {
	//set url
	url := fmt.Sprintf("http://%s/learningTask", rc.pf.endpoint)
	// construct json using runconfig parameter
	j := api.LTConfig{
		PsImage:     imageName.psImage,
		WorkerImage: imageName.workerImage,
		Ns:          rc.params.Ns,
		Scheduler:   rc.params.Scheduler,
		Name:        rc.params.Name,
		NrPS:        rc.params.NrPs,
		NrWorker:    rc.params.NrWorker,
		Gpu:         rc.params.Gpu,
		DryRun:      rc.params.DryRun,
		EntryPoint:  rc.params.EntryPoint,
		Parameters:  rc.params.Parameters,
		Timeout:     rc.params.Timeout,
		Pvc:         rc.params.Pvc,
		MountPath:   rc.params.MountPath,
		Priority:    rc.params.Priority,
		User:        rc.pf.username,
	}

	//encode json
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	//send REST API Request
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO verbose
	//fmt.Println("response Status:", resp.Status)
	return err
}

func (rc *runConfig) displayParams() {
	fmt.Println("** exec parameters ************")
	if rc.script != nil {
		fmt.Printf("| %-30s : %s\n", "workload script", rc.script.Name())
	}
	fmt.Printf("| %-30s : %s\n", "dlkmanager endpoint", rc.pf.endpoint)
	fmt.Printf("| %-30s : %s\n", "docker daemon API endpoint", rc.pf.docker)
	fmt.Printf("| %-30s : %s\n", "docker registry endpoint", rc.pf.registry)
	fmt.Printf("| %-30s : %s\n", "user name", rc.pf.username)
	if rc.params.Image != "" {
		fmt.Printf("| %-30s : %s\n", "docker image", rc.params.Image)
	}
	if rc.params.BaseImage != "" {
		fmt.Printf("| %-30s : %s\n", "docker base image", rc.params.BaseImage)
	}
	if rc.params.PsRawImage != "" {
		fmt.Printf("| %-30s : %s\n", "PS Raw docker image", rc.params.PsRawImage)
	}
	if rc.params.WorkerRawImage != "" {
		fmt.Printf("| %-30s : %s\n", "Worker Raw docker image", rc.params.WorkerRawImage)
	}
	if rc.params.EntryPoint != "" {
		fmt.Printf("| %-30s : \"%s\"\n", "Entry Point", rc.params.EntryPoint)
	}

	if rc.params.Parameters != "" {
		fmt.Printf("| %-30s : \"%s\"\n", "container exec Parameters", rc.params.Parameters)
	}
	fmt.Printf("| %-30s : %s\n", "k8s namespace", rc.params.Ns)
	fmt.Printf("| %-30s : %s\n", "learningTask name", rc.params.Name)
	fmt.Printf("| %-30s : %s\n", "scheduler name", rc.params.Scheduler)
	fmt.Printf("| %-30s : %d\n", "# of PS", rc.params.NrPs)
	fmt.Printf("| %-30s : %d\n", "# of worker", rc.params.NrWorker)
	fmt.Printf("| %-30s : %d\n", "gpu", rc.params.Gpu)
	fmt.Printf("| %-30s : %t\n", "dry-run", rc.params.DryRun)
	fmt.Printf("| %-30s : %d\n", "timeout", rc.params.Timeout)
	fmt.Printf("| %-30s : %s\n", "persistent volume claim", rc.params.Pvc)
	fmt.Printf("| %-30s : %s\n", "nfs mount path", rc.params.MountPath)
	fmt.Printf("| %-30s : %d\n", "priority", rc.params.Priority)
}
