package main

import (
	"context"
	"flag"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/mlkube/katib/api"
)

var server = flag.String("s", "127.0.0.1:6789", "server address")
var confPath = flag.String("f", "", "config file path")

// var verbose = flag.Bool("v", false, "verbose output")

type ManagerAPI struct {
	StudyConf *pb.StudyConfig
}

func (m *ManagerAPI) Createstudy(conn *grpc.ClientConn, args []string) {
	log.Printf("req Createstudy\n")
	c := pb.NewManagerClient(conn)
	req := &pb.CreateStudyRequest{StudyConfig: m.StudyConf}
	r, err := c.CreateStudy(context.Background(), req)
	if err != nil {
		log.Fatalf("CreateStudy failed: %v", err)
	}
	log.Printf("CreateStudy: %v", r)
}

func (m *ManagerAPI) Stopstudy(conn *grpc.ClientConn, args []string) {
	log.Printf("req Stopstudy\n")
	c := pb.NewManagerClient(conn)
	req := &pb.StopStudyRequest{StudyId: args[1]}
	r, err := c.StopStudy(context.Background(), req)
	if err != nil {
		log.Fatalf("StopStudy failed: %v", err)
	}
	log.Printf("StopStudy: %v", r)
}

func (m *ManagerAPI) Getstudies(conn *grpc.ClientConn, args []string) {
	c := pb.NewManagerClient(conn)
	req := &pb.GetStudysRequest{}
	r, err := c.GetStudys(context.Background(), req)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
	}
	fmt.Printf("StudyID         \tName\tOwner\tRunningTrial\tCompletedTrial\n")
	for _, si := range r.StudyInfos {
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", si.StudyId, si.Name, si.Owner, si.RunningTrialNum, si.CompletedTrialNum)
	}
}
func main() {
	flag.Parse()

	if *confPath == "" && flag.Arg(0) == "Createstudy" {
		log.Fatalf("Missing -f <config file path> option")
	}

	log.Printf("connecting %s", *server)
	conn, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	var sc pb.StudyConfig
	var m ManagerAPI
	if *confPath != "" {
		buf, _ := ioutil.ReadFile(*confPath)
		err = yaml.Unmarshal(buf, &sc)
		log.Printf("study conf%v\n", sc)
		m = ManagerAPI{StudyConf: &sc}
	}
	method, ok := reflect.TypeOf(&m).MethodByName(
		string(strings.Title(flag.Arg(0))))
	if !ok {
		log.Fatalf("Method not found: %s", flag.Arg(0))
	}
	ma := []reflect.Value{
		reflect.ValueOf(&m),
		reflect.ValueOf(conn),
		reflect.ValueOf(flag.Args()),
	}

	method.Func.Call(ma)
}
