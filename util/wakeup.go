package util

import (
	"encoding/json"

	"github.com/go-ping/ping"
	"github.com/mdlayher/wol"
	"github.com/r2binx/heimboard-wakeup-go/config"
	"github.com/r2binx/heimboard-wakeup-go/schedule"

	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Wakeup struct {
	iface  *net.Interface
	config config.Config
}

func NewWakeup(config config.Config) *Wakeup {
	iface, err := net.InterfaceByName(config.Iface)
	if err != nil {
		log.Fatal("Failed getting interface", err)
	}
	return &Wakeup{
		config: config,
		iface:  iface,
	}
}

func (w *Wakeup) GetSchedule() (schedule schedule.Schedule) {
	if _, err := os.Stat("schedule.json"); os.IsNotExist(err) {
		log.Println("No schedule found, creating new one")
		w.SetSchedule(schedule)
	}

	file, err := ioutil.ReadFile("schedule.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(file, &schedule)

	return schedule
}

func (w *Wakeup) SetSchedule(schedule schedule.Schedule) error {
	file, err := json.Marshal(schedule)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ioutil.WriteFile("schedule.json", file, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (w *Wakeup) Wake(mac net.HardwareAddr) (bool, error) {
	c, err := wol.NewRawClient(w.iface)
	if err != nil {
		return false, err
	}
	defer c.Close()

	hostOnline := pingHost(w.config.HostIp)
	if !hostOnline {
		log.Println("Waking up:", mac)
		err := c.Wake(mac)
		if err != nil {
			return false, err
		}

		for i := 0; !hostOnline && i < 12; i++ {
			log.Println("Host is offline, waiting for it to come online")
			hostOnline = pingHost(w.config.HostIp)
			time.Sleep(time.Second * 5)
		}

	}
	return hostOnline, nil
}

func pingHost(host net.IP) bool {
	var status bool
	pinger, err := ping.NewPinger(host.String())
	if err != nil {
		log.Fatal("Error creating pinger:", err)
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.Run()

	status = pinger.Statistics().PacketsRecv > 0

	return status
}

func parseMsTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

func (w *Wakeup) CheckSchedule() {
	sched := w.GetSchedule()
	if sched.Time != 0 {

		scheduledTimestamp := parseMsTimestamp(sched.Time)
		now := time.Now()
		scheduledTime := time.Date(now.Year(), now.Month(), now.Day(), scheduledTimestamp.Hour(), scheduledTimestamp.Minute(), 0, 0, time.Local)

		if now.After(scheduledTime) && now.Before(scheduledTime.Add(time.Minute*5)) {
			hostOnline := pingHost(w.config.HostIp)
			if sched.Action == "boot" && !hostOnline {
				log.Println("Host is offline, booting")
				w.Wake(w.config.WolMac)
			} else if sched.Action != "boot" {
				log.Println("Unknown action:", sched.Action)
			}
		}
	}
}
