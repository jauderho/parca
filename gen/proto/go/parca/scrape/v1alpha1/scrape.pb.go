// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: parca/scrape/v1alpha1/scrape.proto

package scrapev1alpha1

import (
	v1alpha1 "github.com/parca-dev/parca/gen/proto/go/parca/profilestore/v1alpha1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// State represents the current state of a target
type TargetsRequest_State int32

const (
	// STATE_ANY_UNSPECIFIED unspecified
	TargetsRequest_STATE_ANY_UNSPECIFIED TargetsRequest_State = 0
	// STATE_ACTIVE target active state
	TargetsRequest_STATE_ACTIVE TargetsRequest_State = 1
	// STATE_DROPPED target dropped state
	TargetsRequest_STATE_DROPPED TargetsRequest_State = 2
)

// Enum value maps for TargetsRequest_State.
var (
	TargetsRequest_State_name = map[int32]string{
		0: "STATE_ANY_UNSPECIFIED",
		1: "STATE_ACTIVE",
		2: "STATE_DROPPED",
	}
	TargetsRequest_State_value = map[string]int32{
		"STATE_ANY_UNSPECIFIED": 0,
		"STATE_ACTIVE":          1,
		"STATE_DROPPED":         2,
	}
)

func (x TargetsRequest_State) Enum() *TargetsRequest_State {
	p := new(TargetsRequest_State)
	*p = x
	return p
}

func (x TargetsRequest_State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TargetsRequest_State) Descriptor() protoreflect.EnumDescriptor {
	return file_parca_scrape_v1alpha1_scrape_proto_enumTypes[0].Descriptor()
}

func (TargetsRequest_State) Type() protoreflect.EnumType {
	return &file_parca_scrape_v1alpha1_scrape_proto_enumTypes[0]
}

func (x TargetsRequest_State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TargetsRequest_State.Descriptor instead.
func (TargetsRequest_State) EnumDescriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{0, 0}
}

// Health are the possible health values of a target
type Target_Health int32

const (
	// HEALTH_UNKNOWN_UNSPECIFIED unspecified
	Target_HEALTH_UNKNOWN_UNSPECIFIED Target_Health = 0
	// HEALTH_GOOD healthy target
	Target_HEALTH_GOOD Target_Health = 1
	// HEALTH_BAD unhealthy target
	Target_HEALTH_BAD Target_Health = 2
)

// Enum value maps for Target_Health.
var (
	Target_Health_name = map[int32]string{
		0: "HEALTH_UNKNOWN_UNSPECIFIED",
		1: "HEALTH_GOOD",
		2: "HEALTH_BAD",
	}
	Target_Health_value = map[string]int32{
		"HEALTH_UNKNOWN_UNSPECIFIED": 0,
		"HEALTH_GOOD":                1,
		"HEALTH_BAD":                 2,
	}
)

func (x Target_Health) Enum() *Target_Health {
	p := new(Target_Health)
	*p = x
	return p
}

func (x Target_Health) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Target_Health) Descriptor() protoreflect.EnumDescriptor {
	return file_parca_scrape_v1alpha1_scrape_proto_enumTypes[1].Descriptor()
}

func (Target_Health) Type() protoreflect.EnumType {
	return &file_parca_scrape_v1alpha1_scrape_proto_enumTypes[1]
}

func (x Target_Health) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Target_Health.Descriptor instead.
func (Target_Health) EnumDescriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{3, 0}
}

// TargetsRequest contains the parameters for the set of targets to return
type TargetsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// state is the state of targets to returns
	State TargetsRequest_State `protobuf:"varint,1,opt,name=state,proto3,enum=parca.scrape.v1alpha1.TargetsRequest_State" json:"state,omitempty"`
}

func (x *TargetsRequest) Reset() {
	*x = TargetsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TargetsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TargetsRequest) ProtoMessage() {}

func (x *TargetsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TargetsRequest.ProtoReflect.Descriptor instead.
func (*TargetsRequest) Descriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{0}
}

func (x *TargetsRequest) GetState() TargetsRequest_State {
	if x != nil {
		return x.State
	}
	return TargetsRequest_STATE_ANY_UNSPECIFIED
}

