module github.com/azhai/gozzo-utils

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190820162420-60c769a6c586
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190731235908-ec7cb31e5a56
	golang.org/x/image => github.com/golang/image v0.0.0-20190823064033-3a9bac650e44
	golang.org/x/net => github.com/golang/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190822191935-b1e2c8edcefd
)

require (
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/kardianos/service v1.0.0
	github.com/kellydunn/golang-geo v0.7.0
	github.com/kylelemons/go-gypsy v0.0.0-20160905020020-08cad365cd28 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/streadway/amqp v0.0.0-20190815230801-eade30b20f1d
	github.com/stretchr/testify v1.4.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
)

go 1.13
