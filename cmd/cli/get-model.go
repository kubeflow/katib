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

//NewCommandGetTrials generate run cmd
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
	var soverviews []*api.StudyOverview
	c := api.NewManagerClient(conn)
	if len(opt.args) > 0 {
		req := &api.GetStudiesRequest{}
		r, err := c.GetStudies(context.Background(), req)
		if err != nil {
			log.Fatalf("GetStudy failed: %v", err)
			return
		}
		if len(r.StudyInfos) > 0 {
			var sInfo []*api.StudyInfo
			for _, si := range r.StudyInfos {
				if len(opt.args) > 0 {
					if utf8.RuneCountInString(opt.args[0]) >= 7 {
						if strings.HasPrefix(si.StudyId, opt.args[0]) {
							soverviews = append(soverviews, &api.StudyOverview{
								Name:  si.Name,
								Owner: si.Owner,
							})
							break
						}
					}
					if si.Name == opt.args[0] {
						soverviews = append(soverviews, &api.StudyOverview{
							Name:  si.Name,
							Owner: si.Owner,
						})
						break
					}
				} else {
					soverviews = append(soverviews, &api.StudyOverview{
						Name:  si.Name,
						Owner: si.Owner,
					})
				}
			}
			if len(sInfo) == 0 {
				log.Fatalf("No Study %v is not saved.", opt.args[0])
				return
			}
		}
	}
	if len(soverviews) == 0 {
		sreq := &api.GetSavedStudiesRequest{}
		sr, err := c.GetSavedStudies(context.Background(), sreq)
		if err != nil {
			log.Fatalf("GetStudy failed: %v", err)
			return
		}
		if len(sr.Studies) == 0 {
			log.Fatalf("No Studies are saved.")
			return
		}
		for _, s := range sr.Studies {
			soverviews = append(soverviews, s)
		}
	}
	for _, si := range soverviews {
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
					if !strings.HasPrefix(m.TrialId, opt.args[1]) {
						continue
					}
				}
				fmt.Printf("TrialID :%v\n", m.TrialId)
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
					if !strings.HasPrefix(m.TrialId, opt.args[1]) {
						continue
					}
				}
				fmt.Fprintf(w, "%s\t%d\t%d\n",
					string([]rune(m.TrialId)[:7]),
					len(m.Parameters),
					len(m.Metrics),
				)
			}
			w.Flush()
		}
	}
}
