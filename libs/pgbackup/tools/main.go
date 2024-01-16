package main

import (
	"log"

	"github.com/Dcarbon/go-shared/libs/pgbackup"
)

type MoveConfig struct {
	Label string
	Src   string
	Dst   string
	Setup []string
}

func (mc *MoveConfig) Execute() error {
	var dump = pgbackup.NewDump("./static")

	var dumpFile, err = dump.Dump(mc.Src)
	if nil != err {
		return err
	}
	log.Printf("[%s] Dump success\n", mc.Label)

	config, err := pgbackup.NewConfigFromUrl(mc.Dst)
	if nil != err {
		return err
	}

	log.Printf("[%s] Setup \n", mc.Label)
	err = config.CreateDb()
	if nil != err {
		return err
	}
	for _, sqlCmd := range mc.Setup {
		_, err = config.Execute(sqlCmd)
		if nil != err {
			return err
		}
	}

	log.Printf("[%s] Setup success\n", mc.Label)

	log.Printf("[%s] Restore \n", mc.Label)
	var restore = pgbackup.Restore{}
	err = restore.Execute(dumpFile, mc.Dst)
	if nil != err {
		return err
	}

	return nil
}

var moveList = []*MoveConfig{
	{
		Label: "IotInfo",
		Src:   "postgres://admin:244466666@10.60.0.58/iott?table=iots",
		Dst:   "postgres://admin:hellosecret@13.228.11.143:5432/iot_info",
		Setup: []string{
			"CREATE EXTENSION IF NOT EXISTS postgis",
		},
	},
	{
		Label: "Sensors",
		Src:   "postgres://admin:244466666@10.60.0.58/iott?table=sensors",
		Dst:   "postgres://admin:hellosecret@13.228.11.143:5432/sensors",
		Setup: []string{
			"CREATE EXTENSION IF NOT EXISTS postgis",
		},
	},
	{
		Label: "Projects",
		Src:   "postgres://admin:244466666@10.60.0.58/iott?table=projects",
		Dst:   "postgres://admin:hellosecret@13.228.11.143:5432/project",
		Setup: []string{
			"CREATE EXTENSION IF NOT EXISTS postgis",
		},
	},
}

func main() {
	for _, it := range moveList {
		err := it.Execute()
		if nil != err {
			log.Println("Move error: ", it.Label, err)
		}

	}
}
