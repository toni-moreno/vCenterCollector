package config

import "fmt"

/***************************
	VCenter backends
	-GetVCenterCfgCfgByID(struct)
	-GetVCenterCfgMap (map - for interna config use
	-GetVCenterCfgArray(Array - for web ui use )
	-AddVCenterCfg
	-DelVCenterCfg
	-UpdateVCenterCfg
  -GetVCenterCfgAffectOnDel
***********************************/

/*GetVCenterCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetVCenterCfgByID(id string) (VCenterCfg, error) {
	cfgarray, err := dbc.GetVCenterCfgArray("id='" + id + "'")
	if err != nil {
		return VCenterCfg{}, err
	}
	if len(cfgarray) > 1 {
		return VCenterCfg{}, fmt.Errorf("Error %d results on get VCenterCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return VCenterCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the VCenter config table", id)
	}
	return *cfgarray[0], nil
}

/*GetVCenterCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetVCenterCfgMap(filter string) (map[string]*VCenterCfg, error) {
	cfgarray, err := dbc.GetVCenterCfgArray(filter)
	cfgmap := make(map[string]*VCenterCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetVCenterCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetVCenterCfgArray(filter string) ([]*VCenterCfg, error) {
	var err error
	var devices []*VCenterCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get VCenterCfg  data filteter with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&devices); err != nil {
			log.Warnf("Fail to get VCenterCfg   data: %v\n", err)
			return nil, err
		}
	}
	return devices, nil
}

/*AddVCenterCfg for adding new devices*/
func (dbc *DatabaseCfg) AddVCenterCfg(dev VCenterCfg) (int64, error) {
	var err error
	var affected int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	//no other relation
	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Added new VCenter backend Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelVCenterCfg for deleting VCenter databases from ID*/
func (dbc *DatabaseCfg) DelVCenterCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()
	// deleting references in VCenterCfg

	affected, err = session.Where("id='" + id + "'").Delete(&VCenterCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully VCenter db with ID %s [ %d Devices Affected  ]", id, affecteddev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateVCenterCfg for adding new VCenter*/
func (dbc *DatabaseCfg) UpdateVCenterCfg(id string, dev VCenterCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + id + "'").UseBool().AllCols().Update(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	err = session.Commit()
	if err != nil {
		return 0, err
	}

	log.Infof("Updated VCenter Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetVCenterCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetVCenterCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	//	var devices []*VCenterCfg
	var obj []*DbObjAction
	/*
		for _, val := range devices {
			obj = append(obj, &DbObjAction{
				Type:     "VCenterCfg",
				TypeDesc: "VCenter Devices",
				ObID:     val.ID,
				Action:   "Reset VCenter Server fro 'default' InfluxDB Server",
			})

		}*/
	return obj, nil
}
