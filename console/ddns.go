package console

import (
	"errors"
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/kdiot/alidns-console/ddns"
	"github.com/kdiot/alidns-console/utility"
)

type CmdDdns struct {
	Cmd
	RR            string
	Type          string
	TTL           int64
	Network       string
	ConfigFile    string
	CheckInterval time.Duration
	RetryInterval time.Duration
	LogFile       string
	LogLevel      utility.LogLevel
	config        *ddns.Config
}

func (cmd *CmdDdns) init() error {
	if err := cmd.Cmd.init("ddns"); err != nil {
		return err
	}

	cmd.LogLevel = utility.LOG_INFO

	cmd.flagSet.StringVar(&cmd.RR, "rr", "@", "RR(Resource-Record)")
	cmd.flagSet.StringVar(&cmd.Type, "type", "", "domain name record type, only 'A' or 'AAAA' can be selected.")
	cmd.flagSet.Int64Var(&cmd.TTL, "ttl", 600, "TTL(Time-To-Live), Retention time of domain name records in DNS servers.")
	cmd.flagSet.StringVar(&cmd.Network, "network", "", "network config, specify the local address and prefix length, IPv6 only")
	cmd.flagSet.StringVar(&cmd.ConfigFile, "conf", "", "config file name")
	cmd.flagSet.DurationVar(&cmd.CheckInterval, "chkIntvl", 10, "check whether the IP address has changed every X seconds")
	cmd.flagSet.DurationVar(&cmd.RetryInterval, "retryIntvl", 30, "retry interval after update domain name record fails")
	cmd.flagSet.StringVar(&cmd.LogFile, "log", "", "log file")
	cmd.flagSet.Var(&cmd.LogLevel, "loglvl", "log level. LogLevel[debug,info,warning,error,fatal]")

	return nil
}

func (cmd *CmdDdns) Check() error {
	return nil
}

func (cmd *CmdDdns) Parse(arguments []string) error {
	var err error

	if err = cmd.Cmd.Parse(arguments); err != nil {
		return err
	}

	if cmd.ConfigFile != "" {
		cmd.config, err = ddns.LoadConfig(cmd.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration file! [%s, %s]", cmd.ConfigFile, err.Error())
		}

		if !cmd.config.LogLevel.IsValid() {
			cmd.config.LogLevel = cmd.LogLevel
		}

		cmd.config.LogFile = utility.DefaultIfEmpty(cmd.config.LogFile, &cmd.LogFile)
		if tea.StringValue(cmd.config.AccessKeyId) == "" || tea.StringValue(cmd.config.AccessKeySecret) == "" {
			cmd.config.AccessKeyId = &cmd.AccessKeyId
			cmd.config.AccessKeySecret = &cmd.AccessKeySecret
		}
		cmd.config.DomainName = utility.DefaultIfEmpty(cmd.config.DomainName, &cmd.DomainName)

	} else {
		if cmd.RR == "" {
			return errors.New("RR must be specified")
		} else if cmd.Type != "A" && cmd.Type != "AAAA" {
			return errors.New("domain name record type must 'A' or 'AAAA'")
		}
		cmd.config = &ddns.Config{
			AccessKeyId:     &cmd.AccessKeyId,
			AccessKeySecret: &cmd.AccessKeySecret,
			DomainName:      &cmd.DomainName,
			CheckInterval:   cmd.CheckInterval,
			RetryInterval:   cmd.RetryInterval,
			DomainList: []*ddns.DDNS{
				{
					RR:      &cmd.RR,
					Type:    &cmd.Type,
					Network: &cmd.Network,
					TTL:     &cmd.TTL,
				},
			},
		}
	}

	for _, d := range cmd.config.DomainList {
		if tea.StringValue(d.AccessKeyId) == "" || tea.StringValue(d.AccessKeySecret) == "" {
			d.AccessKeyId = cmd.config.AccessKeyId
			d.AccessKeySecret = cmd.config.AccessKeySecret
		}
		d.DomainName = utility.DefaultIfEmpty(d.DomainName, cmd.config.DomainName)
	}

	if tea.StringValue(cmd.config.LogFile) != "" {
		if err = utility.SetLogFile(*cmd.config.LogFile); err != nil {
			err = fmt.Errorf("failed to open log file!(%s, %s)", *cmd.config.LogFile, err.Error())
			utility.Warning(err.Error())
		}
	}
	utility.SetLogLevel(cmd.config.LogLevel)

	return nil
}

func (cmd *CmdDdns) Execute() error {

	daemon, err := ddns.NewDaemon(cmd.config)
	if err != nil {
		return err
	}
	daemon.Run()

	return nil
}

func NewCmdDdns() *CmdDdns {
	cmd := CmdDdns{}
	if err := cmd.init(); err != nil {
		panic(err)
	} else {
		return &cmd
	}
}
