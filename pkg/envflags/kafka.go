package envflags

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	// WithSasl enables everything deemed necessary to work with SASL
	WithSasl = &withSaslOption{}

	_ KafkaOption = WithSasl
)

// KafkaOption represents applying a set of options to flags
type KafkaOption interface {
	Apply(*KafkaFlags)
}

type withSaslOption struct{}

// Apply adds all sasl. flags for Kafka
func (o *withSaslOption) Apply(f *KafkaFlags) {
	f.String("sasl.mechanisms", "SASL mechanism to use for authentication. Supported: GSSAPI, PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, OAUTHBEARER. NOTE: Despite the name only one mechanism must be configured.")
	f.String("sasl.username", "SASL username for use with the PLAIN and SASL-SCRAM-.. mechanisms")
	f.String("sasl.password", "SASL password for use with the PLAIN and SASL-SCRAM-.. mechanisms")
}

// KafkaFlags contains a set of flags for configuring Kafka.
type KafkaFlags struct {
	configMap kafka.ConfigMap
	flags     map[string]*string
	Logger    *logging.Logger
}

// ForKafka creates a helper object, which allows to create flags
// connected to Kafka configuration.
// Uses passed in Kafka configuration as a baseline.
func ForKafka(configMap kafka.ConfigMap, opts ...KafkaOption) *KafkaFlags {
	f := &KafkaFlags{
		configMap: configMap,
		flags:     map[string]*string{},
	}
	for _, opt := range opts {
		opt.Apply(f)
	}
	return f
}

// String creates String flag
func (kf *KafkaFlags) String(key, usage string) {
	s := new(string)
	kf.flags[key] = s
	environmentVariable := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	flag.StringVar(s, strings.ReplaceAll(key, ".", "-"), os.Getenv(environmentVariable), usage)
}

// StringWithDefault create String flag with default value
func (kf *KafkaFlags) StringWithDefault(key, def, usage string) {
	s := new(string)
	kf.flags[key] = s
	environmentVariable := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	flag.StringVar(s, strings.ReplaceAll(key, ".", "-"), GetStringDefault(environmentVariable, def), usage)
}

// ToConfigMap uses values from created flags and builds a Kafka configuration.
func (kf *KafkaFlags) ToConfigMap() *kafka.ConfigMap {
	cm := kf.configMap

	for k, v := range kf.flags {
		if *v == "" && kf.Logger != nil {
			kf.Logger.Log(logging.Entry{
				Severity: logging.Debug,
				Payload:  fmt.Sprintf("Setting %s set to empty string", k),
			})
		}
		cm[k] = *v
	}

	return &cm
}
