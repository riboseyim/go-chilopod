package network

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	yrutils "github.com/riboseyim/go-chilopod/utils"
	ping "github.com/sparrc/go-ping"
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

func Exec_ping_task(ctx context.Context, cfg string) {
	ctx, handleSpan := trace.StartSpan(ctx, "exec_ping_task")
	defer handleSpan.End()
	//go_ping("riboseyim.com", 1)
	//go_ping("riboseyim.com", 3)

	var taskno = jodaTime.Format("YYYYMMddhhmss", time.Now())

	log.Println("context create taskno:%s", taskno)

	ctx = context.WithValue(ctx, "taskno", taskno)

	rows, _ := yrutils.ReadLine("host.cfg")

	chan_len := len(rows)
	ch := make(chan pingResult, chan_len)

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		log.Println("-----load row:%s", row)
		go go_ping(ctx, ch, row, 2, time.Second*5, false)
	}

	//print_result(ch, chan_len)
	save_result(ch, chan_len)

}

func save_result(ch chan pingResult, chan_len int) {
	var result pingResult
	var data [][]string = make([][]string, 100, 100)
	data[0] = []string{"任务号", "IP", "Ping丢包率"}
	for i := 0; i < chan_len; i++ {
		result = <-ch
		log.Println("taskno:%s host:%s PacketLoss:%s", result.TaskNo, result.IP, result.PacketLoss)
		data[i+1] = []string{result.TaskNo, result.IP, strconv.FormatFloat(result.PacketLoss, 'f', 0, 32)}
	}
	log.Println("-----start write report-----")
	filename := "ping.csv"
	yrutils.SaveFile(filename, data, false)

}

func print_result(ch chan pingResult, chan_len int) {
	var result pingResult
	for i := 0; i < chan_len; i++ {
		result = <-ch
		log.Println("taskno:%s host:%s PacketLoss:%s", result.TaskNo, result.IP, result.PacketLoss)
	}
}

func print_result_daemon(ch chan pingResult) {
	var result pingResult
	for {
		result = <-ch
		log.Println("chan pingResult host:%s,PacketLoss:%s", result.IP, result.PacketLoss)
	}
}

func go_ping(ctx context.Context, chan_pingResult chan pingResult, host string, pingc int, Timeout time.Duration, debug bool) {

	ctx, handleSpan := trace.StartSpan(ctx, "go_ping")
	defer handleSpan.End()

	TaskNo := ctx.Value("taskno").(string)

	pinger, err := ping.NewPinger(host)
	if err != nil {
		panic(err)
	}
	pinger.Count = pingc
	//	pinger.Interval = time.Second * 3
	//pinger.Debug = true
	pinger.Timeout = Timeout

	pinger.OnRecv = func(pkt *ping.Packet) {
		if debug {
			log.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		if debug {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
			fmt.Printf("--- %s ping statistics end ---\n\n", stats.Addr)
		}

		var result = pingResult{
			IP:          host,
			PacketsSent: stats.PacketsSent,
			PacketsRecv: stats.PacketsRecv,
			PacketLoss:  stats.PacketLoss,
			MinRtt:      stats.MinRtt,
			MaxRtt:      stats.MaxRtt,
			TaskNo:      TaskNo,
		}
		//log.Println("go_ping() before ctx host:%s,PacketLoss:%s", host, result.PacketLoss)
		chan_pingResult <- result
	}

	pinger.Run() // blocks until finished

	if debug {
		stats := pinger.Statistics() // get send/receive/rtt stats
		//	log.Println(stats)
		log.Println("PacketsSent", stats.PacketsSent)
		log.Println("PacketsRecv", stats.PacketsRecv)
		log.Println("PacketLoss", stats.PacketLoss)
		log.Println("MinRtt", stats.MinRtt)
		log.Println("MaxRtt", stats.MaxRtt)
	}

}

func go_ping_simple() {
	pinger, err := ping.NewPinger("www.baidu.com")
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	log.Println(stats)
}

func go_ping_daemon() {
	pinger, err := ping.NewPinger("www.baidu.com")
	if err != nil {
		panic(err)
	}

	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}
