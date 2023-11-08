package dao

import (
	"context"

	etcd "go.etcd.io/etcd/client/v3"
)

var Etcd = &dEtcd{}

type dEtcd struct {
	client *etcd.Client
}

func (d *dEtcd) SetClient(c *etcd.Client) {
	d.client = c
}

func (d *dEtcd) Get(ctx context.Context, k string) ([]byte, int64, error) {
	resp, err := d.client.Get(ctx, k)
	if err != nil {
		return nil, 0, err
	}
	if resp.Count == 0 {
		return nil, resp.Header.Revision, nil
	}
	return resp.Kvs[0].Value, resp.Header.Revision, nil
}

func (d *dEtcd) Put(ctx context.Context, k string, v []byte) error {
	_, err := d.client.Put(ctx, k, string(v))
	return err
}

func (d *dEtcd) Watch(ctx context.Context, k string, rev int64, callback func(data []byte, rev int64)) error {
	c := d.client.Watch(ctx, k, etcd.WithRev(rev))

	var err error
	for resp := range c {
		err = resp.Err()
		if err != nil || resp.Canceled {
			break
		}

		if len(resp.Events) > 0 && resp.Events[0].Kv != nil {
			callback(resp.Events[0].Kv.Value, resp.Header.Revision)
		} else {
			callback(nil, resp.Header.Revision)
		}
	}
	return err
}

func (d *dEtcd) GetWatch(ctx context.Context, k string, callback func(data []byte, rev int64)) error {
	data, rev, err := d.Get(ctx, k)
	if err != nil {
		return err
	}
	callback(data, rev)
	return d.Watch(ctx, k, rev+1, callback)
}

func (d *dEtcd) GetConfig(ctx context.Context, k string) ([]byte, error) {
	data, _, err := d.Get(ctx, k)
	return data, err
}

func (d *dEtcd) SetConfig(ctx context.Context, k string, v []byte) error {
	return d.Put(ctx, k, v)
}
