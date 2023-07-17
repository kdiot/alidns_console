package ddns

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/kdiot/alidns-console/utility"
)

type Network struct {
	IP   net.IP
	Mask net.IPMask
}

type UpdateService struct {
	api           *utility.AlidnsApi
	record        *utility.DomainRecord
	network       *Network
	retryTimer    *time.Timer
	retryInterval time.Duration
	IpAddrChan    chan *net.IP
}

func (s *UpdateService) Type() *string {
	return s.record.Type
}

func (s *UpdateService) Init(d *DDNS) error {

	var err error
	if s.api == nil {
		if s.api, err = utility.NewAlidnsApi(*d.DomainName, *d.AccessKeyId, *d.AccessKeySecret); err != nil {
			return err
		}
	}

	s.IpAddrChan = make(chan *net.IP)

	s.record = &utility.DomainRecord{
		DomainName: d.DomainName,
		RR:         d.RR,
		Type:       d.Type,
		TTL:        d.TTL,
	}

	if d.Network != nil && *d.Network != "" {
		if ip, ipnet, err := net.ParseCIDR(*d.Network); err != nil {
			return fmt.Errorf(`the network address configuration is illegal!["Network": "%s"]`, *d.Network)
		} else {
			s.network = &Network{IP: ip, Mask: ipnet.Mask}
		}
	}

	return nil
}

func (s *UpdateService) Update(ip *net.IP) {

	if s.retryTimer != nil {
		s.retryTimer.Stop()
		s.retryTimer = nil
	}

	if ip == nil {
		return
	}

	if s.network != nil {
		prefix := ip.Mask(s.network.Mask)
		n := len(*ip)
		newip := make(net.IP, n)
		for i := 0; i < n; i++ {
			newip[i] = prefix[i] | s.network.IP[i]
		}
		s.record.Value = tea.String(newip.String())
	} else {
		s.record.Value = tea.String(ip.String())
	}

	utility.Debug("UpdateService.Update: begin update...")
	if err := s.api.AutoUpdate(s.record); err != nil {
		if e, ok := err.(*tea.SDKError); ok {
			utility.Errorf("Update dynamic domain name record '%s.%s' failed! Error message: %s, %s",
				tea.StringValue(s.record.RR),
				tea.StringValue(s.record.DomainName),
				tea.StringValue(e.Code),
				tea.StringValue(e.Message),
			)
		} else {
			utility.Errorf("Update dynamic domain name record '%s.%s' failed! Error message: %s",
				tea.StringValue(s.record.RR),
				tea.StringValue(s.record.DomainName),
				err.Error(),
			)
		}
		if s.retryInterval > 0 {
			s.retryTimer = time.AfterFunc(time.Second*s.retryInterval, func() {
				s.retryTimer = nil
				s.Update(ip)
			})
		}
	} else {
		utility.Infof("Update the dynamic domain name record '%s.%s' successfully! The new IP address is: %s",
			tea.StringValue(s.record.RR),
			tea.StringValue(s.record.DomainName),
			tea.StringValue(s.record.Value),
		)
	}
}

func (s *UpdateService) Close() {
	close(s.IpAddrChan)
}

func (s *UpdateService) Routine(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		s.Close()
		wg.Done()
	}()

	for {
		select {
		case ip := <-s.IpAddrChan:
			s.Update(ip)
		case <-ctx.Done():
			if s.retryTimer != nil {
				s.retryTimer.Stop()
				s.retryTimer = nil
			}
			return
		}
	}
}

func NewUpdateService(d *DDNS, conf *Config) (*UpdateService, error) {
	s := &UpdateService{
		retryInterval: 30,
	}

	if conf != nil {
		s.retryInterval = conf.RetryInterval
	}

	if err := s.Init(d); err != nil {
		return nil, err
	}

	return s, nil
}
