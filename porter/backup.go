package network

import (
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	yrutils "github.com/riboseyim/go-chilopod/utils"
	"github.com/vjeantet/jodaTime"
	"go.opencensus.io/trace"
)

type pingResult struct {
	IP          string        `json:"IP"`
	PacketsSent int           `json:"PacketsSent"`
	PacketsRecv int           `json:"PacketsRecv"`
	PacketLoss  float64       `json:"PacketLoss"`
	MinRtt      time.Duration `json:"MinRtt"`
	MaxRtt      time.Duration `json:"MaxRtt"`
	TaskNo      string        `json:"TaskNo"`
}

func Exec_backup_task(ctx context.Context, cfgfile string) {
	ctx, handleSpan := trace.StartSpan(ctx, "Exec_backup_task")
	defer handleSpan.End()
	//go_ping("riboseyim.com", 1)
	//go_ping("riboseyim.com", 3)

	var taskno = jodaTime.Format("YYYYMMddhhmss", time.Now())

	log.Println("context create taskno:%s", taskno)

	ctx = context.WithValue(ctx, "taskno", taskno)

	rows, _ := yrutils.ReadLine("backup.cfg")

	chan_len := len(rows)
	ch := make(chan pingResult, chan_len)

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		log.Println("-----load row:%s", row)
		go_backup(ctx, ch, row, "./tmp")
	}

	<-ch
}

func go_backup(ctx context.Context, chan_pingResult chan pingResult, srcfile string, dstpath string) {

	ctx, handleSpan := trace.StartSpan(ctx, "go_backup")
	defer handleSpan.End()

	TaskNo := ctx.Value("taskno").(string)
	log.Println("go_backup() TaskNo:%s", TaskNo)

	dstfilename := strings.Replace(srcfile, "/", "_", 10)
	dstfile := dstpath + "/" + dstfilename
	log.Println("dstfile:%s", dstfile)

	cmd := "cp " + srcfile + " " + dstfile
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Printf("Failed to execute command: %s", cmd)
	}
	log.Println("Exec:" + string(out))
}
