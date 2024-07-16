FROM ruimarinho/bitcoin-core:alpine

ADD bitcoin.conf /root/.bitcoin/bitcoin.conf
ADD bitcoin-backup.conf /root/.bitcoin-backup/bitcoin.conf

RUN bitcoind -daemon -datadir=/root/.bitcoin-backup

ENTRYPOINT [ "bitcoind -datadir=/root/.bitcoin" ]