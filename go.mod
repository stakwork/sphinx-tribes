module github.com/stakwork/sphinx-tribes

go 1.2

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/ambelovsky/go-structs v1.1.0 // indirect
	github.com/ambelovsky/gosf v0.0.0-20201109201340-237aea4d6109
	github.com/ambelovsky/gosf-socketio v0.0.0-20220810204405-0f97832ec7af // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2
	github.com/btcsuite/btcwallet/wallet/txauthor v1.3.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/decred/dcrd/lru v1.1.2 // indirect
	github.com/fiatjaf/go-lnurl v1.12.1
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-co-op/gocron v1.25.0
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/gobuffalo/logger v1.0.4 // indirect
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/google/go-github/v39 v39.2.0
	github.com/gorilla/websocket v1.5.0
	github.com/imroc/req v0.3.0
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/joho/godotenv v1.3.0
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/lib/pq v1.10.7
	github.com/lightninglabs/neutrino/cache v1.1.2 // indirect
	github.com/lightningnetwork/lnd v0.17.0-beta.rc6 // indirect
	github.com/lightningnetwork/lnd/tor v1.1.3 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/nbd-wtf/ln-decodepay v1.11.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/redis/go-redis/v9 v9.2.1
	github.com/rs/cors v1.7.0
	github.com/rs/xid v1.4.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/oauth2 v0.4.0
	google.golang.org/api v0.112.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.0
)

replace google.golang.org/api => google.golang.org/api v0.63.0
