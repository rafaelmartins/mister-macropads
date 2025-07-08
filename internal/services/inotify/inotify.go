package inotify

import (
	"fmt"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"
)

type InotifyMask uint32

const (
	IN_ACCESS        InotifyMask = syscall.IN_ACCESS
	IN_MODIFY        InotifyMask = syscall.IN_MODIFY
	IN_ATTRIB        InotifyMask = syscall.IN_ATTRIB
	IN_CLOSE_WRITE   InotifyMask = syscall.IN_CLOSE_WRITE
	IN_CLOSE_NOWRITE InotifyMask = syscall.IN_CLOSE_NOWRITE
	IN_OPEN          InotifyMask = syscall.IN_OPEN
	IN_MOVED_FROM    InotifyMask = syscall.IN_MOVED_FROM
	IN_MOVED_TO      InotifyMask = syscall.IN_MOVED_TO
	IN_CREATE        InotifyMask = syscall.IN_CREATE
	IN_DELETE        InotifyMask = syscall.IN_DELETE
	IN_DELETE_SELF   InotifyMask = syscall.IN_DELETE_SELF
	IN_MOVE_SELF     InotifyMask = syscall.IN_MOVE_SELF
)

type Monitor struct {
	m        sync.Mutex
	fd       int
	wds      map[int]string
	handlers map[string][]func(f string, ev InotifyMask) error
}

func New() (*Monitor, error) {
	fd, err := syscall.InotifyInit1(syscall.IN_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return &Monitor{
		fd:       fd,
		wds:      map[int]string{},
		handlers: map[string][]func(f string, ev InotifyMask) error{},
	}, nil
}

func (m *Monitor) Close() error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.fd > 0 {
		return syscall.Close(m.fd)
	}
	return nil
}

func (m *Monitor) Reset() error {
	m.m.Lock()
	defer m.m.Unlock()

	for wd := range m.wds {
		if _, err := syscall.InotifyRmWatch(m.fd, uint32(wd)); err != nil {
			return err
		}
	}

	m.wds = map[int]string{}
	m.handlers = map[string][]func(f string, ev InotifyMask) error{}
	return nil
}

func (m *Monitor) AddWatch(f string, ev InotifyMask, handler func(f string, ev InotifyMask) error) error {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.handlers[f]; ok {
		m.handlers[f] = append(m.handlers[f], handler)
		return nil
	}

	wd, err := syscall.InotifyAddWatch(m.fd, f, uint32(ev))
	if err != nil {
		return err
	}

	m.wds[wd] = f
	m.handlers[f] = append(m.handlers[f], handler)
	return nil
}

func (m *Monitor) Listen() error {
	buf := make([]byte, 4096)

	for {
		n, err := syscall.Read(m.fd, buf)
		if err != nil {
			return err
		}

		pbuf := buf[:n]
		offset := 0

		for offset < len(pbuf) {
			event := (*syscall.InotifyEvent)(unsafe.Pointer(&pbuf[offset]))
			offset += 16 + int(event.Len)

			f, ok := m.wds[int(event.Wd)]
			if !ok {
				return fmt.Errorf("inotify: file to retrieve watched file/directory from wd: %d", event.Wd)
			}

			if event.Len > 0 {
				fnameB := pbuf[offset : offset+int(event.Len)]
				f = filepath.Join(f, string(fnameB[:clen(fnameB)]))
			}

			if err := func() error {
				m.m.Lock()
				defer m.m.Unlock()

				for _, hnd := range m.handlers[f] {
					if err := hnd(f, InotifyMask(event.Mask)); err != nil {
						return err
					}
				}
				return nil
			}(); err != nil {
				return err
			}
		}
	}
}

func clen(b []byte) int {
	for i := range b {
		if b[i] == 0 {
			return i
		}
	}
	return len(b)
}