// TargetsResponse is the set of targets for the given requested state
type TargetsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// targets is the mapping of targets
	Targets map[string]*Targets `protobuf:"bytes,1,rep,name=targets,proto3" json:"targets,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TargetsResponse) Reset() {
	*x = TargetsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TargetsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TargetsResponse) ProtoMessage() {}

func (x *TargetsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TargetsResponse.ProtoReflect.Descriptor instead.
func (*TargetsResponse) Descriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{1}
}

func (x *TargetsResponse) GetTargets() map[string]*Targets {
	if x != nil {
		return x.Targets
	}
	return nil
}

// Targets is a list of targets
type Targets struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// targets is a list of targets
	Targets []*Target `protobuf:"bytes,1,rep,name=targets,proto3" json:"targets,omitempty"`
}

func (x *Targets) Reset() {
	*x = Targets{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Targets) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Targets) ProtoMessage() {}

func (x *Targets) ProtoReflect() protoreflect.Message {
	mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Targets.ProtoReflect.Descriptor instead.
func (*Targets) Descriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{2}
}

func (x *Targets) GetTargets() []*Target {
	if x != nil {
		return x.Targets
	}
	return nil
}

// Target is the scrape target representation
type Target struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// discovered_labels are the set of labels for the target that have been discovered
	DiscoveredLabels *v1alpha1.LabelSet `protobuf:"bytes,1,opt,name=discovered_labels,json=discoveredLabels,proto3" json:"discovered_labels,omitempty"`
	// labels are the set of labels given for the target
	Labels *v1alpha1.LabelSet `protobuf:"bytes,2,opt,name=labels,proto3" json:"labels,omitempty"`
	// last_error is the error message most recently received from a scrape attempt
	LastError string `protobuf:"bytes,3,opt,name=last_error,json=lastError,proto3" json:"last_error,omitempty"`
	// last_scrape is the time stamp the last scrape request was performed
	LastScrape *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=last_scrape,json=lastScrape,proto3" json:"last_scrape,omitempty"`
	// last_scrape_duration is the duration of the last scrape request
	LastScrapeDuration *durationpb.Duration `protobuf:"bytes,5,opt,name=last_scrape_duration,json=lastScrapeDuration,proto3" json:"last_scrape_duration,omitempty"`
	// url is the url of the target
	Url string `protobuf:"bytes,6,opt,name=url,proto3" json:"url,omitempty"`
	// health indicates the current health of the target
	Health Target_Health `protobuf:"varint,7,opt,name=health,proto3,enum=parca.scrape.v1alpha1.Target_Health" json:"health,omitempty"`
}

func (x *Target) Reset() {
	*x = Target{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Target) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Target) ProtoMessage() {}

func (x *Target) ProtoReflect() protoreflect.Message {
	mi := &file_parca_scrape_v1alpha1_scrape_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Target.ProtoReflect.Descriptor instead.
func (*Target) Descriptor() ([]byte, []int) {
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP(), []int{3}
}

func (x *Target) GetDiscoveredLabels() *v1alpha1.LabelSet {
	if x != nil {
		return x.DiscoveredLabels
	}
	return nil
}

func (x *Target) GetLabels() *v1alpha1.LabelSet {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Target) GetLastError() string {
	if x != nil {
		return x.LastError
	}
	return ""
}

func (x *Target) GetLastScrape() *timestamppb.Timestamp {
	if x != nil {
		return x.LastScrape
	}
	return nil
}

func (x *Target) GetLastScrapeDuration() *durationpb.Duration {
	if x != nil {
		return x.LastScrapeDuration
	}
	return nil
}

func (x *Target) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Target) GetHealth() Target_Health {
	if x != nil {
		return x.Health
	}
	return Target_HEALTH_UNKNOWN_UNSPECIFIED
}

var File_parca_scrape_v1alpha1_scrape_proto protoreflect.FileDescriptor

var file_parca_scrape_v1alpha1_scrape_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61,
	0x70, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x61, 0x72, 0x63,
	0x61, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x0e, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x70,
	0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x22, 0x47, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x41, 0x4e, 0x59, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f,
	0x44, 0x52, 0x4f, 0x50, 0x50, 0x45, 0x44, 0x10, 0x02, 0x22, 0xbc, 0x01, 0x0a, 0x0f, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a,
	0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x33,
	0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x1a, 0x5a, 0x0a, 0x0c,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x34,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x42, 0x0a, 0x07, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x73, 0x12, 0x37, 0x0a, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72,
	0x61, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x52, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x22, 0xdf, 0x03, 0x0a,
	0x06, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x52, 0x0a, 0x11, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x74, 0x52, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x65, 0x64, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x3d, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x61,
	0x72, 0x63, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53,
	0x65, 0x74, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x3b, 0x0a, 0x0b, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74,
	0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x12, 0x4b, 0x0a, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73,
	0x63, 0x72, 0x61, 0x70, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x12, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x3c, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63,
	0x72, 0x61, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x06, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x22, 0x49, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x1e, 0x0a,
	0x1a, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x5f, 0x47, 0x4f, 0x4f, 0x44, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x5f, 0x42, 0x41, 0x44, 0x10, 0x02, 0x32, 0x7b,
	0x0a, 0x0d, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x6a, 0x0a, 0x07, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x12, 0x25, 0x2e, 0x70, 0x61, 0x72,
	0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x26, 0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x10, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0a, 0x12, 0x08, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x42, 0xec, 0x01, 0x0a, 0x19,
	0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0b, 0x53, 0x63, 0x72, 0x61, 0x70,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x70,
	0x61, 0x72, 0x63, 0x61, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x61, 0x72, 0x63, 0x61, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x53, 0x58, 0xaa, 0x02, 0x15, 0x50,
	0x61, 0x72, 0x63, 0x61, 0x2e, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x2e, 0x56, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x15, 0x50, 0x61, 0x72, 0x63, 0x61, 0x5c, 0x53, 0x63, 0x72,
	0x61, 0x70, 0x65, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x21, 0x50,
	0x61, 0x72, 0x63, 0x61, 0x5c, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x5c, 0x56, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x17, 0x50, 0x61, 0x72, 0x63, 0x61, 0x3a, 0x3a, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65,
	0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_parca_scrape_v1alpha1_scrape_proto_rawDescOnce sync.Once
	file_parca_scrape_v1alpha1_scrape_proto_rawDescData = file_parca_scrape_v1alpha1_scrape_proto_rawDesc
)

