module github.com/anyproto/anytype-cli

go 1.23.8

require (
	github.com/anyproto/anytype-heart v0.42.0
	github.com/cheggaaa/mb/v3 v3.0.2
	github.com/chzyer/readline v1.5.1
	github.com/spf13/cobra v1.9.1
	github.com/zalando/go-keyring v0.2.6
	google.golang.org/grpc v1.74.2
)

require (
	al.essio.dev/pkg/shellescape v1.6.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/anyproto/any-store v0.3.3 // indirect
	github.com/anyproto/any-sync v0.28.0 // indirect
	github.com/anyproto/go-slip10 v1.0.0 // indirect
	github.com/anyproto/go-slip21 v1.0.0 // indirect
	github.com/anyproto/go-sqlite v1.4.2-any // indirect
	github.com/btcsuite/btcd v0.24.2 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.5 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20220129005943-27c39e0ab4f9 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ipfs/go-cid v0.5.0 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-libp2p v0.42.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mb0/diff v0.0.0-20131118162322-d8d9a906c24d // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.16.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.9.2 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20250313105119-ba97887b0a25 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/samber/lo v1.51.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/zeebo/errs v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.7 // indirect
	gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
	modernc.org/libc v1.66.6 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.38.2 // indirect
	storj.io/drpc v0.0.34 // indirect
)

replace (
	github.com/anyproto/any-sync => github.com/anyproto/any-sync v0.9.2
	github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.1
	github.com/btcsuite/btcutil => github.com/btcsuite/btcd/btcutil v1.1.5
	gopkg.in/Graylog2/go-gelf.v2 => github.com/anyproto/go-gelf v0.0.0-20210418191311-774bd5b016e7
)
