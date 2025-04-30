package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// 传统基于权重的负载均衡算法

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {
		cc := &conn{
			cc: sc,
		}

		md := sci.Address.Attributes.WithValue("weight", 10)

		weightVal := md.Value("weight")
		weight, _ := weightVal.(float64)
		cc.weight = int(weight)

		if cc.weight == 0 {
			cc.weight = 10
		}

		conns = append(conns, cc)
	}
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*conn
}

// Pick 实现基于权重的负载均衡算法
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable

}

// 代表节点
type conn struct {
	weight int
	cc     balancer.SubConn
}
