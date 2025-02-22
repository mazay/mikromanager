package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

func fetchResources(cfg *PollerCFG) error {
	resource, err := cfg.Client.Run("/system/resource/print")
	if err != nil {
		return err
	}
	logger.Debug("fetched resource data", zap.String("address", cfg.Client.Address))
	inrec, _ := json.Marshal(resource[0].Map)
	return json.Unmarshal(inrec, &cfg.Device)
}

func fetchRbDetails(cfg *PollerCFG) error {
	resource, err := cfg.Client.Run("/system/routerboard/print")
	if err != nil {
		return err
	}
	logger.Debug("fetched rouerboard data", zap.String("address", cfg.Client.Address))
	inrec, _ := json.Marshal(resource[0].Map)
	return json.Unmarshal(inrec, &cfg.Device)
}

func fetchIdentity(cfg *PollerCFG) error {
	identity, err := cfg.Client.Run("/system/identity/print")
	if err != nil {
		return err
	}
	if len(identity) > 0 {
		logger.Debug("identity", zap.String("address", cfg.Client.Address), zap.String("name", identity[0].Map["name"]))
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
		logger.Info("management address discovered", zap.String("device", cfg.Device.Id), zap.String("address", addr))
		cfg.Device.Address = addr
		return nil
	}
	return fmt.Errorf("couldn't find management IP-address")
}

func timeSliceBy(start time.Time, end time.Time, multiplier time.Duration) []time.Time {
	var times []time.Time
	for d := start; !d.After(end); d = d.Add(multiplier) {
		times = append(times, d)
	}
	return times
}
