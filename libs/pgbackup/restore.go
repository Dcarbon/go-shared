package pgbackup

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Restore struct {
}

func (r *Restore) Execute(filename, target string) error {
	// data, err := r.readDump(filename)
	// if nil != err {
	// 	return err
	// }

	filename, err := r.unCompress(filename)
	if nil != err {
		return err
	}

	config, err := parseUrl(target)
	if nil != err {
		return err
	}

	err = r.createDB(config)
	if nil != err {
		return err
	}

	return r.restoreFrom(config, filename)
}

// func (r *Restore) readDump(filename string) ([]byte, error) {
// 	data, err := os.ReadFile(filename)
// 	if nil != err {
// 		return nil, err
// 	}
// 	if len(filename) > 3 && filename[len(filename)-3:] == ".gz" {
// 		decompressed, err := gzipDecompress(bytes.NewBuffer(data))
// 		if nil != err {
// 			return nil, err
// 		}
// 		data = decompressed.Bytes()
// 	}
// 	// fmt.Println(string(data))
// 	return data, nil
// }

func (r *Restore) unCompress(file string) (string, error) {
	if filepath.Ext(file) == ".gz" {
		cmd := exec.Command("gunzip", "-f", "-k", file)
		err := cmd.Run()
		if nil != err {
			return "", err
		}
		return file[:len(file)-3], nil
	}
	return file, nil
}

//	CREATE EXTENSION IF NOT EXISTS postgis;
//
// psql -U $DB_USER -h $DB_HOST -p $DB_PORT postgres -c "CREATE DATABASE $RESTORE_DATABASE;"
func (r *Restore) createDB(c *Config) error {
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
	if nil != err {
		if strings.Contains(buff.String(), "already exists") {
			return nil
		}
		log.Println("Error: ", err)
		return err
	}

	return nil
}

func (r *Restore) restoreFrom(c *Config, file string) error {
	var args = c.GetArgs()
	args = append(args, "-f", file)
	log.Println(strings.Join(args, " "))

	cmd := exec.Command(pathPSQL, args...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.Password)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if nil != err {
		return err
	}
	return nil
}

// func (r *Restore) restore(c *Config, data []byte) error {
// 	var args = c.GetArgs()
// 	// args = append(args, "-c", fmt.Sprintf(`"%s"`, data))
// 	log.Println(strings.Join(args, " "))
// 	args = append(args, "-f", string(data))
// 	// buff := bytes.NewBuffer(nil)
// 	cmd := exec.Command(pathPSQL, args...)
// 	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.Password)
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdout = os.Stdout
// 	// cmd.Stdin = buff
// 	// buff.Write(data)
// 	err := cmd.Run()
// 	if nil != err {
// 		return err
// 	}
// 	return nil
// }