func file_parca_scrape_v1alpha1_scrape_proto_rawDescGZIP() []byte {
	file_parca_scrape_v1alpha1_scrape_proto_rawDescOnce.Do(func() {
		file_parca_scrape_v1alpha1_scrape_proto_rawDescData = protoimpl.X.CompressGZIP(file_parca_scrape_v1alpha1_scrape_proto_rawDescData)
	})
	return file_parca_scrape_v1alpha1_scrape_proto_rawDescData
}

var file_parca_scrape_v1alpha1_scrape_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_parca_scrape_v1alpha1_scrape_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_parca_scrape_v1alpha1_scrape_proto_goTypes = []interface{}{
	(TargetsRequest_State)(0),     // 0: parca.scrape.v1alpha1.TargetsRequest.State
	(Target_Health)(0),            // 1: parca.scrape.v1alpha1.Target.Health
	(*TargetsRequest)(nil),        // 2: parca.scrape.v1alpha1.TargetsRequest
	(*TargetsResponse)(nil),       // 3: parca.scrape.v1alpha1.TargetsResponse
	(*Targets)(nil),               // 4: parca.scrape.v1alpha1.Targets
	(*Target)(nil),                // 5: parca.scrape.v1alpha1.Target
	nil,                           // 6: parca.scrape.v1alpha1.TargetsResponse.TargetsEntry
	(*v1alpha1.LabelSet)(nil),     // 7: parca.profilestore.v1alpha1.LabelSet
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 9: google.protobuf.Duration
}
var file_parca_scrape_v1alpha1_scrape_proto_depIdxs = []int32{
	0,  // 0: parca.scrape.v1alpha1.TargetsRequest.state:type_name -> parca.scrape.v1alpha1.TargetsRequest.State
	6,  // 1: parca.scrape.v1alpha1.TargetsResponse.targets:type_name -> parca.scrape.v1alpha1.TargetsResponse.TargetsEntry
	5,  // 2: parca.scrape.v1alpha1.Targets.targets:type_name -> parca.scrape.v1alpha1.Target
	7,  // 3: parca.scrape.v1alpha1.Target.discovered_labels:type_name -> parca.profilestore.v1alpha1.LabelSet
	7,  // 4: parca.scrape.v1alpha1.Target.labels:type_name -> parca.profilestore.v1alpha1.LabelSet
	8,  // 5: parca.scrape.v1alpha1.Target.last_scrape:type_name -> google.protobuf.Timestamp
	9,  // 6: parca.scrape.v1alpha1.Target.last_scrape_duration:type_name -> google.protobuf.Duration
	1,  // 7: parca.scrape.v1alpha1.Target.health:type_name -> parca.scrape.v1alpha1.Target.Health
	4,  // 8: parca.scrape.v1alpha1.TargetsResponse.TargetsEntry.value:type_name -> parca.scrape.v1alpha1.Targets
	2,  // 9: parca.scrape.v1alpha1.ScrapeService.Targets:input_type -> parca.scrape.v1alpha1.TargetsRequest
	3,  // 10: parca.scrape.v1alpha1.ScrapeService.Targets:output_type -> parca.scrape.v1alpha1.TargetsResponse
	10, // [10:11] is the sub-list for method output_type
	9,  // [9:10] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_parca_scrape_v1alpha1_scrape_proto_init() }
func file_parca_scrape_v1alpha1_scrape_proto_init() {
	if File_parca_scrape_v1alpha1_scrape_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_parca_scrape_v1alpha1_scrape_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TargetsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_parca_scrape_v1alpha1_scrape_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TargetsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_parca_scrape_v1alpha1_scrape_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Targets); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_parca_scrape_v1alpha1_scrape_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Target); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_parca_scrape_v1alpha1_scrape_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_parca_scrape_v1alpha1_scrape_proto_goTypes,
		DependencyIndexes: file_parca_scrape_v1alpha1_scrape_proto_depIdxs,
		EnumInfos:         file_parca_scrape_v1alpha1_scrape_proto_enumTypes,
		MessageInfos:      file_parca_scrape_v1alpha1_scrape_proto_msgTypes,
	}.Build()
	File_parca_scrape_v1alpha1_scrape_proto = out.File
	file_parca_scrape_v1alpha1_scrape_proto_rawDesc = nil
	file_parca_scrape_v1alpha1_scrape_proto_goTypes = nil
	file_parca_scrape_v1alpha1_scrape_proto_depIdxs = nil
}
