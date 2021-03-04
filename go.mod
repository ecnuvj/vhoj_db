module github.com/ecnuvj/vhoj_db

go 1.14

require (
	github.com/ecnuvj/vhoj_common v0.0.0-00010101000000-000000000000
	github.com/jinzhu/gorm v1.9.16
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/ecnuvj/vhoj_common => ../vhoj_common
