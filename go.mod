module github.com/azhai/gozzo-utils

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190313024323-a1f597ede03a
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190312203227-4b39c73a6495
	golang.org/x/image => github.com/golang/image v0.0.0-20190227222117-0694c2d4d067
	golang.org/x/net => github.com/golang/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190312061237-fead79001313
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190312170243-e65039ee4138
)

require (
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/kardianos/service v1.0.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/streadway/amqp v0.0.0-20190312223743-14f78b41ce6d
	github.com/stretchr/testify v1.3.0 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
)
