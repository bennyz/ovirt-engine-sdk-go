//
// Copyright (c) 2020 huihui <huihui.fu@cs2c.com.cn>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package examples

import (
	"fmt"
	"time"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func printVmMacDup() {
	inputRawURL := "https://10.1.111.229/ovirt-engine/api"

	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(inputRawURL).
		Username("admin@internal").
		Password("qwer1234").
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		fmt.Printf("Make connection failed, reason: %v\n", err)
		return
	}
	defer conn.Close()

	// To use `Must` methods, you should recover it if panics
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panics occurs, try the non-Must methods to find the reason")
		}
	}()

	// Get the reference to the vm service:
	vmsService := conn.SystemService().VmsService()

	vmsResp, err := vmsService.List().Send()
	if err != nil {
		fmt.Printf("Failed to get vm list, reason: %v\n", err)
		return
	}

	vmNics := make(map[string]string, 0)
	// Iterate via all virtual machines and print if they have duplicated MAC
	// address with any other virtual machine in the system:
	for _, vm := range vmsResp.MustVms().Slice() {
		fmt.Printf("Vm name: %v\n", vm.MustName())
		vmService := vmsService.VmService(vm.MustId())
		nicsResp, err := vmService.NicsService().List().Send()
		if err != nil {
			fmt.Printf("Failed to get nic list, reason: %v\n", err)
		}
		for _, nic := range nicsResp.MustNics().Slice() {
			if _, exists := vmNics[nic.MustMac().MustAddress()]; exists {
				println("key exists in map")
				fmt.Printf("[%v]: MAC address '%v' is used by following virtual machine already: %v",
					nic.MustName(), nic.MustMac().MustAddress(), vmNics)
			} else {
				vmNics[nic.MustMac().MustAddress()] = vm.MustName()
			}

		}
	}
}
