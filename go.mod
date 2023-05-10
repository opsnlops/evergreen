module github.com/evergreen-ci/evergreen

go 1.19

require (
	github.com/99designs/gqlgen v0.17.20
	github.com/PuerkitoBio/rehttp v1.1.0
	github.com/aws/aws-sdk-go v1.44.210
	github.com/cheynewallace/tabby v1.1.1
	github.com/docker/docker v20.10.12+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/dustin/go-humanize v1.0.0
	github.com/evergreen-ci/birch v0.0.0-20220401151432-c792c3d8e0eb
	github.com/evergreen-ci/certdepot v0.0.0-20211117185134-dbedb3d79a10
	github.com/evergreen-ci/cocoa v0.0.0-20230323210713-a34b158b096e
	github.com/evergreen-ci/gimlet v0.0.0-20230317174941-7e759337ad2c
	github.com/evergreen-ci/juniper v0.0.0-20230119161755-1aced8006202
	github.com/evergreen-ci/pail v0.0.0-20220908201135-8a2090a672b7
	github.com/evergreen-ci/poplar v0.0.0-20220908212406-a5e2aa799def
	github.com/evergreen-ci/shrub v0.0.0-20211025143051-a8d91b2e29fd
	github.com/evergreen-ci/timber v0.0.0-20230413164224-05a3a6e11d78
	github.com/evergreen-ci/utility v0.0.0-20230508194838-873894227dc9
	github.com/google/go-github/v52 v52.0.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gophercloud/gophercloud v0.1.0
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/jpillora/backoff v1.0.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mongodb/amboy v0.0.0-20221207220239-4ab00e3ea9da
	github.com/mongodb/anser v0.0.0-20230501213745-c62f11870fd4
	github.com/mongodb/grip v0.0.0-20230503175150-8b2e395f3111
	github.com/pkg/errors v0.9.1
	github.com/ravilushqa/otelgqlgen v0.11.0
	github.com/robbiet480/go.sns v0.0.0-20210223081447-c7c9eb6836cb
	github.com/robfig/cron v1.2.0
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/smartystreets/goconvey v1.7.2
	github.com/stretchr/testify v1.8.2
	github.com/urfave/cli v1.22.10
	github.com/vektah/gqlparser/v2 v2.5.1
	github.com/vmware/govmomi v0.27.1
	go.mongodb.org/mongo-driver v1.11.2
	go.opentelemetry.io/contrib/detectors/aws/ec2 v1.15.0
	go.opentelemetry.io/contrib/detectors/aws/ecs v1.15.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.40.0
	go.opentelemetry.io/otel v1.15.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.15.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.15.0
	go.opentelemetry.io/otel/metric v0.38.0
	go.opentelemetry.io/otel/sdk v1.15.0
	go.opentelemetry.io/otel/sdk/metric v0.38.0
	go.opentelemetry.io/otel/trace v1.15.0
	golang.org/x/crypto v0.7.0
	golang.org/x/oauth2 v0.7.0
	golang.org/x/tools v0.6.0
	gonum.org/v1/gonum v0.12.0
	google.golang.org/api v0.103.0
	google.golang.org/grpc v1.54.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go/compute v1.15.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20220621081337-cb9428e4ac1e // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/andygrunwald/go-jira v1.14.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/brunoscheufler/aws-ecs-metadata-go v0.0.0-20220812150832-b6b31c6eeeaf // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cloudflare/circl v1.1.0 // indirect
	github.com/containerd/cgroups v1.0.3 // indirect
	github.com/containerd/containerd v1.5.18 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/dghubble/oauth1 v0.7.2 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/evergreen-ci/aviation v0.0.0-20220405151811-ff4a78a4297c // indirect
	github.com/evergreen-ci/baobab v1.0.1-0.20211025210153-3206308845c1 // indirect
	github.com/evergreen-ci/bond v0.0.0-20211109152423-ba2b6b207f56 // indirect
	github.com/evergreen-ci/lru v0.0.0-20211029170532-008d075b972d // indirect
	github.com/evergreen-ci/mrpc v0.0.0-20211025143107-842bca81a3f8 // indirect
	github.com/evergreen-ci/negroni v1.0.1-0.20211028183800-67b6d7c2c035 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/fuyufjh/splunk-hec-go v0.3.4-0.20210909061418-feecd03924b7 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-ldap/ldap/v3 v3.4.4 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/goccy/go-json v0.9.4 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.12.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/jwx v1.2.18 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/mattn/go-xmpp v0.0.0-20211029151415-912ba614897a // indirect
	github.com/mholt/archiver/v3 v3.5.1 // indirect
	github.com/mongodb/ftdc v0.0.0-20220401165013-13e4af55e809 // indirect
	github.com/montanaflynn/stats v0.0.0-20180911141734-db72e6cae808 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/okta/okta-jwt-verifier-golang v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/papertrail/go-tail v0.0.0-20180509224916-973c153b0431 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/peterhellberg/link v1.2.0 // indirect
	github.com/phyber/negroni-gzip v1.0.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/rs/cors v1.8.3 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.4 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/slack-go/slack v0.12.1 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/square/certstrap v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/trivago/tgo v1.0.7 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/urfave/cli/v2 v2.8.1 // indirect
	github.com/urfave/negroni v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.15.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.38.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

require (
	github.com/mholt/archiver v2.0.1-0.20180417220235-e4ef56d48eb0+incompatible
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
)

require (
	github.com/evergreen-ci/evg-lint v0.0.0-20211115144425-3b19c8e83a57
	github.com/evergreen-ci/plank v0.0.0-20230207190607-5f47f8a30da1
	github.com/evergreen-ci/tarjan v0.0.0-20170824211642-fcd3f3321826
	github.com/mongodb/jasper v0.0.0-20220214215554-82e5a72cff6b
	github.com/shirou/gopsutil/v3 v3.23.3
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gopkg.in/yaml.v3 v3.0.1
)
