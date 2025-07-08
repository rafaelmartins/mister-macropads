package dataproviders

import (
	"net"
	"slices"

	"github.com/rafaelmartins/mister-macropads/internal/services"
)

var IpAddr IpAddrType

type IpAddrType struct {
	initialized []string
	valueMap    map[string]string
}

func (i *IpAddrType) Get(backend Backend, itf string) (string, error) {
	if !slices.Contains(i.initialized, itf) {
		if err := services.AddIpAddrWatch(itf, func(itf string, ip net.IP) error {
			if i.valueMap == nil {
				i.valueMap = map[string]string{}
			}
			if ip != nil {
				i.valueMap[itf] = ip.String()
			} else {
				i.valueMap[itf] = ""
			}
			return backend.ScreenRender()
		}); err != nil {
			return "", err
		}
		i.initialized = append(i.initialized, itf)
	}
	return i.valueMap[itf], nil
}
