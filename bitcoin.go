package merry

import "os/exec"

func ForkBitcoin() error {
	bashCmd := exec.Command("docker", "exec", "bitcoin", "bitcoind", "-port=18445", "-rpcport=8333", "-datadir=/data/bitcoin-backup", "-daemon")
	if err := bashCmd.Run(); err != nil {
		return err
	}
	// docker exec bitcoin bitcoind -port=18445 -rpcport=8333 -datadir=/data/bitcoin-backup -daemon
	// docker exec bitcoin bitcoin-cli addnode "127.0.0.1:18445" add
	bashCmd2 := exec.Command("docker", "exec", "bitcoin", "bitcoind", "-port=18445", "-rpcport=8333", "-datadir=/data/bitcoin-backup", "-daemon")
	if err := bashCmd2.Run(); err != nil {
		return err
	}
	// docker exec bitcoin bitcoin-cli --rpcport=8333 getaddednodeinfo
	return nil
	// start bitcoin node
	// connect to our bitcoin node
	// sync nodes
	// disconnect nodes
	// send transaction to the old node
	// mine 5 blocks on the new node
	// add node and sync nodes
	// stop node
}

func ResetBitcoin() {
	// start bitcoin node
	// connect to our bitcoin node
	// sync nodes
	// disconnect nodes
	// send transaction to the old node
	// mine 5 blocks on the new node
	// add node and sync nodes
	// stop node
}
