package console

import (
	"errors"
	"flag"
	"os"
	"os/user"
	"path/filepath"
)

const (
	envAccessKeyId     = "ALIDNS_ACCESSKEYID"
	envAccessKeySecret = "ALIDNS_ACCESSKEYSECRET"
	envDomainName      = "ALIDNS_DOMAINNAME"
)

type Command interface {
	Name() string
	Parse(arguments []string) error
	Check() error
	Execute() error
	Usage()
}

type Cmd struct {
	flagSet *flag.FlagSet
	name    string
	Profile
	ProfileName string
}

func (cmd *Cmd) init(name string) error {
	cmd.name = name
	cmd.flagSet = flag.NewFlagSet(cmd.name, flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.Profile.AccessKeyId, "key", "", "access key id")
	cmd.flagSet.StringVar(&cmd.Profile.AccessKeySecret, "secret", "", "access key secret")
	cmd.flagSet.StringVar(&cmd.Profile.DomainName, "domain", "", "domain name")
	cmd.flagSet.StringVar(&cmd.ProfileName, "profile", "", "A profile containing information such as user authentication.")
	return nil
}

func (cmd *Cmd) Name() string {
	return cmd.name
}

func (cmd *Cmd) Usage() {
	cmd.flagSet.Usage()
}

func (cmd *Cmd) Check() error {

	if cmd.AccessKeyId == "" {
		return errors.New("access key id is not specified")
	}

	if cmd.AccessKeySecret == "" {
		return errors.New("access key secret is not specified")
	}

	if cmd.DomainName == "" {
		return errors.New("domain name must be specified")
	}

	return nil
}

func (cmd *Cmd) Parse(arguments []string) error {

	if err := cmd.flagSet.Parse(arguments); err != nil {
		return err
	}

	profile := Profile{}
	if cmd.ProfileName == "" {
		if user, err := user.Current(); err == nil {
			fileName := filepath.Join(user.HomeDir, ".alidns")
			if _, err := os.Stat(fileName); err == nil {
				profile.Load(fileName)
			}
			if accessKeyId := os.Getenv(envAccessKeyId); accessKeyId != "" {
				profile.AccessKeyId = accessKeyId
			}
			if accessKeySecret := os.Getenv(envAccessKeySecret); accessKeySecret != "" {
				profile.AccessKeySecret = accessKeySecret
			}
			if domainName := os.Getenv(envDomainName); domainName != "" {
				profile.DomainName = domainName
			}
		}
	} else {
		profile.Load(cmd.ProfileName)
	}

	if cmd.AccessKeyId == "" {
		if profile.AccessKeyId != "" {
			cmd.AccessKeyId = profile.AccessKeyId
		} else {
			cmd.AccessKeyId = os.Getenv(envAccessKeyId)
		}
	}

	if cmd.AccessKeySecret == "" {
		if profile.AccessKeySecret != "" {
			cmd.AccessKeySecret = profile.AccessKeySecret
		} else {
			cmd.AccessKeySecret = os.Getenv(envAccessKeySecret)
		}
	}

	if cmd.DomainName == "" {
		if profile.DomainName != "" {
			cmd.DomainName = profile.DomainName
		} else {
			cmd.DomainName = os.Getenv(envDomainName)
		}
	}

	return nil
}
