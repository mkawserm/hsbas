package hsbas

import (
	"bytes"
	"hash/fnv"
	"net"
	"time"
)

type addressScheme struct {
	currentTimeStamp uint64
	sequenceNumber   uint64
	machineNumber    uint64
}

func NewAddressScheme() *addressScheme {
	n := &addressScheme{}
	n.machineNumber = n.GetNodeNumber()
	n.currentTimeStamp = 0
	n.sequenceNumber = 0
	return n
}

func (as *addressScheme) GetCurrentUnixTimestampAsMilliSeconds() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func (as *addressScheme) GetCustomTimestampAsMilliSeconds() uint64 {
	return as.GetCurrentUnixTimestampAsMilliSeconds() - 1546300800000
}

func (as *addressScheme) GetStringHash(data string) uint64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(data))
	return uint64(h.Sum32())
}

func (as *addressScheme) GetMacAddress() (address string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				address = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

func (as *addressScheme) GetNodeNumber() uint64 {
	machineId, err := machineID()
	if err != nil {
		machineId = as.GetMacAddress()
	}

	if len(machineId) == 0 {
		machineId = randStringBytes(32)
	}

	return as.GetStringHash(machineId) % 1024
}

func (as *addressScheme) GetUniqueUint64() uint64 {
	for {
		delta := as.GetCustomTimestampAsMilliSeconds() - as.currentTimeStamp
		if delta > 0 {
			as.currentTimeStamp = as.GetCustomTimestampAsMilliSeconds()
			as.sequenceNumber = as.sequenceNumber + 1
			if as.sequenceNumber > 4095 {
				as.sequenceNumber = 0
			}

			break
		}
	}

	return (as.currentTimeStamp << 22) | (as.machineNumber << 12) | as.sequenceNumber
}
