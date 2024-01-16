package pgbackup

import "testing"

func TestCreateDB(t *testing.T) {
	var target = "postgres://admin:hellosecret@localhost/aaa"
	config, err := parseUrl(target)
	if nil != err {
		panic(err)
	}

	r := &Restore{}

	r.createDB(config)
}

func TestRestore(t *testing.T) {
	var target = "postgres://admin:hellosecret@localhost:5432/projects"
	var r = &Restore{}

	r.Execute("./static/2024/01/15/iott.gz", target)
	// r.Execute("./static/a.sql", target)
}
