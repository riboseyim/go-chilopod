package network

type hostCfg struct {
	hostname  string
	logintype string
	username  string
	password  string
	snmpv     string
	snmprcom  string
}

type deviceCfg struct {
	nodename        string
	devicename      string
	devicecname     string
	loopaddress     string
	devicemodelcode string
	logintype       string
	username        string
	password        string
	snmpv           string
	snmprcom        string
	snmpwcom        string
}
