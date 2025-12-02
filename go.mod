module github.com/mackerelio/mackerel-agent-plugins

go 1.24.0

toolchain go1.24.2

require (
	github.com/Songmu/axslogparser v1.4.0
	github.com/Songmu/postailer v0.0.0-20181014062912-daaa1ba9cc39
	github.com/Songmu/timeout v0.4.0
	github.com/aws/aws-sdk-go-v2 v1.40.0
	github.com/aws/aws-sdk-go-v2/config v1.32.2
	github.com/aws/aws-sdk-go-v2/credentials v1.19.2
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.14
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.52.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.275.0
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5
	github.com/fsouza/go-dockerclient v1.12.3
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/go-ldap/ldap/v3 v3.4.12
	github.com/go-redis/redismock/v9 v9.2.0
	github.com/gosnmp/gosnmp v1.42.1
	github.com/jarcoal/httpmock v1.4.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lestrrat-go/tcptest v0.0.0-20180223004312-f0345789c593
	github.com/lib/pq v1.10.9
	github.com/mackerelio/go-mackerel-plugin v0.1.5
	github.com/mackerelio/go-mackerel-plugin-helper v0.1.4
	github.com/mackerelio/go-osstat v0.2.6
	github.com/mackerelio/golib v1.2.1
	github.com/mackerelio/mackerel-plugin-mongodb v1.1.2
	github.com/mackerelio/mackerel-plugin-mysql v1.3.2
	github.com/mattn/go-pipeline v0.0.0-20190323144519-32d779b32768
	github.com/mattn/go-treasuredata v0.0.0-20170920030233-31758907cfc4
	github.com/michaelklishin/rabbit-hole v1.5.0
	github.com/montanaflynn/stats v0.7.1
	github.com/redis/go-redis/v9 v9.17.0
	github.com/stretchr/testify v1.11.1
	github.com/tomasen/fcgi_client v0.0.0-20180423082037-2bb3d819fd19
	github.com/urfave/cli v1.22.17
	github.com/yusufpapurcu/wmi v1.2.4
	golang.org/x/text v0.31.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Songmu/go-ltsv v0.0.0-20181014062614-c30af2b7b171 // indirect
	github.com/Songmu/wrapcommander v0.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.2 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.8-0.20250403174932-29230038a667 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lestrrat-go/tcputil v0.0.0-20180223003554-d3c7f98154fb // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.1.0 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.3 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
