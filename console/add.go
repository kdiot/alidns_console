package console

import (
	"errors"
	"fmt"

	"github.com/kdiot/alidns-console/utility"
)

type CmdAdd struct {
	Cmd
	RR    string
	Type  string
	Value string
	TTL   int64
}

func (cmd *CmdAdd) init() error {

	if err := cmd.Cmd.init("add"); err != nil {
		return err
	}

	cmd.flagSet.StringVar(&cmd.RR, "rr", "", "RR(Resource-Record)")
	cmd.flagSet.StringVar(&cmd.Type, "type", "", "domain name record type")
	cmd.flagSet.StringVar(&cmd.Value, "value", "", "value of domain name record")
	cmd.flagSet.Int64Var(&cmd.TTL, "ttl", 0, "TTL(Time-To-Live), Retention time of domain name records in DNS servers.")

	return nil
}

func (cmd *CmdAdd) Check() error {

	if err := cmd.Cmd.Check(); err != nil {
		return err
	}

	if cmd.RR == "" {
		return errors.New("RR must be specified")
	}

	if !utility.IsTypeValid(cmd.Type) {
		return errors.New("the domain name record type is invalid or not specified")
	}

	if cmd.Value == "" {
		return errors.New("domain name record value not specified")
	}

	return nil
}

func (cmd *CmdAdd) Execute() error {

	api, err := utility.NewAlidnsApi(cmd.DomainName, cmd.AccessKeyId, cmd.AccessKeySecret)
	if err != nil {
		return err
	}

	record := &utility.DomainRecord{
		RR:    &cmd.RR,
		Type:  &cmd.Type,
		Value: &cmd.Value,
	}

	if cmd.TTL > 0 {
		record.TTL = &cmd.TTL
	}

	if record, err = api.Add(record); err != nil {
		return err
	}

	fmt.Printf(
		"Add domain name record successfully!\n"+
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

func NewCmdAdd() *CmdAdd {
	cmd := CmdAdd{}
	if err := cmd.init(); err != nil {
		panic(err)
	} else {
		return &cmd
	}
}
