package network

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	yrutils "github.com/riboseyim/go-chilopod/utils"
	g "github.com/soniah/gosnmp"
	"github.com/vjeantet/jodaTime"
	"go.opencensus.io/trace"
)

type snmpResult struct {
	IP          string        `json:"IP"`
	PacketsSent int           `json:"PacketsSent"`
	PacketsRecv int           `json:"PacketsRecv"`
	PacketLoss  float64       `json:"PacketLoss"`
	MinRtt      time.Duration `json:"MinRtt"`
	MaxRtt      time.Duration `json:"MaxRtt"`
	TaskNo      string        `json:"TaskNo"`
}

func Exec_snmp_task(ctx context.Context, cfg string) {
	ctx, handleSpan := trace.StartSpan(ctx, "exec_snmp_task")
	defer handleSpan.End()
	//go_snmp("riboseyim.com", 1)
	//go_snmp("riboseyim.com", 3)

	var taskno = jodaTime.Format("YYYYMMddhhmss", time.Now())

	log.Println("context create taskno:%s", taskno)

	ctx = context.WithValue(ctx, "taskno", taskno)

	rows, _ := yrutils.ReadLine("host.cfg")

	chan_len := len(rows)
	ch := make(chan snmpResult, chan_len)

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		log.Println("-----load row:%s", row)
		go Go_snmp(ctx, ch, row, "public", time.Second*5, false)
	}

	print_snmp_result(ch, chan_len)
	//	save_snmp_result(ch, chan_len)

}

func save_snmp_result(ch chan snmpResult, chan_len int) {
	var result snmpResult
	var data [][]string = make([][]string, 100, 100)
	data[0] = []string{"任务号", "IP", "snmp丢包率"}
	for i := 0; i < chan_len; i++ {
		result = <-ch
		log.Println("taskno:%s host:%s PacketLoss:%s", result.TaskNo, result.IP, result.PacketLoss)
		data[i+1] = []string{result.TaskNo, result.IP, strconv.FormatFloat(result.PacketLoss, 'f', 0, 32)}
	}
	log.Println("-----start write report-----")
	filename := "snmp.csv"
	yrutils.SaveFile(filename, data, false)

}

func print_snmp_result(ch chan snmpResult, chan_len int) {
	var result snmpResult
	for i := 0; i < chan_len; i++ {
		result = <-ch
		log.Println("taskno:%s host:%s PacketLoss:%s", result.TaskNo, result.IP, result.PacketLoss)
	}
}

func Go_snmp(ctx context.Context, chan_snmpResult chan snmpResult, host string, snmpc string, Timeout time.Duration, debug bool) {

	ctx, handleSpan := trace.StartSpan(ctx, "go_snmp")
	defer handleSpan.End()

	TaskNo := ctx.Value("taskno").(string)

	envPort := "161" //default snmp port

	port, _ := strconv.ParseUint(envPort, 10, 16)

	// Build our own GoSNMP struct, rather than using g.Default.
	// Do verbose logging of packets.
	params := &g.GoSNMP{
		Target:    host,
		Port:      uint16(port),
		Community: snmpc,
		Version:   g.Version2c,
		Timeout:   Timeout,
		Logger:    log.New(os.Stdout, "", 0),
	}
	err := params.Connect()
	if err != nil {
		log.Println("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	// $ snmpwalk -v 2c -c public 192.168.213.128 1.3.6.1.2.1.1.2
	// SNMPv2-MIB::sysObjectID.0 = OID: NET-SNMP-MIB::netSnmpAgentOIDs.10
	// $ snmpwalk -v 2c -c public 192.168.213.128 1.3.6.1.2.1.1.5
	// SNMPv2-MIB::sysName.0 = STRING: NW-DD-APP

	oids := []string{"1.3.6.1.2.1.1.2", "1.3.6.1.2.1.1.5"}
	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Println("Get() err: %v", err2)
	}

	if result != nil {
		for i, variable := range result.Variables {
			log.Println("----range variable: %d: oid: %s ", i, variable.Name)
			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {
			case g.OctetString:
				log.Println("case g.OctetString string: %s\n", string(variable.Value.([]byte)))
			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				log.Println("variable value:%s default number: %d", variable.Value, g.ToBigInt(variable.Value))
			}
		}
	}

	var snmp_result = snmpResult{
		IP: host,
		// PacketsSent: stats.PacketsSent,
		// PacketsRecv: stats.PacketsRecv,
		// PacketLoss:  stats.PacketLoss,
		// MinRtt:      stats.MinRtt,
		// MaxRtt:      stats.MaxRtt,
		TaskNo: TaskNo,
	}
	//log.Println("go_snmp() before ctx host:%s,PacketLoss:%s", host, result.PacketLoss)
	chan_snmpResult <- snmp_result

}
