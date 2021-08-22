// +build windows

package main

import (
	"time"

	"golang.org/x/sys/windows/svc"
)


type windowsHandler struct {
	ipc chan svc.Status
}

func (w windowsHandler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{
		State: svc.Running, 
		Accepts: svc.AcceptShutdown | svc.AcceptStop,
	}

loop:
	for {
		select {
		case cr := <-r:
			switch cr.Cmd {
			case svc.Interrogate:
				s <- cr.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				s <- cr.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				break loop
			}
		case signal := <-w.ipc:
			s <- signal
		}
	}
	s <- svc.Status{State: svc.StopPending}
	return
}

func init() {
	isManaged, err := svc.IsWindowsService()
	if err != nil {
		panic(err)
	}
	// This enables use of Windows Services, which most users
	// will most likely use
	if isManaged {
		win := windowsHandler{}
		go func() {
			if err = svc.Run("mpdrp", win); err != nil {
				panic(err)
			}
		}()
	}
}