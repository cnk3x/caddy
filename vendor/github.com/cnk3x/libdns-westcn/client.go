package westcn

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cnk3x/k3xurl"
	"github.com/cnk3x/k3xurl/process"
	"github.com/goccy/go-json"
	"github.com/libdns/libdns"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//https://www.west.cn/CustomerCenter/doc/%E8%A5%BF%E9%83%A8%E6%95%B0%E7%A0%81%E5%9F%9F%E5%90%8D%E4%B8%9A%E5%8A%A1API%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3V2.0.html#21u3001u83b7u53d6u8d26u6237u53efu7528u4f59u989d0a3ca20id3d21u3001u83b7u53d6u8d26u6237u53efu7528u4f59u989d3e203ca3e

type Client struct {
	Endpoint string
	Username string
	Password string
}

func (p *Client) sign(data url.Values) url.Values {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	hash := md5.New()
	hash.Write([]byte(p.Username + p.Password + ts))
	token := hex.EncodeToString(hash.Sum(nil))
	data.Set("token", token)
	data.Set("time", ts)
	data.Set("username", p.Username)
	return data
}

func (p *Client) pre(next k3xurl.Process) k3xurl.Process {
	return func(resp *http.Response) error {
		body := transform.NewReader(resp.Body, simplifiedchinese.HZGB2312.NewEncoder())
		var result Result
		err := json.NewDecoder(body).Decode(&result)
		if err != nil {
			return err
		}
		if result.Result != 200 {
			return &result
		}
		nr := new(http.Response)
		*nr = *resp
		nr.Body = io.NopCloser(bytes.NewReader(result.Data))
		nr.ContentLength = int64(len(result.Data))
		return next(nr)
	}
}

func (p *Client) request(ctx context.Context, api string, params url.Values, out any) error {
	proc := k3xurl.NoProcess
	if out != nil {
		proc = process.JSON(out)
	}
	return k3xurl.With(ctx).
		Url(p.Endpoint + api).
		Method(k3xurl.MethodPost).
		Form(p.sign(params)).
		BeforeProcess(p.pre).
		Process(proc)
}

func (p *Client) GetRecords(ctx context.Context, zone string) (records []libdns.Record, err error) {
	const api = "/domain/?act=getdnsrecord"
	params := url.Values{
		"domain": {zone},
		"limit":  {"100"},
		"pageno": {"1"},
	}

	type Record struct {
		ID    int64  `json:"id"`
		Host  string `json:"host"`
		Value string `json:"value"`
		Type  string `json:"type"`
		Level int    `json:"level"`
		TTL   int    `json:"ttl"`
		Line  string `json:"line"`
		Pause int    `json:"pause"`
	}

	var data struct {
		Limit      int      `json:"limit"`
		Total      int      `json:"total"`
		PageNo     int      `json:"pageno"`
		TotalPages int      `json:"totalpages"`
		Items      []Record `json:"items"`
	}

	if err = p.request(ctx, api, params, &data); err != nil {
		return
	}

	for _, item := range data.Items {
		records = append(records, libdns.Record{
			ID:       strconv.FormatInt(item.ID, 10),
			Type:     item.Type,
			Name:     item.Host,
			Value:    item.Value,
			TTL:      time.Second * time.Duration(item.TTL),
			Priority: item.Level,
		})
	}

	return
}

func (p *Client) AddRecord(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	ttl := strconv.FormatFloat(record.TTL.Seconds(), 'f', 0, 0)
	params := url.Values{
		"domain": {zone},                          //	string	是	west.cn	域名
		"host":   {record.Name},                   //	string	是	@	主机名称
		"type":   {record.Type},                   //	string	是	A	解析类型 限 A,CNAME,MX,TXT,AAAA,SRV
		"value":  {record.Value},                  //	string	是	127.0.0.1	解析值
		"ttl":    {ttl},                           //	number	是	900	解析生效时间值 60~86400 单位秒 (默认900)
		"level":  {strconv.Itoa(record.Priority)}, //	number	是	10	优先级别 1-100 默认(10)
		"line":   {""},                            //	string	否		线路: 必须先添加默认线程后才能添加其它线路(默认="" ,电信="LTEL" ,联通="LCNC" ,移动="LMOB" ,教育网="LEDU" ,搜索引擎="LSEO")
	}

	var data struct {
		ID int64 `json:"id"`
	}
	err := p.request(ctx, "/domain/?act=adddnsrecord", params, &data)
	record.ID = strconv.FormatInt(data.ID, 10)
	return record, err
}

func (p *Client) UpdateRecord(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	ttl := strconv.FormatFloat(record.TTL.Seconds(), 'f', 0, 0)
	params := url.Values{
		"domain": {zone},
		"id":     {record.ID},
		"value":  {record.Value},
		"ttl":    {ttl},
		"host":   {record.Name},
		"type":   {record.Type},
		"level":  {strconv.Itoa(record.Priority)},
		"line":   {""},
	}
	err := p.request(ctx, "/domain/?act=moddnsrecord", params, nil)
	return record, err
}

func (p *Client) DeleteRecord(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	params := url.Values{
		"domain": {zone},
		"id":     {record.ID},
	}
	err := p.request(ctx, "/domain/?act=deldnsrecord", params, nil)
	return record, err
}

type Result struct {
	Result   int             `json:"result"`
	ClientID string          `json:"clientid"`
	Msg      string          `json:"msg"`
	ErrCode  int             `json:"errcode"`
	Data     json.RawMessage `json:"data"`
}

func (r *Result) Error() string {
	return fmt.Sprintf("请求%q失败:%d: %s(%d)", r.ClientID, r.Result, r.Msg, r.ErrCode)
}
