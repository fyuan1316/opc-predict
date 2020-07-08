package option

import "opcdata-predict/pkg/scopelog"

type OptionFunc func(spec *Options)

type Options struct {
	// predict
	PredictDomain  string
	PredictAuth    string
	PredictIp      string
	PredictTimeout int // ms
	// kafka
	Brokers       string
	Topic         string
	Partition     int
	SaslEnable    bool
	SaslUsername  string
	SaslPassword  string
	SaslMechanism string
	//ws server
	WsPort int
}

const MinTimeout = 200 // the unit is milliseconds

func SetWsPort(p int) OptionFunc {
	return func(spec *Options) {
		spec.WsPort = p
	}
}

func SetPredictDomain(domain string) OptionFunc {
	return func(spec *Options) {
		spec.PredictDomain = domain
	}
}
func SetPredictAuth(auth string) OptionFunc {
	return func(spec *Options) {
		spec.PredictAuth = auth
	}
}
func SetServerIp(p string) OptionFunc {
	return func(spec *Options) {
		spec.PredictIp = p
	}
}
func SetTimeout(p int) OptionFunc {
	return func(spec *Options) {
		if p < MinTimeout {
			scopelog.Printf("Options", "predict timeout: %d (ms) is too short, use default %d(ms) instead.", p, MinTimeout)
			p = MinTimeout
		}
		spec.PredictTimeout = p
	}
}
func SetBrokers(p string) OptionFunc {
	return func(spec *Options) {
		spec.Brokers = p
	}
}
func SetTopic(p string) OptionFunc {
	return func(spec *Options) {
		spec.Topic = p
	}
}
func SetPartition(p int) OptionFunc {
	return func(spec *Options) {
		spec.Partition = p
	}
}

func SetSaslEnable(p bool) OptionFunc {
	return func(spec *Options) {
		spec.SaslEnable = p
	}
}
func SetSaslUser(user string) OptionFunc {
	return func(spec *Options) {
		spec.SaslUsername = user
	}
}
func SetSaslPassword(p string) OptionFunc {
	return func(spec *Options) {
		spec.SaslPassword = p
	}
}
func SetSaslMechanism(p string) OptionFunc {
	return func(spec *Options) {
		spec.SaslMechanism = p
	}
}

func NewOptions(fns ...OptionFunc) Options {
	opts := &Options{}
	for _, fn := range fns {
		fn(opts)
	}

	return *opts
}
