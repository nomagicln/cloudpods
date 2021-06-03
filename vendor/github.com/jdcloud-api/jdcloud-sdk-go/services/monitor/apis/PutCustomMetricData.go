// Copyright 2018 JDCLOUD.COM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// NOTE: This class is auto generated by the jdcloud code generator program.

package apis

import (
    "github.com/jdcloud-api/jdcloud-sdk-go/core"
    monitor "github.com/jdcloud-api/jdcloud-sdk-go/services/monitor/models"
)

type PutCustomMetricDataRequest struct {

    core.JDCloudRequest

    /* 数据参数 (Optional) */
    MetricDataList []monitor.MetricDataCm `json:"metricDataList"`
}

/*
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewPutCustomMetricDataRequest(
) *PutCustomMetricDataRequest {

	return &PutCustomMetricDataRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/customMetrics",
			Method:  "POST",
			Header:  nil,
			Version: "v2",
		},
	}
}

/*
 * param metricDataList: 数据参数 (Optional)
 */
func NewPutCustomMetricDataRequestWithAllParams(
    metricDataList []monitor.MetricDataCm,
) *PutCustomMetricDataRequest {

    return &PutCustomMetricDataRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/customMetrics",
            Method:  "POST",
            Header:  nil,
            Version: "v2",
        },
        MetricDataList: metricDataList,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewPutCustomMetricDataRequestWithoutParam() *PutCustomMetricDataRequest {

    return &PutCustomMetricDataRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/customMetrics",
            Method:  "POST",
            Header:  nil,
            Version: "v2",
        },
    }
}

/* param metricDataList: 数据参数(Optional) */
func (r *PutCustomMetricDataRequest) SetMetricDataList(metricDataList []monitor.MetricDataCm) {
    r.MetricDataList = metricDataList
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r PutCustomMetricDataRequest) GetRegionId() string {
    return ""
}

type PutCustomMetricDataResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result PutCustomMetricDataResult `json:"result"`
}

type PutCustomMetricDataResult struct {
    Success bool `json:"success"`
    ErrMetricDataList []monitor.MetricDataList `json:"errMetricDataList"`
}