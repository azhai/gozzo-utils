module github.com/azhai/gozzo-utils

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190320223903-b7391e95e576
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190316020145-860388717186
	golang.org/x/image => github.com/golang/image v0.0.0-20190321063152-3fc05d484e9f
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190321052220-f7bb7a8bee54
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190320215829-36c10c0a621f
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
