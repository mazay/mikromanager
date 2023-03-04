package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func fetchResources(cfg *PollerCFG) error {
	resource, err := cfg.Client.Run("/system/resource/print")
	if err != nil {
		return err
	}
	logger.Debugf("fetched resource data for %s", cfg.Client.Address)
	inrec, _ := json.Marshal(resource[0].Map)
	return json.Unmarshal(inrec, &cfg.Device)
}

func fetchIdentity(cfg *PollerCFG) error {
	identity, err := cfg.Client.Run("/system/identity/print")
	if err != nil {
		return err
	}
	if len(identity) > 0 {
		logger.Debugf("identity for %s is %s", cfg.Client.Address, identity[0].Map["name"])
		cfg.Device.Identity = string(identity[0].Map["name"])
		return nil
	}
	return fmt.Errorf("got an empty identity data")
}

func fetchManagementIp(cfg *PollerCFG) error {
	ipaddr, err := cfg.Client.Run("/ip/address/print ?comment=MGMT")
	if err != nil {
		return err
	}
	if len(ipaddr) > 0 {
		addr := strings.Split(ipaddr[0].Map["address"], "/")[0]
		logger.Debugf("management address for device ID %s will be set to %s", cfg.Device.Id, addr)
		cfg.Device.Address = addr
		return nil
	}
	return fmt.Errorf("couldn't find management IP-address")
}

func writeBackupFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		logger.Fatal(err)
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	return nil
}

func timeSliceBy(start time.Time, end time.Time, multiplier time.Duration) []time.Time {
	var times []time.Time
	for d := start; !d.After(end); d = d.Add(multiplier) {
		times = append(times, d)
	}
	return times
}
