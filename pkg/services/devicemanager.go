package services

import (
	"strings"
	"time"
	"fmt"
)

type MyQService interface {
	SetDoorState(string, string) error
	DeviceState(serialNumber string) (string, error)
}
type DeviceManager struct {
	options DeviceManagerOptions
	myqService MyQService
}

type DeviceManagerOptions struct {
	MaxRetries int
	RetryInterval int
}

type DeviceDesiredState struct {
	SerialNumber string
	DesiredState string
	Action string
}

func NewDeviceManager(options DeviceManagerOptions, service MyQService) *DeviceManager {
	if options.MaxRetries <= 0 {
		options.MaxRetries = 3
	}
	if options.RetryInterval <= 0 {
		options.RetryInterval = 30
	}
	return &DeviceManager{
		options: options,
		myqService: service,
	}
}

func (dm DeviceManager) SetDesiredState(desiredState DeviceDesiredState) error {
	retries := dm.options.MaxRetries
	retry := 0
	state, err := dm.myqService.DeviceState(desiredState.SerialNumber)
	
	for !strings.Contains(state, desiredState.DesiredState) && retry < retries {
		retry++
		err = dm.myqService.SetDoorState(desiredState.SerialNumber, desiredState.Action)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		time.Sleep(time.Duration(dm.options.RetryInterval) * time.Second)
		state, err = dm.myqService.DeviceState(desiredState.SerialNumber)
		fmt.Printf("State: %s, Error: %s\n", state, err)
	}

	if err != nil {
		return err
	}

	if !strings.Contains(state, desiredState.DesiredState) {
		return fmt.Errorf("device not %s", desiredState.DesiredState)
	}

	return nil

}



