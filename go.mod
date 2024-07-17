module github.com/gopi-frame/collection

go 1.22.2

//replace (
//	github.com/gopi-frame/contract => ../contract
//	github.com/gopi-frame/contract/exception => ../contract/exception
//	github.com/gopi-frame/future => ../future
//	github.com/gopi-frame/utils => ../utils
//)

require (
	github.com/gopi-frame/contract v0.0.0-20240715070535-dae72b3c8cd0
	github.com/gopi-frame/exception v0.0.0-20240628085057-b605370ef1c5
	github.com/gopi-frame/future v0.0.0-20240715070619-7c24479cfa9b
	github.com/gopi-frame/utils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
