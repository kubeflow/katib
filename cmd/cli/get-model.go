package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type getModelOpt struct {
	detail bool
	args   []string
}

//NewCommandGetModel generate get model cmd
func NewCommandGetModel() *cobra.Command {
	var opt getModelOpt
	cmd := &cobra.Command{
		Use:     "model",
		Args:    cobra.MaximumNArgs(2),
		Short:   "Display Model Info",
		Long:    `Display Information of saved model`,
		Aliases: []string{"md"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			getModel(cmd, &opt)
		},
	}
	cmd.Flags().BoolVarP(&opt.detail, "detail", "d", false, "Display detail information of Model")
	return cmd
}

func getModel(cmd *cobra.Command, opt *getModelOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}

	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	// Search study if Study ID or name is set
	req := &api.GetStudyListRequest{}
	r, err := c.GetStudyList(context.Background(), req)
	if err != nil {
		log.Fatalf("GetModels failed: %v", err)
	}
	if len(r.StudyOverviews) == 0 {
		log.Println("No Study fond")
		return
	}
	for _, si := range r.StudyOverviews {
		if len(opt.args) > 0 {
			if utf8.RuneCountInString(opt.args[0]) >= 7 {
				if !strings.HasPrefix(si.Id, opt.args[0]) {
					break
				}
			}
			if si.Name != opt.args[0] {
				break
			}
		}
		// Search Models from ModelDB
		mreq := &api.GetSavedModelsRequest{StudyName: si.Name}
		mr, err := c.GetSavedModels(context.Background(), mreq)
		if err != nil {
			log.Fatalf("GetModels failed: %v", err)
			return
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', tabwriter.TabIndent)
		fmt.Printf("Study %v Owner %v Saved Model Num %v:\n", si.Name, si.Owner, len(mr.Models))
		if opt.detail {
			for _, m := range mr.Models {
				if len(opt.args) > 1 {
					if !strings.HasPrefix(m.WorkerId, opt.args[1]) {
						continue
					}
				}
				fmt.Printf("WorkerID :%v\n", m.WorkerId)
				fmt.Printf("Model Path: %s\n", m.ModelPath)
				fmt.Println("Parameters:")
				for _, p := range m.Parameters {
					fmt.Fprintf(w, "   %s:\t%v\n", p.Name, p.Value)
				}
				w.Flush()
				fmt.Println("Metrics:")
				for _, m := range m.Metrics {
					fmt.Fprintf(w, "   %s:\t%v\n", m.Name, m.Value)
				}
				w.Flush()
			}
		} else {
			fmt.Fprintln(w, "TrialID\tParamNum\tMetricsNum")
			for _, m := range mr.Models {
				if len(opt.args) > 1 {
					if !strings.HasPrefix(m.WorkerId, opt.args[1]) {
						continue
					}
				}
				fmt.Fprintf(w, "%s\t%d\t%d\n",
					string([]rune(m.WorkerId)[:7]),
					len(m.Parameters),
					len(m.Metrics),
				)
			}
			w.Flush()
		}
	}
}
