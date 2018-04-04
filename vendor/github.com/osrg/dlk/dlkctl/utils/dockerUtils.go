package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//DockerImage is data domain for docker image
type DockerImage struct {
	ID          string   `json:"Id"`
	ParentID    string   `json:"ParentID"`
	RepoTags    []string `json:"RepoTags"`
	RepoDigests []string `json:"RepoDigests"`
	Created     int      `json:"Created"`
	Size        int      `json:"Size"`
	VirtualSize int      `json:"VirtualSize"`
	SharedSize  int      `json:"SharedSize"`
	Containers  int      `json:"Containers"`
}

//Tags is structure for store docker registry Search API result
type Tags struct {
	SearchName string   `json:"name"`
	TagNames   []string `json:"tags"`
}

// Repositories is list of docker images on private registry
type Repositories struct {
	ImgNames []string `json:"repositories"`
}

//DockerClient provides docker operation functions
type DockerClient struct {
	Cli               *http.Client
	DockerInterface   string
	RegistryInterface string
}

// GetImageListOnRegistry method related constructor
const (
	Catalog = "/v2/_catalog"              // catalog
	N       = "n=100"                     // Index and number to get images
	Last    = "last="                     // last image name index
	Errmsg  = "Can't get last image name" // error message
)

//NewDockerClient returns docker client Strunct instance
func NewDockerClient(dockerAPI string, registryAPI string) *DockerClient {
	//insecure setting
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//init Client with default value
	cli := &http.Client{Transport: tr}
	dockerInterface := dockerAPI
	registryInterface := registryAPI
	rtn := DockerClient{
		Cli:               cli,
		DockerInterface:   dockerInterface,
		RegistryInterface: registryInterface,
	}

	return &rtn
}

// ListImage equals to docker images command
func (c *DockerClient) ListImage() ([]DockerImage, error) {
	//construct REST API
	url := "http://" + c.DockerInterface + "/images/json"

	//create request GET + url
	req, _ := http.NewRequest("GET", url, nil)

	//send request
	res, err := c.Cli.Do(req)
	fmt.Println("SEND: " + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//get result
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var Images []DockerImage
	err = json.Unmarshal(body, &Images)

	return Images, err
}

//SearchOnRegistry return tag list of the specified images on registry
func (c *DockerClient) SearchOnRegistry(image string, user string) (*Tags, error) {
	url := fmt.Sprintf("https://%s/v2/%s/%s/tags/list", c.RegistryInterface, user, image)

	//create request GET + url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//send request
	res, err := c.Cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	//get result
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rtn := &Tags{}

	err = json.Unmarshal(body, rtn)
	if err != nil {
		return nil, err
	}
	//TODO display when  verbose flag is true
	//fmt.Println("Search word: " + rtn.SearchName)
	//for _, tag := range rtn.TagNames {
	//	fmt.Println(tag)
	//}

	return rtn, err

}

//GetImageListOnRegistry returns list of all images on registry
func (c *DockerClient) GetImageListOnRegistry() ([]Repositories, error) {

	var err error
	var last string // "&last=<last image name from previous response>"

	// https://<registry address>/v2/_catalog?n=<number>
	url := fmt.Sprintf("https://%s", c.RegistryInterface)
	url += Catalog + "?" + N

	rtn := []Repositories{} // image list

	for {
		//create request GET + url
		req, err := http.NewRequest("GET", url+last, nil)
		if err != nil {
			return nil, err
		}

		//send request
		res, err := c.Cli.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		//get result
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		rep := &Repositories{}

		err = json.Unmarshal(body, rep)
		if err != nil {
			return nil, err
		}

		rtn = append(rtn, *rep)

		// if there are no images to get on registry any more,
		// finish to get ones
		if len(res.Header["Link"]) == 0 {
			break
		}

		// if there are still images to get on registry,
		// retrieve last image name from GET response
		for _, l := range res.Header["Link"] {

			if strings.Contains(l, Catalog) &&
				strings.Contains(l, Last) &&
				strings.Contains(l, N) &&
				strings.Contains(l, "rel=\"next\"") {

				fr := strings.LastIndex(l, Last)
				if fr == -1 {
					err = errors.New(Errmsg)
					return nil, err
				}
				fr += len(Last)

				to := strings.LastIndex(l, "&"+N)
				if to == -1 {
					err = errors.New(Errmsg)
					return nil, err
				}

				img := l[fr:to]

				img = strings.Replace(img, "%2F", "/", 1)

				last = "&" + Last + img

				break
			}
		}
	}
	return rtn, err
}

//IsImageExistOnRegistry check whether image exist on registry or not
func (c *DockerClient) IsImageExistOnRegistry(imageName string, user string) (bool, error) {

	var tag string
	var image string

	//init return valiable with false. set it true when image found on registry
	exist := false

	//check imageName contains both image name and tag
	i := strings.LastIndex(imageName, ":")
	//get imagename and tag name from passed arg
	if i < 0 {
		image = imageName
		tag = "latest"
	} else {
		image = imageName[:i]
		tag = imageName[i+1:]
	}
	// search images using user/image as search keyword and get list of tags
	tags, err := c.SearchOnRegistry(image, user)
	if err != nil {
		return false, err
	}

	//check expected tag exist or not. on exists,return true
	for _, item := range tags.TagNames {
		if item == tag {
			exist = true
			break
		}
	}

	return exist, err
}

//IsImageExistLocally check whether image exist locally or not
func (c *DockerClient) IsImageExistLocally(target string) (bool, error) {

	//if passed name not contains tag,then add :latest
	t := target
	if !strings.Contains(target, ":") {
		t += ":latest"
	}

	//get local docker image list
	list, err := c.ListImage()
	if err != nil {
		return false, err
	}
	rtn := false

	// if list contains target(t) image,then return true
	for _, images := range list {
		for _, image := range images.RepoTags {
			if image == t {
				rtn = true
				break
			}

		}

	}
	return rtn, err
}

//BuildNewImage generate new docker image based on passed parameter
func (c *DockerClient) BuildNewImage(tag string, contextTar *os.File) error {
	//set url
	url := fmt.Sprintf("http://%s/build?t=%s", c.DockerInterface, tag)

	//send build request
	//TODO display when  verbose flag is true
	//fmt.Printf("BUILD REQUEST: %s, binary: %s\n", url, contextTar.Name())
	resp, err := http.Post(url, "binary/octet-stream", contextTar)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Read result
	s, err := ioutil.ReadAll(resp.Body)
	//error handling
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("docker build error: context=%s \n%s", contextTar.Name(), s)
	}
	return err
}

