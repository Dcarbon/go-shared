package pgbackup

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"testing"
)

func TestParseUrl(t *testing.T) {
	URL := "postgres://admin:hellosecret@postgis/aaaa"
	parseUrl(URL)
}

func TestDump(t *testing.T) {
	var opath = "./static"
	var x = NewDump(opath)
	var file, err = x.Dump("postgres://admin:244466666@10.60.0.58/iott?table=projects_specs")
	if nil != err {
		panic(err)
	}
	log.Println("Filename: ", file)
}

func TestDumpAndRestore(t *testing.T) {
	var opath = "./static"
	var dumpUrls = []string{
		"postgres://admin:244466666@10.60.0.58/iott?table=iot_shape",
		// "postgres://admin:244466666@10.60.0.58/iott?table=projects_specs",
	}

	var targetUrl = "postgres://admin:hellosecret@13.228.11.143:5432/iot_shape"
	for _, durl := range dumpUrls {
		var x = NewDump(opath)
		var file, err = x.Dump(durl)
		if nil != err {
			panic(err)
		}
		log.Println("Filename: ", file)

		var r = &Restore{}
		r.Execute(file, targetUrl)
	}
}

func TestGzip(t *testing.T) {
	var txt = "this is content to zip"
	var wbuff = bytes.NewBuffer(nil)
	var w = gzip.NewWriter(wbuff)
	_, err := w.Write([]byte(txt))
	if nil != err {
		panic(err)
	}
	w.Flush()
	w.Close()

	log.Println("GZip writer success")

	r, err := gzip.NewReader(wbuff)
	if nil != err {
		panic(err)
	}

	var rbuff = bytes.NewBuffer(nil)
	io.Copy(rbuff, r)
	r.Close()

	log.Println("GZip reade success: ", rbuff.String())
}

func TestGzip2(t *testing.T) {
	var txt = "this is content to zip 2"
	compressed, err := gzipCompress([]byte(txt))
	if nil != err {
		panic(err)
	}

	decompressed, err := gzipDecompress(compressed)
	if nil != err {
		panic(err)
	}

	log.Println("GZip reade success: ", decompressed.String())
}
