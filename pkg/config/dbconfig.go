package config

//Real Time Filtering by device/alertid/or other tags

// InfluxCfg is the main configuration for any InfluxDB TSDB
type InfluxCfg struct {
	ID                 string `xorm:"'id' unique" binding:"Required"`
	Host               string `xorm:"host" binding:"Required"`
	Port               int    `xorm:"port" binding:"Required;IntegerNotZero"`
	DB                 string `xorm:"db" binding:"Required"`
	User               string `xorm:"user" binding:"Required"`
	Password           string `xorm:"password" binding:"Required"`
	Retention          string `xorm:"'retention' default 'autogen'" binding:"Required"`
	Precision          string `xorm:"'precision' default 's'" binding:"Default(s);OmitEmpty;In(h,m,s,ms,u,ns)"` //posible values [h,m,s,ms,u,ns] default seconds for the nature of data
	Timeout            int    `xorm:"'timeout' default 30" binding:"Default(30);IntegerNotZero"`
	UserAgent          string `xorm:"useragent" binding:"Default(vcentercollector)"`
	EnableSSL          bool   `xorm:"enable_ssl"`
	SSLCA              string `xorm:"ssl_ca"`
	SSLCert            string `xorm:"ssl_cert"`
	SSLKey             string `xorm:"ssl_key"`
	InsecureSkipVerify bool   `xorm:"insecure_skip_verify"`
	Description        string `xorm:"description"`
}

//http://www-01.ibm.com/support/docview.wss?uid=nas8N1019111

// VCenterCfg contains all related device definitions
type VCenterCfg struct {
	ID string `xorm:"'id' unique" binding:"Required"`
	//https://+Host+:12443
	Host     string `xorm:"host" binding:"Required"`
	Port     int    `xorm:"port" binding:"Required"`
	User     string `xorm:"user" binding:"Required"`
	Password string `xorm:"password" binding:"Required"`

	Active bool `xorm:"'active' default 1"`

	Freq           int `xorm:"'freq' default 60" binding:"Default(60);IntegerNotZero"`
	UpdateScanFreq int `xorm:"'update_scan_freq' default 60" binding:"Default(60);UIntegerAndLessOne"`

	OutDB    string `xorm:"outdb"`
	LogLevel string `xorm:"loglevel" binding:"Default(info)"`
	APIDebug bool   `xorm:"api_debug"`
	LogFile  string `xorm:"logfile"`

	//influx tags
	DeviceTagName  string   `xorm:"devicetagname" binding:"Default(hostname)"`
	DeviceTagValue string   `xorm:"devicetagvalue" binding:"Default(id)"`
	ExtraTags      []string `xorm:"extra-tags"`

	Description string `xorm:"description"`
}

// DEVICE TABLE

// DeviceCfg contains all  related device definitions
type DeviceCfg struct {
	//LogicalPartition
	ID             string `xorm:"'id' unique" binding:"Required"` //VMID
	Name           string `xorm:"name" binding:"Required"`        //
	SerialNumber   string `xorm:"serial_number"`                  //
	OSVersion      string `xorm:"os_version"`                     //
	Type           string `xorm:"type"`                           //
	PartitionState string `xorm:"-"`

	Location string `xorm:"location"`

	ExtraTags []string `xorm:"extra-tags"` //common tags for devices stats

	Description string `xorm:"description"`
}

// DBConfig read from DB
type DBConfig struct {
	Influxdb map[string]*InfluxCfg
	VCenter  map[string]*VCenterCfg
	Devices  map[string]*DeviceCfg
}

// Init initialices the DB
func Init(cfg *DBConfig) error {

	log.Debug("--------------------Initializing Config-------------------")

	log.Debug("-----------------------END Config metrics----------------------")
	return nil
}
