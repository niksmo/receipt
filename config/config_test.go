package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		config := LoadConfig()
		assert.Equal(t, defaultLogLevel, config.LogLevel)
		assert.Equal(t, defaultHTTPServerAddr, config.HTTPServerAddr)
		assert.Equal(t, defaultSeedBrokers, config.BrokerConfig.SeedBrokers)
		assert.Equal(t, defaultTopic, config.BrokerConfig.Topic)
		assert.Equal(t, minPartitions, config.BrokerConfig.Partitions)
		assert.Equal(t, minReplicationFactor, config.BrokerConfig.ReplicationFactor)
		assert.Equal(t, defaultConsumerGroup, config.BrokerConfig.ConsumerGroup)
	})

	t.Run("should_set_values", func(t *testing.T) {
		t.Setenv("RECEIPT_LOG_LEVEL", "myLevel")
		t.Setenv("RECEIPT_HTTP_ADDR", "127.0.0.1:4000")
		t.Setenv("RECEIPT_SEED_BROKERS", "localhost:3001,localhost:3002")
		t.Setenv("RECEIPT_TOPIC", "myTopic")
		t.Setenv("RECEIPT_PARTITIONS", "8")
		t.Setenv("RECEIPT_REPLICATION_FACTOR", "3")
		t.Setenv("RECEIPT_CONSUMER_GROUP", "myGroup")

		config := LoadConfig()
		assert.Equal(t, "myLevel", config.LogLevel)
		assert.Equal(t, "127.0.0.1:4000", config.HTTPServerAddr)
		assert.Equal(t, []string{"localhost:3001", "localhost:3002"}, config.BrokerConfig.SeedBrokers)
		assert.Equal(t, "myTopic", config.BrokerConfig.Topic)
		assert.Equal(t, 8, config.BrokerConfig.Partitions)
		assert.Equal(t, 3, config.BrokerConfig.ReplicationFactor)
		assert.Equal(t, "myGroup", config.BrokerConfig.ConsumerGroup)
	})

	t.Run("should_panic", func(t *testing.T) {
		t.Setenv("RECEIPT_HTTP_ADDR", "notvalidaddr123456")
		t.Setenv("RECEIPT_SEED_BROKERS", "notvalidbrokeraddr1,notvalidbrokeraddr2")

		require.Panics(t, func() {
			LoadConfig()
		})
	})
}
