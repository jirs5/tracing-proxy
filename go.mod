module github.com/jirs5/tracing-proxy

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/facebookgo/startstop v0.0.0-20161013234910-bc158412526d
	github.com/facebookgo/structtag v0.0.0-20150214074306-217e25fb9691 // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.3
	github.com/gomodule/redigo v1.8.8
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/honeycombio/dynsampler-go v0.2.1
	github.com/honeycombio/husky v0.9.0
	github.com/honeycombio/libhoney-go v1.15.8
	github.com/jessevdk/go-flags v1.5.0
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.13.6
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.26.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	github.com/vmihailenco/msgpack/v4 v4.3.11
	go.opentelemetry.io/proto/otlp v0.9.0
	google.golang.org/grpc v1.50.1
	gopkg.in/alexcesaro/statsd.v2 v2.0.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

//replace github.com/honeycombio/libhoney-go v1.15.8 => github.com/jirs5/libtrace-go v0.0.0-20220209113356-39ae92fc19f4
replace github.com/honeycombio/libhoney-go v1.15.8 => github.com/jirs5/libtrace-go v1.15.9-0.20221219105703-796cb39b0512

//replace github.com/honeycombio/libhoney-go v1.15.8 => github.com/jirs5/libtrace-go v1.15.9-0.20221215130906-ffb6698e9c86

//replace github.com/honeycombio/husky v0.9.0 => github.com/jirs5/husky v0.9.1-0.20220302161820-fe16f58d3996
replace github.com/honeycombio/husky v0.9.0 => github.com/jirs5/husky v0.9.1-0.20220616112458-7bb2625f28df
