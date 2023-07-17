package console

import (
	"errors"
	"fmt"

	"github.com/kdiot/alidns-console/utility"
)

type CmdMod struct {
	Cmd
	RecordId string
	RR       string
	Type     string
	Value    string
	TTL      int64
}

func (cmd *CmdMod) init() error {

	if err := cmd.Cmd.init("mod"); err != nil {
		return err
	}

	cmd.flagSet.StringVar(&cmd.RecordId, "id", "", "id of domain name record")
	cmd.flagSet.StringVar(&cmd.RR, "rr", "", "resource record")
	cmd.flagSet.StringVar(&cmd.Type, "type", "", "domain name record type")
	cmd.flagSet.StringVar(&cmd.Value, "value", "", "value of domain name record")
	cmd.flagSet.Int64Var(&cmd.TTL, "ttl", 0, "TTL(Time-To-Live), Retention time of domain name records in DNS servers.")

	return nil
}

func (cmd *CmdMod) Check() error {

	if err := cmd.Cmd.Check(); err != nil {
		return err
	}

	if cmd.RecordId == "" {
		return errors.New("domain name record id must be specified")
	}

	return nil
}

func (cmd *CmdMod) Execute() error {

	api, err := utility.NewAlidnsApi(cmd.DomainName, cmd.AccessKeyId, cmd.AccessKeySecret)
	if err != nil {
		return err
	}

	record, err := api.Retrieve(cmd.RecordId)
	if err != nil {
		return err
	}

	if cmd.RR != "" {
		record.RR = &cmd.RR
	}

	if cmd.Type != "" {
		if !utility.IsTypeValid(cmd.Type) {
			return errors.New("the domain name record type is invalid or not specified")
		}
		record.Type = &cmd.Type
	}

	if cmd.Value != "" {
		record.Value = &cmd.Value
	}

	if cmd.TTL > 0 {
		record.TTL = &cmd.TTL
	}

	if err = api.Update(record); err != nil {
		return err
	}

	fmt.Printf(
		"Domain name record successfully updated!\n"+
			"ID:     %s\n"+
			"DOMAIN: %s\n"+
			"RR:     %s\n"+
			"TYPE:   %s\n"+
			"VALUE:  %s\n"+
			"TTL:    %d\n"+
			"STATUS: %s\n",
		*record.RecordId,
		*record.DomainName,
		*record.RR,
		*record.Type,
		*record.Value,
		*record.TTL,
		*record.Status,
	)

	return nil
}

func NewCmdMod() *CmdMod {
	cmd := CmdMod{}
	if err := cmd.init(); err != nil {
		panic(err)
	} else {
		return &cmd
	}
}
