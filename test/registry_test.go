package test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-kratos/nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func TestRegistry(t *testing.T) {
	ip := ""
	serviceName := "golang-sms@grpc"
	ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	_, e := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})

	time.Sleep(time.Second)

	is, e := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})

	t.Logf("%#v, %v", is, e)

	time.Sleep(time.Second)
	r := registry.New(client)

	go func() {
		w, e := r.Watch(ctx, "golang-sms@grpc")
		if e != nil {
			log.Fatal(e)
		}
		for {
			res, err := w.Next()
			if err != nil {
				return
			}
			log.Printf("watch: %d", len(res))
			for _, r := range res {
				log.Printf("next: %+v", r)
			}
		}
	}()

	time.Sleep(time.Second)

	ins, e := r.GetService(ctx, serviceName)
	t.Logf("e:%v", e)
	for _, in := range ins {
		t.Logf("ins: %#v", in)
	}

	time.Sleep(time.Second)

	select {}
}

func TestRegistryMany(t *testing.T) {
	ip := ""
	serviceName := "golang-sms@grpc"
	//ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	_, e := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f1",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})

	_, e = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f2",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})

	_, e = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f3",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})

	time.Sleep(time.Second)

	is, e := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})

	for _, host := range is.Hosts {
		t.Logf("host: %#v,e: %v", host, e)
	}

	select {}
}
