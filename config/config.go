package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/niksmo/receipt/pkg/env"
)

const (
	defaultLogLevel       = "info"
	defaultHTTPServerAddr = ":8080"

	defaultTopic         = "mail-receipt"
	minPartitions        = 1
	minReplicationFactor = -1

	defaultConsumerGroup = "mail-group"
)

var (
	defaultSeedBrokers = []string{
		"localhost:19094", "localhost:29094", "localhost:39094",
	}
)

type BrokerConfig struct {
	SeedBrokers       []string
	Topic             string
	Partitions        int
	ReplicationFactor int
	ConsumerGroup     string
}

type Config struct {
	LogLevel       string
	HTTPServerAddr string
	BrokerConfig
}

func LoadConfig() Config {
	var errs []error

	httpSrvAddr, err := loadHTTPServerAddr()
	if err != nil {
		errs = append(errs, err)
	}

	brokerCfg, err := loadBrokerConfig()
	if err != nil {
		errs = append(errs, err)
	}

	if errsOnLoad(errs) {
		panic(errors.Join(errs...))
	}

	cfg := Config{
		LogLevel:       loadLogLevel(),
		HTTPServerAddr: httpSrvAddr,
		BrokerConfig:   brokerCfg,
	}
	return cfg
}

func (c Config) Print() {
	fmt.Printf(
		`Configuration:
LogLevel:          %q
HTTPServerAddress: %q
SeedBrokers:       %q
Topic:             %q
Partitions:        % d
ReplicationFactor: % d
ConsumerGroup:     %q

`,
		c.LogLevel,
		c.HTTPServerAddr,
		c.SeedBrokers,
		c.Topic,
		c.Partitions,
		c.ReplicationFactor,
		c.ConsumerGroup,
	)
}

func LoadLogLevel(envValue string, defaultValue string) string {
	logLevel, err := env.String(envValue, nil)
	if errors.Is(err, env.ErrNotSet) {
		return defaultValue
	}
	return logLevel
}

func LoadHTTPServerAddr(envValue string, defaultValue string) (string, error) {
	httpSrvAddr, err := env.String(
		envValue,
		func(v string) error {
			_, err := net.ResolveTCPAddr("tcp", v)
			return err
		},
	)

	if err != nil {
		if errors.Is(err, env.ErrNotSet) {
			return defaultValue, nil
		}
		return "", err
	}
	return httpSrvAddr, nil
}

func loadLogLevel() string {
	return LoadLogLevel("RECEIPT_LOG_LEVEL", defaultLogLevel)
}

func loadHTTPServerAddr() (string, error) {
	return LoadHTTPServerAddr("RECEIPT_HTTP_ADDR", defaultHTTPServerAddr)
}

func loadBrokerConfig() (BrokerConfig, error) {
	var errs []error

	seedBrokers, err := loadSeedBrokers()
	if err != nil {
		errs = append(errs, err)
	}

	if errsOnLoad(errs) {
		return BrokerConfig{}, errors.Join(errs...)
	}

	brokerCfg := BrokerConfig{
		SeedBrokers:       seedBrokers,
		Topic:             loadTopic(),
		Partitions:        loadPartitions(),
		ReplicationFactor: loadReplicationFactor(),
		ConsumerGroup:     loadConsumerGroup(),
	}

	return brokerCfg, nil
}

func loadSeedBrokers() ([]string, error) {
	v, err := env.StringS(
		"RECEIPT_SEED_BROKERS",
		func(v []string) error {
			for _, brokerAddr := range v {
				_, err := net.ResolveTCPAddr("tcp", brokerAddr)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

	if err != nil {
		if errors.Is(err, env.ErrNotSet) {
			return defaultSeedBrokers, nil
		}
		return nil, err
	}

	return v, nil
}

func loadTopic() string {
	v, err := env.String("RECEIPT_TOPIC", nil)
	if errors.Is(err, env.ErrNotSet) {
		return defaultTopic
	}
	return v
}

func loadPartitions() int {
	v, err := env.Int(
		"RECEIPT_PARTITIONS",
		func(v int) error {
			if v < minPartitions {
				return errors.New("invalid number of partitions")
			}
			return nil
		},
	)

	if err != nil {
		return minPartitions
	}

	return v
}

func loadReplicationFactor() int {
	v, err := env.Int(
		"RECEIPT_REPLICATION_FACTOR",
		func(v int) error {
			if v < minReplicationFactor {
				return errors.New("invalid replication factor")
			}
			return nil
		},
	)
	if err != nil {
		return minReplicationFactor
	}

	return v
}

func loadConsumerGroup() string {
	v, err := env.String("RECEIPT_CONSUMER_GROUP", nil)
	if errors.Is(err, env.ErrNotSet) {
		return defaultConsumerGroup
	}
	return v
}

func errsOnLoad(errs []error) bool {
	return len(errs) != 0
}
