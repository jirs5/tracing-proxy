//go:build all || race
// +build all race

package sharder

import (
	"testing"

	"github.com/jirs5/tracing-proxy/config"
	"github.com/jirs5/tracing-proxy/internal/peer"
	"github.com/jirs5/tracing-proxy/logger"
	"github.com/stretchr/testify/assert"
)

func TestWhichShard(t *testing.T) {
	const (
		selfAddr = "127.0.0.1:8081"
		traceID  = "test"
	)

	peers := []string{
		"http://" + selfAddr,
		"http://2.2.2.2:8081",
		"http://3.3.3.3:8081",
	}
	config := &config.MockConfig{
		GetPeerListenAddrVal: selfAddr,
		GetPeersVal:          peers,
		PeerManagementType:   "file",
	}
	filePeers, err := peer.NewPeers(config)
	assert.Equal(t, nil, err)
	sharder := DeterministicSharder{
		Config: config,
		Logger: &logger.NullLogger{},
		Peers:  filePeers,
	}

	assert.NoError(t, sharder.Start(),
		"starting deterministic sharder should not error")

	shard := sharder.WhichShard(traceID)
	assert.Contains(t, peers, shard.GetAddress(),
		"should select a peer for a trace")

	config.GetPeersVal = []string{}
	config.ReloadConfig()
	assert.Equal(t, shard.GetAddress(), sharder.WhichShard(traceID).GetAddress(),
		"should select the same peer if peer list becomes empty")
}
