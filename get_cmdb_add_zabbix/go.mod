module get_cmdb_add_zabbix

go 1.12

replace github.com/zabbixoperating => ./zabbixoperating

replace github.com/auth => ./auth

require github.com/zabbixoperating v0.0.0-00010101000000-000000000000

require (
	github.com/chenhg5/collection v0.0.0-20191118032303-cb21bccce4c3
	github.com/shopspring/decimal v0.0.0-20200227202807-02e2044944cc // indirect
)