//PushImage push image
func (c *DockerClient) PushImage(image string) error {
	url := fmt.Sprintf("http://%s/images/%s/push", c.DockerInterface, image)
	auth := fmt.Sprintf("{serveraddress: %s}", c.RegistryInterface)

	//create request GET + url
	//TODO display when  verbose flag is true
	//fmt.Println("PUSH: " + url)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Registry-Auth", auth)

	//send request
	res, err := c.Cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//get result
	_, err = ioutil.ReadAll(res.Body)
	//TODO verbose
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Push Error: status code: %s, url:%s", res.Status, url)
	} else {
		fmt.Println(res.Status)
	}

	return err

}

//GetImageDigest return passed docker image digest on registry
func (c *DockerClient) GetImageDigest(image string, tag string) (string, error) {

	//remove "/" if imagename start with it
	img := image
	if strings.HasPrefix(image, "/") {
		img = image[1:]
	}

	//remove ":" if tagname start with it
	tg := tag
	if strings.HasPrefix(image, ":") {
		img = image[1:]
	}

	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", c.RegistryInterface, img, tg)
	req, _ := http.NewRequest("GET", url, nil)

	//header
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	//send request
	res, err := c.Cli.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	//check response code
	if res.StatusCode != http.StatusOK {
		fmt.Println(url)
		return "", fmt.Errorf("no such image, image: %s tag: %s", img, tg)
	}

	digest := res.Header.Get("Docker-Content-Digest")

	//should be dead code
	if digest == "" {
		return "", fmt.Errorf("No 'Docker-Content-Digest' in response header, url: %s", url)
	}

	return digest, nil
}

//SeparateTagrepoIntoImageAndTag return image name and tag separately
func (c *DockerClient) SeparateTagrepoIntoImageAndTag(repotag string) (string, string) {

	tag := ""
	image := repotag

	//remove registry url
	if strings.HasPrefix(image, c.RegistryInterface) {
		image = image[len(c.RegistryInterface):]
	}

	if strings.HasPrefix(image, "/") {
		image = image[1:]
	}

	fmt.Println(image)
	i := strings.Index(image, ":")
	if i != -1 {
		tag = image[i+1:]
		image = image[:i]
	}

	return image, tag
}

//GetImageDigest return passed docker image digest on registry
func (c *DockerClient) DeleteImage(image string, digest string) error {
	//remove "/" if imagename start with it
	img := image
	if strings.HasPrefix(image, "/") {
		img = image[1:]
	}

	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", c.RegistryInterface, img, digest)
	req, _ := http.NewRequest("DELETE", url, nil)

	//send request
	res, err := c.Cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//check response code
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("delete image error image: %s digest:%s code: %s", img, digest, res.Status)
	}

	return nil
}
