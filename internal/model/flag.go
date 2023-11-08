package model

import "errors"

type Flag struct {
	Storage       string
	LocalConfig   string
	EtcdEndpoints []string
	EtcdConfigKey string
	EtcdCA        string
	EtcdCert      string
	EtcdKey       string
	EtcdUsername  string
	EtcdPassword  string
}

func (f *Flag) Validate() error {
	if f.Storage == "local" {
		if f.LocalConfig == "" {
			return errors.New("local-config required when storage is local")
		}
		return nil
	}

	if f.Storage == "etcd" {
		if len(f.EtcdEndpoints) == 0 {
			return errors.New("etcd-endpoints required when storage is etcd")
		}

		if f.EtcdConfigKey == "" {
			return errors.New("etcd-config-key required when storage is etcd")
		}

		if f.EtcdCert != "" && f.EtcdKey == "" {
			return errors.New("etcd-key required when etcd-cert set")
		}

		if f.EtcdCert == "" && f.EtcdKey != "" {
			return errors.New("etcd-cert required when etcd-key set")
		}
		return nil
	}

	return errors.New("unknown storage")
}
