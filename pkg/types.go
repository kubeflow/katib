package pkg

import (
	"os"
)

const (
	VizierServiceIPEnvName = "VIZIER_CORE_PORT_6789_TCP_ADDR"
	VizierServicePortEnvName = "VIZIER_CORE_PORT_6789_TCP_PORT"
	VizierServiceNamespaceEnvName = "VIZIER_CORE_NAMESPACE"
	VizierService = "vizier-core"
	VizierPort = "6789"
	ManagerAddr = VizierService + ":" + VizierPort
)

func GetManagerAddr() string {
	ns := os.Getenv(VizierServiceNamespaceEnvName)
	if len(ns) == 0 {
		addr := os.Getenv(VizierServiceIPEnvName)
		port := os.Getenv(VizierServicePortEnvName)
		if len(addr) > 0 && len(port) > 0 {
			return addr + ":" + port
		} else {
			return ManagerAddr
		}
	} else {
		return VizierService + "." + ns + ":"+ VizierPort
	}
}