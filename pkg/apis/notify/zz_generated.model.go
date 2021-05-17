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

package notify

import (
	time "time"

	"yunion.io/x/onecloud/pkg/apis"
)

// SConfig is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SConfig.
type SConfig struct {
	apis.SStandaloneResourceBase
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// SNotification is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SNotification.
type SNotification struct {
	apis.SStatusStandaloneResourceBase
	ContactType string `json:"contact_type"`
	// swagger:ignore
	Topic    string `json:"topic"`
	Priority string `json:"priority"`
	// swagger:ignore
	Message     string    `json:"message"`
	ReceivedAt  time.Time `json:"received_at"`
	Event       string    `json:"event"`
	AdvanceDays int       `json:"advance_days"`
	SendTimes   int       `json:"send_times"`
	Tag         string    `json:"tag"`
}

// SReceiver is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SReceiver.
type SReceiver struct {
	apis.SStatusStandaloneResourceBase
	apis.SDomainizedResourceBase
	apis.SEnabledResourceBase
	Email string `json:"email"`
	// swagger:ignore
	Mobile string `json:"mobile"`
	Lang   string `json:"lang"`
	// swagger:ignore
	EnabledEmail *bool `json:"enabled_email,omitempty"`
	// swagger:ignore
	VerifiedEmail *bool `json:"verified_email,omitempty"`
	// swagger:ignore
	EnabledMobile *bool `json:"enabled_mobile,omitempty"`
	// swagger:ignore
	VerifiedMobile *bool `json:"verified_mobile,omitempty"`
}

// SReceiverNotification is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SReceiverNotification.
type SReceiverNotification struct {
	apis.SJointResourceBase
	ReceiverID     string `json:"receiver_id"`
	NotificationID string `json:"notification_id"`
	// ignore if ReceiverID is not empty or default
	Contact      string `json:"contact"`
	SendBy       string `json:"send_by"`
	Status       string `json:"status"`
	FailedReason string `json:"failed_reason"`
}

// SSubContact is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SSubContact.
type SSubContact struct {
	apis.SStandaloneResourceBase
	// id of receiver user
	ReceiverID        string `json:"receiver_id"`
	Type              string `json:"type"`
	Contact           string `json:"contact"`
	ParentContactType string `json:"parent_contact_type"`
	Enabled           *bool  `json:"enabled,omitempty"`
	Verified          *bool  `json:"verified,omitempty"`
	VerifiedNote      string `json:"verified_note"`
}

// SSubscription is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SSubscription.
type SSubscription struct {
	apis.SStandaloneResourceBase
	apis.SEnabledResourceBase
	Type        string `json:"type"`
	Resources   uint64 `json:"resources"`
	Actions     uint32 `json:"actions"`
	AdvanceDays int    `json:"advance_days"`
}

// SSubscriptionReceiver is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SSubscriptionReceiver.
type SSubscriptionReceiver struct {
	apis.SStandaloneResourceBase
	SubscriptionID string `json:"subscription_id"`
	// role id or receiver id or other and the value type is determined by the ReceiverType
	Receiver     string `json:"receiver"`
	ReceiverType string `json:"receiver_type"`
	RoleScope    string `json:"role_scope"`
}

// SSubscriptionReceiverDis is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SSubscriptionReceiverDis.
type SSubscriptionReceiverDis struct {
	SSubscriptionReceiver
	ReceiverName string `json:"receiver_name"`
	RoleName     string `json:"role_name"`
}

// STemplate is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.STemplate.
type STemplate struct {
	apis.SStandaloneResourceBase
	ContactType string `json:"contact_type"`
	Topic       string `json:"topic"`
	// title | content | remote
	TemplateType string `json:"template_type"`
	Content      string `json:"content"`
	Lang         string `json:"lang"`
	Example      string `json:"example"`
}

// SVerification is an autogenerated struct via yunion.io/x/onecloud/pkg/notify/models.SVerification.
type SVerification struct {
	apis.SStandaloneResourceBase
	ReceiverId  string `json:"receiver_id"`
	ContactType string `json:"contact_type"`
	Token       string `json:"token"`
}
