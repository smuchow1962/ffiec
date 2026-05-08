module github.com/mmpworks/ffiec/ledger

go 1.23

require (
	github.com/mmpworks/ffiec/core v0.0.0
	github.com/mmpworks/ffiec/testkit v0.0.0
)

replace (
	github.com/mmpworks/ffiec/core => ../core
	github.com/mmpworks/ffiec/testkit => ../testkit
)
