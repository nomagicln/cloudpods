// Copyright 2019 Yunion
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

package storagemon

import (
	"yunion.io/x/onecloud/pkg/cloudmon/collectors/common"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

func init() {
	factory := SCloudAccountFactory{}
	common.RegisterFactory(&factory)
}

type SCloudAccountFactory struct {
}

func (self *SCloudAccountFactory) NewCloudReport(provider *common.SProvider, session *mcclient.ClientSession,
	args *common.ReportOptions,
	operatorType string) common.ICloudReport {
	return &SStorageReport{
		common.CloudReportBase{
			SProvider: nil,
			Session:   session,
			Args:      args,
			Operator:  string(common.STORAGE),
		},
	}
}

func (S SCloudAccountFactory) GetId() string {
	return string(common.STORAGE)
}

type SStorageReport struct {
	common.CloudReportBase
}

func (self *SStorageReport) Report() error {
	accounts, err := self.GetAllStorage(&modules.Storages)
	if err != nil {
		return err
	}
	return self.collectMetric(accounts)
}
