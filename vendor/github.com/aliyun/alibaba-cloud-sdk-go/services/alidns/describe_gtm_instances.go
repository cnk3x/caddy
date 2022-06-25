package alidns

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeGtmInstances invokes the alidns.DescribeGtmInstances API synchronously
func (client *Client) DescribeGtmInstances(request *DescribeGtmInstancesRequest) (response *DescribeGtmInstancesResponse, err error) {
	response = CreateDescribeGtmInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeGtmInstancesWithChan invokes the alidns.DescribeGtmInstances API asynchronously
func (client *Client) DescribeGtmInstancesWithChan(request *DescribeGtmInstancesRequest) (<-chan *DescribeGtmInstancesResponse, <-chan error) {
	responseChan := make(chan *DescribeGtmInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeGtmInstances(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeGtmInstancesWithCallback invokes the alidns.DescribeGtmInstances API asynchronously
func (client *Client) DescribeGtmInstancesWithCallback(request *DescribeGtmInstancesRequest, callback func(response *DescribeGtmInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeGtmInstancesResponse
		var err error
		defer close(result)
		response, err = client.DescribeGtmInstances(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeGtmInstancesRequest is the request struct for api DescribeGtmInstances
type DescribeGtmInstancesRequest struct {
	*requests.RpcRequest
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	ResourceGroupId      string           `position:"Query" name:"ResourceGroupId"`
	UserClientIp         string           `position:"Query" name:"UserClientIp"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	Lang                 string           `position:"Query" name:"Lang"`
	Keyword              string           `position:"Query" name:"Keyword"`
	NeedDetailAttributes requests.Boolean `position:"Query" name:"NeedDetailAttributes"`
}

// DescribeGtmInstancesResponse is the response struct for api DescribeGtmInstances
type DescribeGtmInstancesResponse struct {
	*responses.BaseResponse
	RequestId    string                             `json:"RequestId" xml:"RequestId"`
	PageNumber   int                                `json:"PageNumber" xml:"PageNumber"`
	PageSize     int                                `json:"PageSize" xml:"PageSize"`
	TotalItems   int                                `json:"TotalItems" xml:"TotalItems"`
	TotalPages   int                                `json:"TotalPages" xml:"TotalPages"`
	GtmInstances GtmInstancesInDescribeGtmInstances `json:"GtmInstances" xml:"GtmInstances"`
}

// CreateDescribeGtmInstancesRequest creates a request to invoke DescribeGtmInstances API
func CreateDescribeGtmInstancesRequest() (request *DescribeGtmInstancesRequest) {
	request = &DescribeGtmInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "DescribeGtmInstances", "alidns", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeGtmInstancesResponse creates a response to parse from DescribeGtmInstances response
func CreateDescribeGtmInstancesResponse() (response *DescribeGtmInstancesResponse) {
	response = &DescribeGtmInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
