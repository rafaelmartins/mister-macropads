package ipaddr

import (
	"net"
	"slices"
	"sync"
	"time"
)

type Monitor struct {
	m        sync.Mutex
	ch       chan struct{}
	itfs     []string
	ipMap    map[string]net.IP
	handlers map[string][]func(itf string, ip net.IP) error
}

func New() *Monitor {
	return &Monitor{
		ch:       make(chan struct{}),
		itfs:     []string{},
		ipMap:    map[string]net.IP{},
		handlers: map[string][]func(itf string, ip net.IP) error{},
	}
}

func (m *Monitor) Close() error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.ch != nil {
		close(m.ch)
		m.ch = nil
	}
	return nil
}

func (m *Monitor) AddWatch(itf string, handler func(itf string, ip net.IP) error) error {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.handlers[itf]; ok {
		m.handlers[itf] = append(m.handlers[itf], handler)
		return nil
	}

	if _, err := net.InterfaceByName(itf); err != nil {
		return err
	}

	m.handlers[itf] = append(m.handlers[itf], handler)
	m.itfs = append(m.itfs, itf)
	return nil
}

func (m *Monitor) Listen() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		m.m.Lock()
		itfs := slices.Clone(m.itfs)
		m.m.Unlock()

		for _, itfn := range itfs {
			itf, err := net.InterfaceByName(itfn)
			if err != nil {
				return err
			}

			addrs, err := itf.Addrs()
			if err != nil {
				return err
			}

			found := false
			for _, addr := range addrs {
				if netaddr, ok := addr.(*net.IPNet); ok {
					if ip := netaddr.IP.To4(); ip != nil {
						m.m.Lock()
						if oip, exist := m.ipMap[itfn]; exist && ip.Equal(oip) {
							m.m.Unlock()
							found = true
							break
						}
						for _, handler := range m.handlers[itfn] {
							if err := handler(itf.Name, ip); err != nil {
								m.m.Unlock()
								return err
							}
						}
						m.ipMap[itfn] = ip
						m.m.Unlock()
						found = true
						break
					}
				}
			}
			if !found {
				m.m.Lock()
				if oip, exist := m.ipMap[itfn]; exist && oip == nil {
					m.m.Unlock()
					continue
				}
				for _, handler := range m.handlers[itfn] {
					if err := handler(itf.Name, nil); err != nil {
						m.m.Unlock()
						return err
					}
				}
				m.ipMap[itfn] = nil
				m.m.Unlock()
			}
		}

		select {
		case <-m.ch:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
