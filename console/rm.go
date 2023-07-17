package console

import (
	"errors"
	"fmt"

	"github.com/kdiot/alidns-console/utility"
)

type CmdRm struct {
	Cmd
	RecordId string
}

func (cmd *CmdRm) init() error {

	if err := cmd.Cmd.init("rm"); err != nil {
		return err
	}

	cmd.flagSet.StringVar(&cmd.RecordId, "id", "", "id of domain name record")

	return nil
}

func (cmd *CmdRm) Check() error {

	if err := cmd.Cmd.Check(); err != nil {
		return err
	}

	if cmd.RecordId == "" {
		return errors.New("domain name record id must be specified")
	}

	return nil
}

func (cmd *CmdRm) Execute() error {

	api, err := utility.NewAlidnsApi(cmd.DomainName, cmd.AccessKeyId, cmd.AccessKeySecret)
	if err != nil {
		return err
	}

	if err = api.Delete(cmd.RecordId); err != nil {
		return err
	} else {
		fmt.Printf("The domain name record with ID '%s' was successfully deleted!\n", cmd.RecordId)
		return nil
	}
}

func NewCmdRm() *CmdRm {
	cmd := CmdRm{}
	if err := cmd.init(); err != nil {
		panic(err)
	} else {
		return &cmd
	}
}
