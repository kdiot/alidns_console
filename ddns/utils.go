package ddns

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/kdiot/alidns-console/utility"
)

func RegExpIPv4() *regexp.Regexp {
	result, _ := regexp.Compile(
		"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\." +
			"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\." +
			"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\." +
			"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)",
	)
	return result
}

func RegExpIPv6() *regexp.Regexp {
	result, _ := regexp.Compile(
		`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|` +
			`([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:)` +
			`{1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1` +
			`,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}` +
			`:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{` +
			`1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA` +
			`-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a` +
			`-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0` +
			`-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,` +
			`4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}` +
			`:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9` +
			`])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0` +
			`-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]` +
			`|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]` +
			`|1{0,1}[0-9]){0,1}[0-9]))`,
	)

	return result
}

func GetPublicIPv4(url string, re *regexp.Regexp) (string, error) {

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	var content string
	if data, err := io.ReadAll(response.Body); err != nil {
		return "", err
	} else {
		content = string(data)
	}

	if re == nil {
		re = RegExpIPv4()
	}

	if result := re.FindStringSubmatch(content); len(result) == 0 {
		return "", errors.New("no IPv4 address matched")
	} else {
		return result[0], nil
	}
}

func GetPublicIPv6(url string, re *regexp.Regexp) (string, error) {

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	var content string
	if data, err := io.ReadAll(response.Body); err != nil {
		return "", err
	} else {
		content = string(data)
	}

	if re == nil {
		re = RegExpIPv6()
	}

	if result := re.FindStringSubmatch(content); len(result) == 0 {
		return "", errors.New("no IPv6 address matched")
	} else {
		return result[0], nil
	}

}

func GetLocalIP(address string) *net.IP {
	sock, err := net.Dial("udp", address)
	if err != nil {
		return nil
	}
	addr := sock.LocalAddr().(*net.UDPAddr)
	return &addr.IP
}

type ExternalIP interface {
	GetIP(url string) (net.IP, error)
	Refresh() (net.IP, bool)
}

type ExternalIPv4 struct {
	IP        net.IP
	re        *regexp.Regexp
	providers []string
}

func (ipv4 *ExternalIPv4) GetIP(url string) (net.IP, error) {
	if ip, err := GetPublicIPv4(url, ipv4.re); err != nil {
		return nil, err
	} else {
		return net.ParseIP(ip), nil
	}
}

func (ipv4 *ExternalIPv4) Refresh() (net.IP, bool) {
	count := len(ipv4.providers)
	for i := 0; i < count; i++ {
		url := ipv4.providers[0]
		if ip, err := ipv4.GetIP(url); err != nil {
			ipv4.providers = append(ipv4.providers[1:], url)
			utility.Errorf("Failed to obtain public IPv4 address.[vender:%s, error:%s]", url, err.Error())
			continue
		} else {
			if !ip.Equal(ipv4.IP) {
				ipv4.IP = ip
				return ipv4.IP, true
			} else {
				return ipv4.IP, false
			}
		}
	}
	return ipv4.IP, false
}

func NewExternalIPv4(providers []string) *ExternalIPv4 {
	ipv4 := &ExternalIPv4{
		IP: net.IPv4zero,
		re: RegExpIPv4(),
	}
	if len(providers) == 0 {
		ipv4.providers = []string{
			"https://myip.ipip.net",
			"https://ip.tool.lu",
			"https://myip.dnsomatic.com",
			"https://api4.ipify.org",
			"https://ipv4.jsonip.com",
		}
	} else {
		ipv4.providers = providers
	}
	return ipv4
}

type ExternalIPv6 struct {
	IP        net.IP
	re        *regexp.Regexp
	providers []string
}

func (ipv6 *ExternalIPv6) GetIP(address string) (net.IP, error) {

	if address[:4] == "http" {
		if ip, err := GetPublicIPv6(address, ipv6.re); err != nil {
			return nil, err
		} else {
			return net.ParseIP(ip), nil
		}
	} else {
		sock, err := net.Dial("udp", fmt.Sprintf("[%s]:0", address))
		if err != nil {
			return net.IPv6zero, err
		} else {
			defer sock.Close()
			return sock.LocalAddr().(*net.UDPAddr).IP, nil
		}
	}
}

func (ipv6 *ExternalIPv6) Refresh() (net.IP, bool) {
	count := len(ipv6.providers)
	for i := 0; i < count; i++ {
		url := ipv6.providers[0]
		if ip, err := ipv6.GetIP(url); err != nil {
			ipv6.providers = append(ipv6.providers[1:], url)
			continue
		} else {
			if !ip.Equal(ipv6.IP) {
				ipv6.IP = ip
				return ipv6.IP, true
			} else {
				return ipv6.IP, false
			}
		}
	}
	return ipv6.IP, false
}

func NewExternalIPv6(providers []string) *ExternalIPv6 {
	ipv6 := &ExternalIPv6{
		IP: net.IPv6zero,
		re: RegExpIPv6(),
	}
	if len(providers) == 0 {
		ipv6.providers = []string{
			"2400:3200::1",      // ALIBABA1
			"2400:3200:baba::1", // ALIBABA2
			"2400:da00::6666",   // BAIDU
			"240e:4c:4008::1",   // CT1
			"240e:4c:4808::1",   // CT2
			//"https://v6.ident.me",
			//"https://api6.ipify.org",
			//"https://ipv6.jsonip.com",
		}
	} else {
		ipv6.providers = providers
	}
	return ipv6
}
