package dao

import (
	"context"
	"fmt"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	//"github.com/go-kit/kit/transport/grpc"
	callapi "callServer/api" // <-direct use server proto api client.go content
	"callServer/internal/model"
	svrapi "callServer/smallapi" // <-- use for call it!
	"time"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/bilibili/kratos/pkg/time"
	grpcempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&model.Article{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Article(c context.Context, id int64) (*model.Article, error)
	SayHello(c context.Context, req *callapi.HelloReq) (resp *grpcempty.Empty, err error)
}

// dao dao.
type dao struct {
	db              *sql.DB
	redis           *redis.Redis
	mc              *memcache.Memcache
	remotesvrClient svrapi.AsmallsClient
	cache           *fanout.Fanout
	demoExpire      int32
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}

	grpccfg := &warden.ClientConfig{}
	//paladin.Get("grpc.toml").UnmarshalTOML(grpccfg)
	if err = paladin.Get("grpc.toml").UnmarshalTOML(grpccfg); err != nil {
		return
	}
	var grpcClient svrapi.AsmallsClient
	if grpcClient, err = NewClient(grpccfg); err != nil {
		return
	}
	d = &dao{
		db:              db,
		redis:           r,
		mc:              mc,
		remotesvrClient: grpcClient,
		cache:           fanout.New("cache"),
		demoExpire:      int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}

// SayHello say hello.
func (d *dao) SayHello(c context.Context, req *callapi.HelloReq) (resp *grpcempty.Empty, err error) {
	var svrReq *svrapi.Req
	name := req.Name
	svrReq = new(svrapi.Req)
	fmt.Print("hello origin name is ", name)
	svrReq.Name = fmt.Sprint("Do it yourself", name)
	fmt.Print("hello new name is \n", svrReq.Name)

	//req1 := new(svrapi.Req)
	//req1.Name = "grpc goofly" // If we don't provide req1 name, it will be error, since service will check name field!
	//	r1, err1 := d.remotesvrClient.Create(context.Background(), req1)
	r1, err := d.remotesvrClient.Create(c, svrReq)
	if err != nil {
		fmt.Println(err)
		err = errors.Wrapf(err, "%v", svrReq.Name)
		return
	}
	fmt.Println(r1.Content)

	return
}
