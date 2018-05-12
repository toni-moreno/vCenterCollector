package webui

import (
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/vCenterCollector/pkg/agent"
	"github.com/toni-moreno/vCenterCollector/pkg/agent/devices/vcenter"
	"github.com/toni-moreno/vCenterCollector/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgVCenterServer VCenterServer API REST creator
func NewAPICfgVCenterServer(m *macaron.Macaron) error {

	bind := binding.Bind

	m.Group("/api/cfg/vcenterserver", func() {
		m.Get("/", reqSignedIn, GetVCenterServer)
		m.Post("/", reqSignedIn, bind(config.VCenterCfg{}), AddVCenterServer)
		m.Put("/:id", reqSignedIn, bind(config.VCenterCfg{}), UpdateVCenterServer)
		m.Delete("/:id", reqSignedIn, DeleteVCenterServer)
		m.Get("/:id", reqSignedIn, GetVCenterServerByID)
		m.Get("/checkondel/:id", reqSignedIn, GetVCenterServerAffectOnDel)
		m.Post("/ping/", reqSignedIn, bind(config.VCenterCfg{}), PingVCenterServer)
		m.Post("/import", reqSignedIn, bind(config.VCenterCfg{}), ImportVCenterDevices)
	})

	return nil
}

// GetVCenterServer Return Server Array
func GetVCenterServer(ctx *Context) {
	cfgarray, err := agent.MainConfig.Database.GetVCenterCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get VCenterServer db :%+s", err)
		return
	}
	ctx.JSON(200, &cfgarray)
	log.Debugf("Getting DEVICEs %+v", &cfgarray)
}

// AddVCenterServer Insert new measurement groups to de internal BBDD --pending--
func AddVCenterServer(ctx *Context, dev config.VCenterCfg) {
	log.Printf("ADDING VCenterServer Backend %+v", dev)
	affected, err := agent.MainConfig.Database.AddVCenterCfg(dev)
	if err != nil {
		log.Warningf("Error on insert new Backend %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateVCenterServer --pending--
func UpdateVCenterServer(ctx *Context, dev config.VCenterCfg) {
	id := ctx.Params(":id")
	log.Debugf("Tying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateVCenterCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update VCenterServer db %s  , affected : %+v , error: %s", dev.ID, affected, err)
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteVCenterServer --pending--
func DeleteVCenterServer(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Tying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelVCenterCfg(id)
	if err != nil {
		log.Warningf("Error on delete influx db %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetVCenterServerByID --pending--
func GetVCenterServerByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetVCenterCfgByID(id)
	if err != nil {
		log.Warningf("Error on get VCenterServer db data for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetVCenterServerAffectOnDel --pending--
func GetVCenterServerAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetVCenterCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for influx device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}

//PingVCenterServer Return ping result
func PingVCenterServer(ctx *Context, cfg config.VCenterCfg) {
	log.Infof("trying to ping influx server %s : %+v", cfg.ID, cfg)
	_, elapsed, message, err := vcenter.Ping(&cfg, log, false, "")
	type result struct {
		Result  string
		Elapsed time.Duration
		Message string
	}
	if err != nil {
		log.Debugf("ERROR on ping VCenterServerDB Server : %s", err)
		res := result{Result: "NOOK", Elapsed: elapsed, Message: err.Error()}
		ctx.JSON(400, res)
	} else {
		log.Debugf("OK on ping VCenterServerDB Server %+v, %+v", elapsed, message)
		res := result{Result: "OK", Elapsed: elapsed, Message: message}
		ctx.JSON(200, res)
	}
}

// ImportVCenterDevices new snmpdevice to de internal BBDD --pending--
func ImportVCenterDevices(ctx *Context, dev config.VCenterCfg) {
	log.Warningf("Importing VCenter devices for VCenter: %s", dev.ID)
	ses, _, _, err := vcenter.Ping(&dev, log, false, "")
	if err != nil {
		log.Errorf("Error on Ping VCenter Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}

	devices, err := vcenter.ScanVCenter(ses)
	if err != nil {
		log.Errorf("Error on Scan VCenter Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}

	for smid, sm := range devices {
		d := &config.DeviceCfg{
			ID:        sm.ID,
			Name:      sm.Name,
			OSVersion: "osversion",
			Type:      "Managed",
			Location:  dev.ID,
			//Nmon properties doesn't apply here
		}

		log.Infof("Importing Managed System: %s | %s", smid, sm.Name)

		agent.MainConfig.Database.AddOrUpdateDeviceCfg(d)

	}

	ctx.JSON(200, &devices)
}
