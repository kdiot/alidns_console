package console

import (
	"fmt"
	"os"

	"github.com/kdiot/alidns-console/utility"
	"github.com/olekukonko/tablewriter"
)

type CmdLs struct {
	Cmd
	RR     string
	Type   string
	Line   string
	Status string
}

func (cmd *CmdLs) init() error {

	if err := cmd.Cmd.init("ls"); err != nil {
		return err
	}

	cmd.flagSet.StringVar(&cmd.RR, "rr", "", "resource record")
	cmd.flagSet.StringVar(&cmd.Type, "type", "", "domain name record type, A|AAAA|CNAME|TXT")
	cmd.flagSet.StringVar(&cmd.Line, "line", "", "Line")
	cmd.flagSet.StringVar(&cmd.Status, "status", "", "status")

	return nil
}

func (cmd *CmdLs) Check() error {

	if err := cmd.Cmd.Check(); err != nil {
		return err
	}
	return nil
}

func (cmd *CmdLs) Execute() error {

	query := &utility.QueryInfo{}
	if cmd.RR != "" {
		query.RR = &cmd.RR
	}
	if cmd.Type != "" {
		query.Type = &cmd.Type
	}
	if cmd.Line != "" {
		query.Line = &cmd.Line
	}
	if cmd.Status != "" {
		query.Status = &cmd.Status
	}

	api, err := utility.NewAlidnsApi(cmd.DomainName, cmd.AccessKeyId, cmd.AccessKeySecret)
	if err != nil {
		return err
	}

	records, err := api.Query(query)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "RR", "TYPE", "VALUE", "TTL", "STATUS"})
	for _, record := range records {
		table.Append([]string{
			*record.RecordId,
			*record.RR,
			*record.Type,
			*record.Value,
			fmt.Sprintf("%d", *record.TTL),
			*record.Status,
		})
	}
	table.Render()

	return nil
}

func NewCmdLs() *CmdLs {
	cmd := CmdLs{}
	if err := cmd.init(); err != nil {
		panic(err)
	} else {
		return &cmd
	}
}
