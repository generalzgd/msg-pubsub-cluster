/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @time: 2019/12/6 4:55 下午
 * @project: packagesubscribesvr
 */

package config

import (
	`os`
	`path/filepath`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/svr-config/ymlcfg`
)

//type ClusterConfig struct {
//	NodeType int      `yaml:"type"`
//	SerfPort int      `yaml:"serf"`
//	RaftPort int      `yaml:"raft"`
//	RpcPort  int      `yaml:"rpc"`
//	HttpPort int      `yaml:"http"`
//	Except   int      `yaml:"except"`
//	Peers    []string `yaml:"peers"` // serf终端地址
//}

// type PostConfig struct {
// 	TcpCfg    svrcfg.TcpCfg `yaml:"link"`
// 	LogLevel  uint8          `yaml:"loglevel"`
// 	ExcludeIp []string      `yaml:"exclude"`
// }

//type GrpcPoolCfg struct {
//	Init           int           `yaml:"init"`
//	Capacity       int           `yaml:"capacity"`
//	IdleTimeout    time.Duration `yaml:"idle"`
//	MaxLifeTimeout time.Duration `yaml:"maxlife"`
//}

type SubscribeCfg struct {
	Host string `yaml:"host"` //
	Port int    `yaml:"port"`
}

type DecodeCfg struct {
	HeadSize int `yaml:"headsize"`
	CmdPos   int `yaml:"cmdpos"`
	CmdSize  int `yaml:"cmdsize"`
	LenPos   int `yaml:"lenpos"`
	LenSize  int `yaml:"lensize"`
}

type AppConfig struct {
	Name         string        `yaml:"name"`
	Ver          string        `yaml:"ver"`
	Memo         string        `yaml:"memo"`
	LogLevel     int           `yaml:"loglevel"`
	PostLogLevel uint8         `yaml:"postloglevel"`
	ExcludeIp    []string      `yaml:"exclude"`
	TreeDegree   int           `yaml:"degree"`
	MaxPackDelay time.Duration `yaml:"maxpackdelay"`
	// 订阅者异常断开，延时清理订阅信息。要大于死信延时
	CleanDelay    time.Duration `yaml:"cleandelay"`
	Retry         int           `yaml:"retry"`
	QueueMode     int           `yaml:"queuemode"`
	BatchSize     int           `yaml:"batch"`
	DeadBatchSize int           `yaml:"deadbatch"`
	// 死信队列运行间隔
	DeadDelay    time.Duration        `yaml:"deaddelay"`
	Consul       ymlcfg.ConsulConfig  `yaml:"consul"`
	PostCfg      ymlcfg.TcpCfg        `yaml:"post"`
	SubscribeCfg SubscribeCfg         `yaml:"subscribe"`
	Cluster      ymlcfg.ClusterConfig `yaml:"cluster"`
	GrpcPool     ymlcfg.ConnPool      `yaml:"grpcpool"`
	Decode       DecodeCfg            `yaml:"decode"`
}

func (p *AppConfig) GetLogLevel() int {
	if p.LogLevel > 0 {
		return p.LogLevel
	}
	return logs.LevelInfo
}

func (p *AppConfig) Load() error {
	path := filepath.Join(filepath.Dir(os.Args[0]), "config", "config_dev.yml")
	logs.Info("load config:", path)
	return ymlcfg.LoadYaml(path, p)
}
