package http

import (
	"git.meiqia.com/business_platform/tool"
	"io/ioutil"
	"net/http"
	"net/url"
)

func pipelineWithRespToJson(respi interface{}) (interface{}, error) {
	return tool.PipelineWith(respi).Do(func(respi interface{}) (interface{}, error) {
		resp := respi.(http.Response)
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}).Do(func(rsti interface{}) (interface{}, error) {
		rst := rsti.([]byte)
		var dat map[string]interface{}
		err := tool.JSON.Unmarshal(rst, dat)
		return dat, err
	}).Done()
}

func GetJSON(urlS string, params map[string][]string) (map[string]interface{}, error) {
	rst, err := tool.PipelineWith(nil).Do(func(_ interface{}) (interface{}, error) {
		return http.NewRequest("GET", urlS, nil)
	}).Do(func(reqi interface{}) (interface{}, error) {
		req := reqi.(http.Request)
		req.URL.RawQuery = url.Values(params).Encode()
		client := &http.Client{}
		return client.Do(&req)
	}).Do(pipelineWithRespToJson).Done()
	return rst.(map[string]interface{}), err
}

func PostJSON(urlS string, params map[string][]string) (map[string]interface{}, error) {
	rst, err := tool.PipelineWith(nil).Do(func(interface{}) (interface{}, error) {
		return http.PostForm(urlS, url.Values(params))
	}).Do(pipelineWithRespToJson).Done()
	return rst.(map[string]interface{}), err
}
