package ddns

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kdiot/alidns-console/utility"
)

type Daemon struct {
	config   *Config
	services []*UpdateService
	ipv4     ExternalIP
	ipv6     ExternalIP
}

func (daemon *Daemon) init(conf *Config) error {
	daemon.config = conf
	daemon.ipv4 = NewExternalIPv4(nil)
	daemon.ipv6 = NewExternalIPv6(nil)

	for _, d := range conf.DomainList {
		if err := d.Check(); err != nil {
			utility.Errorf("Dynamic domain name configuration error: %s", err.Error())
			continue
		}
		service, err := NewUpdateService(d, conf)
		if err != nil {
			utility.Errorf("An error occurred while creating the AliDDNS object, the reason for the error: %s", utility.ErrMsg(err))
			continue
		}
		daemon.services = append(daemon.services, service)
	}

	return nil
}

func (daemon *Daemon) doCheck() {
	if ip, changed := daemon.ipv4.Refresh(); changed {
		utility.Infof("Detected that the IPv4 address(%s) has changed, preparing to update the domain name record...", ip.String())
		for _, d := range daemon.services {
			if d.record.Type != nil && *d.record.Type == "A" {
				d.IpAddrChan <- &ip
			}
		}
	} else {
		utility.Debug("IPv4 addresses have not changed, no need to update domain name records.")
	}

	if ip, changed := daemon.ipv6.Refresh(); changed {
		utility.Infof("Detected that the IPv6 address(%s) has changed, preparing to update the domain name record...", ip.String())
		for _, d := range daemon.services {
			if d.record.Type != nil && *d.record.Type == "AAAA" {
				d.IpAddrChan <- &ip
			}
		}
	} else {
		utility.Debug("IPv6 addresses have not changed, no need to update domain name records.")
	}
}

func (daemon *Daemon) Run() {

	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for _, d := range daemon.services {
		wg.Add(1)
		go d.Routine(ctx, &wg)
	}

	loop := true
	for loop {
		select {
		case s := <-sigs:
			utility.Infof("System signal: %s", s.String())
			loop = false
		case <-time.After(daemon.config.CheckInterval * time.Second):
			daemon.doCheck()
		}
	}

	cancel()
	wg.Wait()

	utility.Info("Daemon exit safely!")
}

func NewDaemon(conf *Config) (*Daemon, error) {
	daemon := &Daemon{}
	if err := daemon.init(conf); err != nil {
		return nil, err
	}
	return daemon, nil
}
