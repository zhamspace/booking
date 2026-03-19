package model

type Main struct {
	Legend []DictIdName

	AuditResources  []DictIdName
	AuditChangeType []DictIdName
}

type DictIdName struct {
	Id   string
	Name string
}
