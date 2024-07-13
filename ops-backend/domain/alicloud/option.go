package alicloud

import "fmt"

type Option func(d *DomainManager) error

func WithDomain(domain string) Option {
	return func(d *DomainManager) error {
		if domain == "" {
			return fmt.Errorf("domain is empty")
		}

		d.Domain = domain
		return nil
	}
}

func WithRegionID(regionID string) Option {
	return func(d *DomainManager) error {
		d.RegionID = regionID
		return nil
	}
}

func WithAccessKey(accessKey string) Option {
	return func(d *DomainManager) error {
		if accessKey == "" {
			return fmt.Errorf("access key is empty")
		}

		d.AccessKey = accessKey
		return nil
	}
}

func WithAccessSecret(accessSecret string) Option {
	return func(d *DomainManager) error {
		if accessSecret == "" {
			return fmt.Errorf("access secret is empty")
		}

		d.AccessSecret = accessSecret
		return nil
	}
}
