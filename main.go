package main

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

type BluetoothResult struct {
	MAC            string
	AdvertisedName string
	RSSI           int16
	Time           time.Time
}

var resultMap map[string][]BluetoothResult

func printUniqNames() {
	for MAC, resultList := range resultMap {
		fmt.Printf("%v: \n", MAC)
		visited := make(map[string]bool)
		for _, item := range resultList {
			name := item.AdvertisedName
			if name != "" && !visited[name] {
				fmt.Printf("    %v\n", item.AdvertisedName)
				visited[name] = true
			}
		}
	}
}

func main() {
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	resultMap = make(map[string][]BluetoothResult)

	// Start scanning.
	println("scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		result := BluetoothResult{
			MAC:            device.Address.String(),
			AdvertisedName: device.LocalName(),
			RSSI:           device.RSSI,
			Time:           time.Now(),
		}
		if resultMap[device.Address.String()] == nil {
			fmt.Printf("%v: %v - %v\n", device.Address.String(), device.LocalName(), device.RSSI)
		}
		resultMap[device.Address.String()] = append(resultMap[device.Address.String()], result)
		// a, _ := adapter.Connect(bluetooth.Address{}, bluetooth.ConnectionParams{})
		// a.DiscoverServices()
		// printUniqNames()
	})
	must("start scan", err)
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
