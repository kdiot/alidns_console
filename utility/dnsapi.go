package utility

import (
	"errors"
	"fmt"
	"os"
	"strings"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type DomainRecord = alidns.DescribeDomainRecordsResponseBodyDomainRecordsRecord

var RecordTypes = []string{
	"A",
	"NS",
	"MX",
	"TXT",
	"CNAME",
	"SRV",
	"AAAA",
	"CAA",
	"REDIRECT_URL",
	"FORWARD_URL",
}

func IsTypeValid(v string) bool {
	for _, t := range RecordTypes {
		if v == t {
			return true
		}
	}
	return false
}

type QueryInfo struct {
	RR     *string
	Type   *string
	Status *string
	Line   *string
}

type AlidnsApi struct {
	DomainName string
	client     *alidns.Client
	options    *util.RuntimeOptions
}

func NewAlidnsApi(domainName string, accessKeyId string, accessKeySecret string) (*AlidnsApi, error) {

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}

	if v := os.Getenv("ALIDNS_ENDPOINT"); v != "" {
		config.Endpoint = tea.String(v)
	} else {
		config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	}

	client, err := alidns.NewClient(config)
	if err != nil {
		return nil, err
	}

	api := &AlidnsApi{
		DomainName: domainName,
		client:     client,
		options: &util.RuntimeOptions{
			ConnectTimeout: tea.Int(5000),
			ReadTimeout:    tea.Int(5000),
			Autoretry:      tea.Bool(true),
			MaxAttempts:    tea.Int(3),
		},
	}

	return api, nil
}

func (api *AlidnsApi) describeDomainRecords(request *alidns.DescribeDomainRecordsRequest) (*alidns.DescribeDomainRecordsResponse, error) {
	if request == nil {
		request = &alidns.DescribeDomainRecordsRequest{}
	}
	request.DomainName = &api.DomainName
	response, err := func() (result *alidns.DescribeDomainRecordsResponse, e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				result = nil
				e = r
			}
		}()
		return api.client.DescribeDomainRecordsWithOptions(request, api.options)
	}()

	return response, err
}

func (api *AlidnsApi) describeDomainRecordInfo(request *alidns.DescribeDomainRecordInfoRequest) (*alidns.DescribeDomainRecordInfoResponse, error) {
	if request == nil {
		return nil, errors.New("AlidnsApi.describeDomainRecordInfo: The parameter request cannot be nil")
	}
	response, err := func() (result *alidns.DescribeDomainRecordInfoResponse, e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				result = nil
				e = r
			}
		}()
		return api.client.DescribeDomainRecordInfoWithOptions(request, api.options)
	}()

	return response, err
}

func (api *AlidnsApi) updateDomainRecord(request *alidns.UpdateDomainRecordRequest) (*alidns.UpdateDomainRecordResponse, error) {
	if request == nil {
		return nil, errors.New("AlidnsApi.updateDomainRecord: The parameter request cannot be nil")
	}
	response, err := func() (result *alidns.UpdateDomainRecordResponse, e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				result = nil
				e = r
			}
		}()
		return api.client.UpdateDomainRecordWithOptions(request, api.options)
	}()

	return response, err
}

func (api *AlidnsApi) addDomainRecord(request *alidns.AddDomainRecordRequest) (*alidns.AddDomainRecordResponse, error) {
	if request == nil {
		return nil, errors.New("AlidnsApi.addDomainRecord: The parameter request cannot be nil")
	}
	request.DomainName = &api.DomainName
	response, err := func() (result *alidns.AddDomainRecordResponse, e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				result = nil
				e = r
			}
		}()
		return api.client.AddDomainRecordWithOptions(request, api.options)
	}()

	return response, err
}

func (api *AlidnsApi) deleteDomainRecord(request *alidns.DeleteDomainRecordRequest) (*alidns.DeleteDomainRecordResponse, error) {
	if request == nil {
		return nil, errors.New("AlidnsApi.deleteDomainRecord: The parameter request cannot be nil")
	}
	response, err := func() (result *alidns.DeleteDomainRecordResponse, e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				result = nil
				e = r
			}
		}()
		return api.client.DeleteDomainRecordWithOptions(request, api.options)
	}()

	return response, err
}

func (api *AlidnsApi) Query(query *QueryInfo) ([]*DomainRecord, error) {
	request := &alidns.DescribeDomainRecordsRequest{
		RRKeyWord:  query.RR,
		Type:       query.Type,
		Status:     query.Status,
		Line:       query.Line,
		PageNumber: tea.Int64(1),
		PageSize:   tea.Int64(500),
	}

	var result []*DomainRecord
	for {
		response, err := api.describeDomainRecords(request)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = response.Body.DomainRecords.Record
		} else {
			result = append(result, response.Body.DomainRecords.Record...)
		}

		if len(response.Body.DomainRecords.Record) == 0 || len(result) >= int(*response.Body.TotalCount) {
			break
		}

		*request.PageNumber = *request.PageNumber + 1
	}

	return result, nil
}

