package vcenter

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
	"time"

	"github.com/toni-moreno/vCenterCollector/pkg/agent/devices"
	"github.com/toni-moreno/vCenterCollector/pkg/agent/output"
	"github.com/toni-moreno/vCenterCollector/pkg/config"
	"github.com/toni-moreno/vCenterCollector/pkg/data/pointarray"
	"github.com/toni-moreno/vCenterCollector/pkg/data/utils"
)

var (
	cfg    *config.DBConfig
	db     *config.DatabaseCfg
	logDir string
)

// SetDBConfig set agent config
func SetDBConfig(c *config.DBConfig, d *config.DatabaseCfg) {
	cfg = c
	db = d
}

// SetLogDir set log dir
func SetLogDir(l string) {
	logDir = l
}

// Server contains all runtime device related device configu ns and state
type Server struct {
	devices.Base
	client  *http.Client
	cfg     *config.VCenterCfg
	Devices map[string]*config.DeviceCfg
}

// Ping check connection to the
func Ping(c *config.VCenterCfg, log *logrus.Logger, apidbg bool, filename string) (*http.Client, time.Duration, string, error) {
	return nil, 0, "error", nil
}

// ScanVCenterDevices scan VCenter
func (d *Server) ScanVCenterDevices() error {
	d.Infof("Scanning  managed systems")

	var err error
	d.Devices, err = ScanVCenter(d.client)
	if err != nil {
		d.Infof("ERROR on get Managed Systems: %s", err)
		return err
	}
	return nil
}

//ScanVCenter scan Device
func ScanVCenter(client *http.Client) (map[string]*config.DeviceCfg, error) {
	return nil, nil
}

// New create and Initialice a device Object
func New(c *config.VCenterCfg) *Server {
	dev := Server{}
	dev.Init(c)
	return &dev
}

// ToJSON return a JSON version of the device data
func (d *Server) ToJSON() ([]byte, error) {
	d.DataLock()
	defer d.DataUnlock()
	result, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		d.Errorf("Error on Get JSON data from device")
		dummy := []byte{}
		return dummy, nil
	}
	return result, err
}

// GetOutSenderFromMap to get info about the sender will use
func (d *Server) GetOutSenderFromMap(influxdb map[string]*output.InfluxDB) (*output.InfluxDB, error) {
	if len(d.cfg.OutDB) == 0 {
		d.Warnf("GetOutSenderFromMap No OutDB configured on the device")
	}
	var ok bool
	name := d.cfg.OutDB
	if d.Influx, ok = influxdb[name]; !ok {
		//we assume there is always a default db
		if d.Influx, ok = influxdb["default"]; !ok {
			//but
			return nil, fmt.Errorf("No influx config for the device: %s", d.cfg.ID)
		}
	}
	d.Debugf("GetOutSenderFromMap: This VCenter server has configured the %s influxdb : %#+v", name, d.Influx)

	return d.Influx, nil
}

func (d *Server) handleMessages(id string, data interface{}) {

}

func (d *Server) setProtocolDebug(debug bool) {

}

// GetVCenterData get data from Device
func (d *Server) GetVCenterData() {

	bpts, _ := d.Influx.BP()
	startStats := time.Now()

	points := pointarray.New(d.GetLogger(), bpts)
	//prepare batchpoint
	err := d.ImportData(points)
	if err != nil {
		d.Errorf("Error in  import VCenter Data from Device %s: ERROR: %s", d.cfg.ID, err)
		return
	}
	points.Flush()
	elapsedStats := time.Since(startStats)

	d.RtStats.SetGatherDuration(startStats, elapsedStats)
	d.RtStats.AddMeasStats(points.MetSent, points.MetError, points.MeasSent, points.MeasError)

	/*************************
	 *
	 * Send data to InfluxDB process
	 *
	 ***************************/

	startInfluxStats := time.Now()
	if bpts != nil {
		d.Influx.Send(bpts)
	} else {
		d.Warnf("Can not send data to the output DB becaouse of batchpoint creation error")
	}
	elapsedInfluxStats := time.Since(startInfluxStats)
	d.RtStats.AddSentDuration(startInfluxStats, elapsedInfluxStats)
}

/*
Init  does the following

- Initialize not set variables to some defaults
- Initialize logfile for this device
- Initialize comunication channels and initial device state
*/
func (d *Server) Init(c *config.VCenterCfg) error {
	if c == nil {
		return fmt.Errorf("Error on initialice device, configuration struct is nil")
	}
	// Set ALL methods IMPORTANT!!! (review if interface could be better here)
	d.Gather = d.GetVCenterData
	d.Scan = d.ScanVCenterDevices
	d.ReleaseClient = d.releaseClient
	d.Reconnect = d.reconnect
	d.SetProtocolDebug = d.setProtocolDebug
	d.HandleMessages = d.handleMessages
	d.CheckDeviceConnectivity = d.checkDeviceConnectivity

	d.cfg = c

	//Init Freq
	d.Freq = d.cfg.Freq
	if d.cfg.Freq == 0 {
		d.Freq = 60
	}

	//Init Logger

	d.Base.Init(d, d.cfg.ID)
	d.InitLog(logDir+"/"+d.cfg.ID+".log", d.cfg.LogLevel)

	d.DeviceActive = d.cfg.Active

	//Init Device Tags

	conerr := d.Reconnect()
	if conerr != nil {
		d.Errorf("First Device connect error: %s", conerr)
		d.DeviceConnected = false
	} else {
		d.DeviceConnected = true
	}

	//Init TagMap

	d.TagMap = make(map[string]string)
	d.TagMap["device"] = d.cfg.ID

	ExtraTags, err := utils.KeyValArrayToMap(d.cfg.ExtraTags)
	if err != nil {
		d.Warnf("Warning on Device  %s Tag gathering: %s", err)
	}
	utils.MapAdd(d.TagMap, ExtraTags)

	// Init stats
	d.InitStats(d.cfg.ID)

	d.SetScanFreq(d.cfg.UpdateScanFreq)

	return nil
}

// ReleaseClient release connections
func (d *Server) releaseClient() {
	/*if d.client != nil {
		d.client.Close()
	}*/
}

// reconnect does HTTP connection  protocol
func (d *Server) reconnect() error {
	var t time.Duration
	var id string
	var err error
	d.Debugf("Trying Reconnect again....")
	/*if d.client != nil {
		d.client.Close()
	}*/
	d.client, t, id, err = Ping(d.cfg, d.GetLogger(), d.cfg.APIDebug, d.cfg.ID)
	if err != nil {
		d.Errorf("Error on Device connection %s", err)
		return err
	}
	d.Infof("Connected to Device  OK : ID: %s : Duration %s ", id, t.String())
	return nil
}

// checkDeviceConnectivity check if Device connection is ok
func (d *Server) checkDeviceConnectivity() bool {
	d.Debugf("Check Device Connectivity: Nothing to do in the Device %s", d.cfg.ID)

	return true
}
