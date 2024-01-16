package pgbackup

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	pathPSQL   = "psql"
	pathPGDump = "pg_dump"
	// pathPGRestore = "pg_restore"
)

type Dump struct {
	UseGzip bool
	OutPath string
}

func NewDump(opath string) *Dump {
	var x = &Dump{
		OutPath: opath,
		UseGzip: true,
	}
	return x
}

func (d *Dump) Dump(targetUrl string) (string, error) {
	config, err := parseUrl(targetUrl)
	if nil != err {
		return "", err
	}

	var args = config.GetArgs()
	// args = append(args, "-")
	// fmt.Println("Args: ", strings.Join(args, " "))
	buff := bytes.NewBuffer(nil)

	cmd := exec.Command(pathPGDump, args...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+config.Password)
	cmd.Stderr = os.Stderr
	cmd.Stdout = buff

	err = cmd.Run()
	if nil != err {
		log.Println("Error: ", err)
		return "", err
	}

	return d.writeFile(buff, config)
}

func (d *Dump) writeFile(buff *bytes.Buffer, config *Config) (string, error) {
	var now = time.Now()
	var folder = fmt.Sprintf("%s/%d/%02d/%02d", d.OutPath, now.Year(), now.Month(), now.Day())
	err := os.MkdirAll(folder, 0777)
	if nil != err {
		return "", err
	}

	var filePath = fmt.Sprintf("%s/%s.sql", folder, config.Database)
	if d.UseGzip {
		filePath = fmt.Sprintf("%s/%s.gz", folder, config.Database)
		compressed, err := gzipCompress(buff.Bytes())
		if nil != err {
			return "", err
		}
		err = os.WriteFile(filePath, compressed.Bytes(), 0666)
		if nil != err {
			return "", err
		}
	} else {
		err = os.WriteFile(filePath, buff.Bytes(), 0666)
		if nil != err {
			return "", err
		}
	}

	return filePath, nil

}

// type GZipEncoder struct {
// }

func gzipCompress(data []byte) (*bytes.Buffer, error) {
	var wbuff = bytes.NewBuffer(nil)

	w, err := gzip.NewWriterLevel(wbuff, gzip.BestCompression)
	if nil != err {
		return nil, err
	}

	_, err = w.Write(data)
	if nil != err {
		return nil, err
	}

	w.Flush()
	w.Close()

	return wbuff, nil
}

func gzipDecompress(buff io.Reader) (*bytes.Buffer, error) {
	r, err := gzip.NewReader(buff)
	if nil != err {
		return nil, err
	}
	defer r.Close()

	var rbuff = bytes.NewBuffer(nil)
	_, err = io.Copy(rbuff, r)
	if nil != err {
		return nil, err
	}

	return rbuff, nil
}
