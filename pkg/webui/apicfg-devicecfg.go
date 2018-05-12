package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/vCenterCollector/pkg/agent"
	"github.com/toni-moreno/vCenterCollector/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgDevice DeviceCfg API REST creator
func NewAPICfgDevice(m *macaron.Macaron) error {

	bind := binding.Bind

	m.Group("/api/cfg/devices", func() {
		m.Get("/", reqSignedIn, GetDeviceCfg)
		m.Post("/", reqSignedIn, bind(config.DeviceCfg{}), AddDeviceCfg)
		m.Put("/:id", reqSignedIn, bind(config.DeviceCfg{}), UpdateDeviceCfg)
		m.Delete("/:id", reqSignedIn, DeleteDeviceCfg)
		m.Get("/:id", reqSignedIn, GetDeviceCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetDeviceCfgAffectOnDel)
	})

	return nil
}

// GetDeviceCfg Return Server Array
func GetDeviceCfg(ctx *Context) {
	cfgarray, err := agent.MainConfig.Database.GetDeviceCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Device :%+s", err)
		return
	}
	ctx.JSON(200, &cfgarray)
	log.Debugf("Getting DEVICEs %+v", &cfgarray)
}

// AddDeviceCfg Insert new measurement groups to de internal BBDD --pending--
func AddDeviceCfg(ctx *Context, dev config.DeviceCfg) {
	log.Printf("ADDING Device %+v", dev)
	affected, err := agent.MainConfig.Database.AddDeviceCfg(dev)
	if err != nil {
		log.Warningf("Error on insert new Backend %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateDeviceCfg --pending--
func UpdateDeviceCfg(ctx *Context, dev config.DeviceCfg) {
	id := ctx.Params(":id")
	log.Debugf("Tying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateDeviceCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update device %s  , affected : %+v , error: %s", dev.ID, affected, err)
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteDeviceCfg --pending--
func DeleteDeviceCfg(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Tying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelDeviceCfg(id)
	if err != nil {
		log.Warningf("Error on delete influx db %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetDeviceCfgByID --pending--
func GetDeviceCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetDeviceCfgByID(id)
	if err != nil {
		log.Warningf("Error on get device db data for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

// GetDeviceCfgAffectOnDel --pending--
func GetDeviceCfgAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetDeviceCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for influx device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
