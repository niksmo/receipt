package config

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	defaultLogLevel       = "info"
	defaultHTTPServerAddr = "localhost:8080"

	defaultTopic         = "mail-receipt"
	minPartitions        = 1
	minReplicationFactor = -1
)

var (
	defaultSeedBrokers = []string{
		"localhost:19094", "localhost:29094", "localhost:39094",
	}
)

var errEnvNotSet = errors.New("the env variable is not specified")

type BrokerConfig struct {
	SeedBrokers       []string
	Topic             string
	Partitions        int
	ReplicationFactor int
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

func (c Config) Print(w io.Writer) {
	fmt.Fprintf(w,
		`Configuration:
LogLevel:          %q
HTTPServerAddress: %q
SeedBrokers:       %q
Topic:             %q
Partitions:        % d
ReplicationFactor: % d

`,
		c.LogLevel,
		c.HTTPServerAddr,
		c.BrokerConfig.SeedBrokers,
		c.BrokerConfig.Topic,
		c.BrokerConfig.Partitions,
		c.BrokerConfig.ReplicationFactor,
	)
}

func loadLogLevel() string {
	logLevel, err := envString("RECEIPT_LOG_LEVEL", nil)
	if errors.Is(err, errEnvNotSet) {
		return defaultLogLevel
	}
	return logLevel
}

func loadHTTPServerAddr() (string, error) {
	httpSrvAddr, err := envString(
		"RECEIPT_HTTP_ADDR",
		func(v string) error {
			_, err := net.ResolveTCPAddr("tcp", v)
			return err
		},
	)

	if err != nil {
		if errors.Is(err, errEnvNotSet) {
			return defaultHTTPServerAddr, nil
		}
		return "", err
	}
	return httpSrvAddr, nil
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
	}

	return brokerCfg, nil
}

func loadSeedBrokers() ([]string, error) {
	v, err := envStringS(
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
		if errors.Is(err, errEnvNotSet) {
			return defaultSeedBrokers, nil
		}
		return nil, err
	}

	return v, nil
}

func loadTopic() string {
	v, err := envString("RECEIPT_TOPIC", nil)
	if errors.Is(err, errEnvNotSet) {
		return defaultTopic
	}
	return v
}

func loadPartitions() int {
	v, err := envInt(
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
	v, err := envInt(
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

func errsOnLoad(errs []error) bool {
	return len(errs) != 0
}

func envString(name string, validationFunc func(v string) error) (string, error) {
	v, set := os.LookupEnv(name)
	if !set {
		return "", errEnvNotSet
	}
	if validationFunc != nil {
		if err := validationFunc(v); err != nil {
			return "", err
		}
	}

	return v, nil
}

func envInt(name string, validationFunc func(v int) error) (int, error) {
	vStr, set := os.LookupEnv(name)
	if !set {
		return 0, errEnvNotSet
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, fmt.Errorf("env %s invalid 'int' value: %w", name, err)
	}

	if validationFunc != nil {
		if err := validationFunc(v); err != nil {
			return 0, err
		}
	}

	return v, nil
}

func envStringS(name string, validationFunc func(v []string) error) ([]string, error) {
	v, set := os.LookupEnv(name)
	if !set {
		return nil, errEnvNotSet
	}

	s := strings.Split(v, ",")

	if validationFunc != nil {
		if err := validationFunc(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
