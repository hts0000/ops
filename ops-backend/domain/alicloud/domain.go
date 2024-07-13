package alicloud

import (
	"log"

	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	domainpb "github.com/hts0000/ops-backend/domain/api/gen/v1"
	"github.com/hts0000/ops-backend/shared/server"
	"github.com/hts0000/ops-backend/shared/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DomainManager struct {
	AccessKey    string
	AccessSecret string
	RegionID     string
	Domain       string
	Logger       *zap.Logger
}

func NewDomainManager(opts ...Option) (*DomainManager, error) {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	d := &DomainManager{Logger: logger}
	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}
	return d, nil
}

func (d *DomainManager) GetDomains(page, pageSize int64) (*domainpb.GetDomainsResponse, error) {
	dnsClient, err := dns.NewClient(&client.Config{
		AccessKeyId:     &d.AccessKey,
		AccessKeySecret: &d.AccessSecret,
		RegionId:        &d.RegionID,
	})
	if err != nil {
		d.Logger.Error("cannot create dns client", zap.Error(err))
		return nil, errors.Wrap(err, "cannot create dns client")
	}

	descResp, err := dnsClient.DescribeDomainRecordsWithOptions(&dns.DescribeDomainRecordsRequest{
		PageNumber: tea.Int64(page),
		PageSize:   tea.Int64(pageSize),
		DomainName: &d.Domain,
	}, &service.RuntimeOptions{
		ConnectTimeout: tea.Int(10),
		ReadTimeout:    tea.Int(10),
	})
	if err != nil {
		d.Logger.Error("cannot describe domain records", zap.Error(err))
		return nil, err
	}

	records := make([]*domainpb.DomainRecord, 0)
	for _, record := range descResp.Body.DomainRecords.Record {
		r := &domainpb.DomainRecord{
			Rr:     util.GetOrZero(record.RR),
			Line:   util.GetOrZero(record.Line),
			Remark: util.GetOrZero(record.Remark),
			Value:  util.GetOrZero(record.Value),
			TtlSec: util.GetOrZero(record.TTL),
		}
		if record.Status != nil {
			r.Status = domainpb.Status(domainpb.Status_value[*record.Status])
		}
		if record.Type != nil {
			r.Type = domainpb.Type(domainpb.Type_value[*record.Type])
		}
		records = append(records, r)
	}

	return &domainpb.GetDomainsResponse{
		Total:        uint64(len(records)),
		Page:         uint64(1),
		PageSize:     uint64(len(records)),
		DomainRecord: records,
	}, nil
}
