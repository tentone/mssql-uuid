package mssql-uuid

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

var generator = newRFC4122Generator()

// NewV4 returns random generated UUID.
func NewV4() UUID {
	var uuid, _ = generator.NewV4()
	return uuid
}

// Default generator implementation.
type Generator struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex
	rand              io.Reader
	epochFunc         epochFunc
	hwAddrFunc        hwAddrFunc
	lastTime          uint64
	clockSequence     uint16
	hardwareAddr      [6]byte
}

func newRFC4122Generator() *Generator {
	return &Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

// NewV4 returns random generated UUID.
func (g *Generator) NewV4() (UUID, error) {
	var u = UUID{}
	var err error

	_, err = io.ReadFull(g.rand, u[:])
	if err != nil {
		return NilUUID, err
	}

	u.SetVersion(V4)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// Returns hardware address.
func defaultHWAddrFunc() (net.HardwareAddr, error) {
	var ifaces, err = net.Interfaces()
	if err != nil {
		return []byte{}, err
	}

	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}

	return []byte{}, fmt.Errorf("No HW address found.")
}
