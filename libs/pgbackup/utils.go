package pgbackup

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Username            string
	Password            string
	Host                string
	Port                string
	Database            string
	TablePattern        string
	TablePatternExclude string
}

func NewConfigFromUrl(URL string) (*Config, error) {
	return parseUrl(URL)
}

func (c *Config) GetArgs() []string {
	var args = []string{
		"-h", c.Host,
		"-U", c.Username,
		"-p", c.Port,
		"-d", c.Database,
	}

	if c.TablePattern != "" {
		args = append(args, "--table", c.TablePattern)
	}

	if c.TablePatternExclude != "" {
		args = append(args, "--exclude-table", c.TablePatternExclude)
	}

	return args
}

func (c *Config) Execute(sqlCmd string) (*bytes.Buffer, error) {
	var args = []string{
		"-h", c.Host,
		"-U", c.Username,
		"-p", c.Port,
		"postgres",
		"-c", fmt.Sprintf(`%s;`, sqlCmd),
	}
	buff := bytes.NewBuffer(nil)

	cmd := exec.Command(pathPSQL, args...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.Password)
	cmd.Stderr = buff
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if nil != err {
		return buff, err
	}
	return buff, nil
}

func (c *Config) CreateDb() error {
	var args = []string{
		"-h", c.Host,
		"-U", c.Username,
		"-p", c.Port,
		"postgres",
		"-c", fmt.Sprintf("CREATE DATABASE %s;", c.Database),
	}
	buff := bytes.NewBuffer(nil)

	cmd := exec.Command(pathPSQL, args...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.Password)
	cmd.Stderr = buff
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	// buff, err := c.Execute("CREATE DATABASE IF NOT EXISTS " + c.Database)
	if nil != err {
		if buff != nil && strings.Contains(buff.String(), "already exists") {
			return nil
		}
		log.Println("Error: ", err)
		return err
	}

	return nil
}

func (c *Config) CreateExtension(xname string) error {
	buff, err := c.Execute("CREATE EXTENSION IF NOT EXISTS  " + xname)
	if nil != err {
		if buff != nil && strings.Contains(buff.String(), "already exists") {
			return nil
		}
		log.Println("Error: ", err)
		return err
	}

	return nil
}

func parseUrl(target string) (*Config, error) {
	urlParsed, err := url.Parse(target)
	if nil != err {
		return nil, err
	}

	var rs = &Config{
		Host: urlParsed.Host,
	}

	rs.Username = urlParsed.User.Username()
	rs.Password, _ = urlParsed.User.Password()
	if rs.Username == "" || rs.Password == "" {
		return nil, errors.New("postgres url missing username or password")
	}

	// Parse host
	idx := strings.Index(urlParsed.Host, ":")
	if idx > 0 {
		rs.Host = urlParsed.Host[:idx]
		rs.Port = urlParsed.Host[idx+1:]
	} else {
		rs.Host = urlParsed.Host
		rs.Port = "5432"
	}

	// Db
	if urlParsed.Path[0] == '/' {
		rs.Database = urlParsed.Path[1:]
	} else {
		rs.Database = urlParsed.Path
	}

	// Table
	rs.TablePattern = urlParsed.Query().Get("table")

	return rs, nil
}
