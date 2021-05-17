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

// Code generated by model-api-gen. DO NOT EDIT.

package cloudid

import (
	"yunion.io/x/onecloud/pkg/apis"
)

// SAMLProviderResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SAMLProviderResourceBase.
type SAMLProviderResourceBase struct {
	SAMLProviderId string `json:"saml_provider_id"`
}

// SCloudaccount is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudaccount.
type SCloudaccount struct {
	apis.SStandaloneResourceBase
	apis.SDomainizedResourceBase
	AccountId   string `json:"account_id"`
	Provider    string `json:"provider"`
	Brand       string `json:"brand"`
	IamLoginUrl string `json:"iam_login_url"`
	SAMLAuth    *bool  `json:"saml_auth,omitempty"`
	AccessUrl   string `json:"access_url"`
}

// SCloudaccountResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudaccountResourceBase.
type SCloudaccountResourceBase struct {
	// 云账号Id
	CloudaccountId string `json:"cloudaccount_id"`
}

// SCloudgroup is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroup.
type SCloudgroup struct {
	apis.SStatusInfrasResourceBase
	Provider string `json:"provider"`
}

// SCloudgroupJointsBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroupJointsBase.
type SCloudgroupJointsBase struct {
	apis.SJointResourceBase
	// 用户组Id
	CloudgroupId string `json:"cloudgroup_id"`
}

// SCloudgroupPolicy is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroupPolicy.
type SCloudgroupPolicy struct {
	SCloudgroupJointsBase
	SCloudpolicyResourceBase
}

// SCloudgroupResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroupResourceBase.
type SCloudgroupResourceBase struct {
	CloudgroupId string `json:"cloudgroup_id"`
}

// SCloudgroupUser is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroupUser.
type SCloudgroupUser struct {
	SCloudgroupJointsBase
	SClouduserResourceBase
}

// SCloudgroupcache is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudgroupcache.
type SCloudgroupcache struct {
	apis.SStatusStandaloneResourceBase
	apis.SExternalizedResourceBase
	SCloudaccountResourceBase
	// 用户组Id
	CloudgroupId string `json:"cloudgroup_id"`
}

// SCloudpolicy is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudpolicy.
type SCloudpolicy struct {
	apis.SStatusInfrasResourceBase
	apis.SExternalizedResourceBase
	// 权限类型
	// | 权限类型      |  说明                |
	// |---------------|----------------------|
	// | system        | 平台内置权限         |
	// | custom        | 用户自定义权限       |
	PolicyType string `json:"policy_type"`
	// 平台
	// | 云平台   | 说明                                        |
	// |----------|---------------------------------------------|
	// | Google   | 支持                                        |
	// | Aliyun   | 支持										|
	// | Huawei   | 支持                                        |
	// | Azure    | 支持                                        |
	// | 腾讯云   | 支持                                        |
	Provider string `json:"provider"`
	// 策略内容
	Document interface{} `json:"document"`
	// 是否锁定, 若锁定后, 此策略不允许被绑定到用户或权限组, 仅管理员可以设置是否锁定
	Locked   *bool  `json:"locked,omitempty"`
	CloudEnv string `json:"cloud_env"`
}

// SCloudpolicyResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudpolicyResourceBase.
type SCloudpolicyResourceBase struct {
	// 权限Id
	CloudpolicyId string `json:"cloudpolicy_id"`
}

// SCloudpolicycache is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudpolicycache.
type SCloudpolicycache struct {
	apis.SStatusStandaloneResourceBase
	apis.SExternalizedResourceBase
	SCloudaccountResourceBase
	SCloudproviderResourceBase
	// 权限Id
	CloudpolicyId string `json:"cloudpolicy_id"`
}

// SCloudprovider is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudprovider.
type SCloudprovider struct {
	apis.SStandaloneResourceBase
	Provider       string `json:"provider"`
	CloudaccountId string `json:"cloudaccount_id"`
}

// SCloudproviderResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudproviderResourceBase.
type SCloudproviderResourceBase struct {
	// 子订阅Id
	CloudproviderId string `json:"cloudprovider_id"`
}

// SCloudrole is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SCloudrole.
type SCloudrole struct {
	apis.SEnabledStatusInfrasResourceBase
	apis.SExternalizedResourceBase
	SCloudaccountResourceBase
	SAMLProviderResourceBase
	SCloudgroupResourceBase
	Document interface{} `json:"document"`
	OwnerId  string      `json:"owner_id"`
}

// SClouduser is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SClouduser.
type SClouduser struct {
	apis.SStatusDomainLevelUserResourceBase
	apis.SExternalizedResourceBase
	SCloudaccountResourceBase
	Secret string `json:"secret"`
	// 是否可以控制台登录
	IsConsoleLogin *bool `json:"is_console_login,omitempty"`
	// 手机号码
	MobilePhone string `json:"mobile_phone"`
	// 邮箱地址
	Email string `json:"email"`
}

// SClouduserJointsBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SClouduserJointsBase.
type SClouduserJointsBase struct {
	apis.SJointResourceBase
	SClouduserResourceBase
}

// SClouduserPolicy is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SClouduserPolicy.
type SClouduserPolicy struct {
	SClouduserJointsBase
	SCloudproviderResourceBase
	SCloudpolicyResourceBase
}

// SClouduserResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SClouduserResourceBase.
type SClouduserResourceBase struct {
	ClouduserId string `json:"clouduser_id"`
}

// SSAMLProvider is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SSAMLProvider.
type SSAMLProvider struct {
	apis.SStatusInfrasResourceBase
	apis.SExternalizedResourceBase
	SCloudaccountResourceBase
	EntityId         string `json:"entity_id"`
	MetadataDocument string `json:"metadata_document"`
	AuthUrl          string `json:"auth_url"`
}

// SSamluser is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudid/models.SSamluser.
type SSamluser struct {
	apis.SStatusDomainLevelUserResourceBase
	apis.SExternalizedResourceBase
	SCloudgroupResourceBase
	SCloudaccountResourceBase
	// 邮箱地址
	Email string `json:"email"`
}