func (api *AlidnsApi) Retrieve(recordId string) (*DomainRecord, error) {
	response, err := api.describeDomainRecordInfo(&alidns.DescribeDomainRecordInfoRequest{
		RecordId: tea.String(recordId),
	})
	if err != nil {
		return nil, err
	}

	record := &DomainRecord{
		DomainName: response.Body.DomainName,
		Line:       response.Body.Line,
		Locked:     response.Body.Locked,
		Priority:   response.Body.Priority,
		RR:         response.Body.RR,
		RecordId:   response.Body.RecordId,
		Status:     response.Body.Status,
		TTL:        response.Body.TTL,
		Type:       response.Body.Type,
		Value:      response.Body.Value,
	}

	return record, nil
}

func (api *AlidnsApi) Add(record *DomainRecord) (*DomainRecord, error) {
	response, err := api.addDomainRecord(&alidns.AddDomainRecordRequest{
		RR:    record.RR,
		TTL:   record.TTL,
		Type:  record.Type,
		Value: record.Value,
	})

	if err != nil {
		return nil, err
	} else {
		return api.Retrieve(*response.Body.RecordId)
	}
}

func (api *AlidnsApi) Update(record *DomainRecord) error {
	_, err := api.updateDomainRecord(&alidns.UpdateDomainRecordRequest{
		RecordId: record.RecordId,
		RR:       record.RR,
		TTL:      record.TTL,
		Type:     record.Type,
		Value:    record.Value,
	})
	return err
}

func (api *AlidnsApi) UpdateAndRetrieve(record *DomainRecord) (*DomainRecord, error) {
	if err := api.Update(record); err != nil {
		return nil, err
	}
	return api.Retrieve(*record.RecordId)
}

func (api *AlidnsApi) AutoUpdate(record *DomainRecord) error {

	// If RecordId is not null, directly call updateDomainRecord to update the domain name record
	if record.RecordId != nil {
		request := &alidns.UpdateDomainRecordRequest{
			RR:       record.RR,
			RecordId: record.RecordId,
			TTL:      record.TTL,
			Type:     record.Type,
			Value:    record.Value,
		}
		_, err := api.updateDomainRecord(request)
		if err != nil {
			if e, ok := err.(*tea.SDKError); ok {
				switch *e.Code {
				case "DomainRecordDuplicate":
					// There is no change in the domain name record, the update is considered successful
					return nil
				case "DomainRecordNotBelongToUser":
					// The original domain name record could not be found, it may be deleted
					record.RecordId = nil
				default:
					return err
				}
			} else {
				return err
			}
		} else {
			return nil
		}
	}

	response, err := api.describeDomainRecords(&alidns.DescribeDomainRecordsRequest{
		Type:      record.Type,
		RRKeyWord: record.RR,
	})
	if err != nil {
		return err
	}
	if len(response.Body.DomainRecords.Record) > 0 {
		old := response.Body.DomainRecords.Record[0]
		_, err := api.updateDomainRecord(&alidns.UpdateDomainRecordRequest{
			RR:       record.RR,
			RecordId: old.RecordId,
			TTL:      record.TTL,
			Type:     record.Type,
			Value:    record.Value,
		})
		if err != nil {
			if e, ok := err.(*tea.SDKError); ok {
				if strings.Compare(*e.Code, "DomainRecordDuplicate") != 0 {
					return err
				}
			} else {
				return err
			}
		}
		record.RecordId = old.RecordId
		record.DomainName = old.DomainName
		record.Line = old.Line
		record.Locked = old.Locked
		record.Priority = old.Priority
		record.Remark = old.Remark
		record.Status = old.Status
		record.Weight = old.Weight
		return nil
	} else {
		// domain name record does not exist, add a domain name record
		response, err := api.addDomainRecord(&alidns.AddDomainRecordRequest{
			RR:    record.RR,
			TTL:   record.TTL,
			Type:  record.Type,
			Value: record.Value,
		})
		if err != nil {
			return err
		}
		record.RecordId = response.Body.RecordId
		new, _ := api.Retrieve(*response.Body.RecordId)
		if new != nil {
			*record = *new
		}
		return nil
	}
}

func (api *AlidnsApi) Delete(recordId string) error {
	_, err := api.deleteDomainRecord(&alidns.DeleteDomainRecordRequest{
		RecordId: &recordId,
	})
	return err
}

func ErrMsg(err error) string {
	if e, ok := err.(*tea.SDKError); ok {
		return fmt.Sprintf("%s, %s", *e.Code, *e.Message)
	} else {
		return err.Error()
	}
}
