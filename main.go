package main

import (
	"context"
	"flag"
	"log"

	openzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/http"
	chilopod_net "github.com/riboseyim/go-chilopod/network"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

const (
	MODULE  = "Chilopod Module"
	VERSION = "-1.0-release-201812"
	SERVICE = "127.0.1.1:8080"
	AUTHOR  = "@RiboseYim"
)

func main() {
	//模式：listener|监听捕捉，probe|主动探测
	model := flag.String("m", "", "[listener | probe  ].eg")
	//任务：
	task := flag.String("t", "", "[ping | snmp | ssh | all ].eg")
	//配置文件
	cfgfile := flag.String("cfg", "", "")

	flag.Parse()
	log.Printf("Welcome to [ Chilopod System %s ] \n Author:%s \n\n", VERSION, AUTHOR)
	log.Println("model:", *model)
	log.Println("task:", *task)
	log.Println("cfgfile:", *cfgfile)

	zipkinEndPoint, err := openzipkin.NewEndpoint(MODULE, SERVICE)
	if err != nil {
		log.Println(err)
	}
	reporter := http.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	exporter := zipkin.NewExporter(reporter, zipkinEndPoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	ctx := context.Background()
	ctx, handSpan := trace.StartSpan(ctx, "main")
	defer handSpan.End()

	new_task(ctx, *model, *task, *cfgfile)

	log.Println("--------finish----------")
}

func new_task(ctx context.Context, model string, task string, cfgfile string) {
	ctx, handleSpan := trace.StartSpan(ctx, "new_task")
	defer handleSpan.End()

	if model == "listener" {
		log.Println("建设中的功能 model:[%s]", model)
	} else if model == "probe" {
		if task == "ping" {
			chilopod_net.Exec_ping_task(ctx, cfgfile)
		} else if task == "snmp" {
			chilopod_net.Exec_snmp_task(ctx, cfgfile)
		} else if task == "ssh" {
		} else if task == "all" {
		} else {
			log.Println("cannot found this taskId:[%s]", task)
		}
	}

}
