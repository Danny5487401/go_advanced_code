package _3_abstract_factory

import "fmt"

type AbstractFactory interface {
	CreateTelevision() ITelevision
	CreateAirConditioner() IAirConditioner
}

type ITelevision interface {
	Watch()
}

type IAirConditioner interface {
	SetTemperature(int)
}

type HuaWeiFactory struct{}

func (hf *HuaWeiFactory) CreateTelevision() ITelevision {
	return &HuaWeiTV{}
}
func (hf *HuaWeiFactory) CreateAirConditioner() IAirConditioner {
	return &HuaWeiAirConditioner{}
}

type HuaWeiTV struct{}

func (ht *HuaWeiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type HuaWeiAirConditioner struct{}

func (ha *HuaWeiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("HuaWei AirConditioner set temperature to %d ℃\n", temp)
}

type MiFactory struct{}

func (mf *MiFactory) CreateTelevision() ITelevision {
	return &MiTV{}
}
func (mf *MiFactory) CreateAirConditioner() IAirConditioner {
	return &MiAirConditioner{}
}

type MiTV struct{}

func (mt *MiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type MiAirConditioner struct{}

func (ma *MiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("Mi AirConditioner set temperature to %d ℃\n", temp)
}
