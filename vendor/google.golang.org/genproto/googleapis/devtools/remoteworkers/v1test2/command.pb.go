// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/devtools/remoteworkers/v1test2/command.proto

package remoteworkers // import "google.golang.org/genproto/googleapis/devtools/remoteworkers/v1test2"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Describes a shell-style task to execute.
type CommandTask struct {
	// The inputs to the task.
	Inputs *CommandTask_Inputs `protobuf:"bytes,1,opt,name=inputs,proto3" json:"inputs,omitempty"`
	// The expected outputs from the task.
	ExpectedOutputs *CommandTask_Outputs `protobuf:"bytes,4,opt,name=expected_outputs,json=expectedOutputs,proto3" json:"expected_outputs,omitempty"`
	// The timeouts of this task.
	Timeouts             *CommandTask_Timeouts `protobuf:"bytes,5,opt,name=timeouts,proto3" json:"timeouts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *CommandTask) Reset()         { *m = CommandTask{} }
func (m *CommandTask) String() string { return proto.CompactTextString(m) }
func (*CommandTask) ProtoMessage()    {}
func (*CommandTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{0}
}
func (m *CommandTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandTask.Unmarshal(m, b)
}
func (m *CommandTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandTask.Marshal(b, m, deterministic)
}
func (dst *CommandTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandTask.Merge(dst, src)
}
func (m *CommandTask) XXX_Size() int {
	return xxx_messageInfo_CommandTask.Size(m)
}
func (m *CommandTask) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandTask.DiscardUnknown(m)
}

var xxx_messageInfo_CommandTask proto.InternalMessageInfo

func (m *CommandTask) GetInputs() *CommandTask_Inputs {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *CommandTask) GetExpectedOutputs() *CommandTask_Outputs {
	if m != nil {
		return m.ExpectedOutputs
	}
	return nil
}

func (m *CommandTask) GetTimeouts() *CommandTask_Timeouts {
	if m != nil {
		return m.Timeouts
	}
	return nil
}

// Describes the inputs to a shell-style task.
type CommandTask_Inputs struct {
	// The command itself to run (e.g., argv)
	Arguments []string `protobuf:"bytes,1,rep,name=arguments,proto3" json:"arguments,omitempty"`
	// The input filesystem to be set up prior to the task beginning. The
	// contents should be a repeated set of FileMetadata messages though other
	// formats are allowed if better for the implementation (eg, a LUCI-style
	// .isolated file).
	//
	// This field is repeated since implementations might want to cache the
	// metadata, in which case it may be useful to break up portions of the
	// filesystem that change frequently (eg, specific input files) from those
	// that don't (eg, standard header files).
	Files []*Digest `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
	// All environment variables required by the task.
	EnvironmentVariables []*CommandTask_Inputs_EnvironmentVariable `protobuf:"bytes,3,rep,name=environment_variables,json=environmentVariables,proto3" json:"environment_variables,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *CommandTask_Inputs) Reset()         { *m = CommandTask_Inputs{} }
func (m *CommandTask_Inputs) String() string { return proto.CompactTextString(m) }
func (*CommandTask_Inputs) ProtoMessage()    {}
func (*CommandTask_Inputs) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{0, 0}
}
func (m *CommandTask_Inputs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandTask_Inputs.Unmarshal(m, b)
}
func (m *CommandTask_Inputs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandTask_Inputs.Marshal(b, m, deterministic)
}
func (dst *CommandTask_Inputs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandTask_Inputs.Merge(dst, src)
}
func (m *CommandTask_Inputs) XXX_Size() int {
	return xxx_messageInfo_CommandTask_Inputs.Size(m)
}
func (m *CommandTask_Inputs) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandTask_Inputs.DiscardUnknown(m)
}

var xxx_messageInfo_CommandTask_Inputs proto.InternalMessageInfo

func (m *CommandTask_Inputs) GetArguments() []string {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *CommandTask_Inputs) GetFiles() []*Digest {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *CommandTask_Inputs) GetEnvironmentVariables() []*CommandTask_Inputs_EnvironmentVariable {
	if m != nil {
		return m.EnvironmentVariables
	}
	return nil
}

// An environment variable required by this task.
type CommandTask_Inputs_EnvironmentVariable struct {
	// The envvar name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The envvar value.
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandTask_Inputs_EnvironmentVariable) Reset() {
	*m = CommandTask_Inputs_EnvironmentVariable{}
}
func (m *CommandTask_Inputs_EnvironmentVariable) String() string { return proto.CompactTextString(m) }
func (*CommandTask_Inputs_EnvironmentVariable) ProtoMessage()    {}
func (*CommandTask_Inputs_EnvironmentVariable) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{0, 0, 0}
}
func (m *CommandTask_Inputs_EnvironmentVariable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable.Unmarshal(m, b)
}
func (m *CommandTask_Inputs_EnvironmentVariable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable.Marshal(b, m, deterministic)
}
func (dst *CommandTask_Inputs_EnvironmentVariable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable.Merge(dst, src)
}
func (m *CommandTask_Inputs_EnvironmentVariable) XXX_Size() int {
	return xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable.Size(m)
}
func (m *CommandTask_Inputs_EnvironmentVariable) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable.DiscardUnknown(m)
}

var xxx_messageInfo_CommandTask_Inputs_EnvironmentVariable proto.InternalMessageInfo

func (m *CommandTask_Inputs_EnvironmentVariable) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CommandTask_Inputs_EnvironmentVariable) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Describes the expected outputs of the command.
type CommandTask_Outputs struct {
	// A list of expected files, relative to the execution root.
	Files []string `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	// A list of expected directories, relative to the execution root.
	Directories          []string `protobuf:"bytes,2,rep,name=directories,proto3" json:"directories,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandTask_Outputs) Reset()         { *m = CommandTask_Outputs{} }
func (m *CommandTask_Outputs) String() string { return proto.CompactTextString(m) }
func (*CommandTask_Outputs) ProtoMessage()    {}
func (*CommandTask_Outputs) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{0, 1}
}
func (m *CommandTask_Outputs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandTask_Outputs.Unmarshal(m, b)
}
func (m *CommandTask_Outputs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandTask_Outputs.Marshal(b, m, deterministic)
}
func (dst *CommandTask_Outputs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandTask_Outputs.Merge(dst, src)
}
func (m *CommandTask_Outputs) XXX_Size() int {
	return xxx_messageInfo_CommandTask_Outputs.Size(m)
}
func (m *CommandTask_Outputs) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandTask_Outputs.DiscardUnknown(m)
}

var xxx_messageInfo_CommandTask_Outputs proto.InternalMessageInfo

func (m *CommandTask_Outputs) GetFiles() []string {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *CommandTask_Outputs) GetDirectories() []string {
	if m != nil {
		return m.Directories
	}
	return nil
}

// Describes the timeouts associated with this task.
type CommandTask_Timeouts struct {
	// This specifies the maximum time that the task can run, excluding the
	// time required to download inputs or upload outputs. That is, the worker
	// will terminate the task if it runs longer than this.
	Execution *duration.Duration `protobuf:"bytes,1,opt,name=execution,proto3" json:"execution,omitempty"`
	// This specifies the maximum amount of time the task can be idle - that is,
	// go without generating some output in either stdout or stderr. If the
	// process is silent for more than the specified time, the worker will
	// terminate the task.
	Idle *duration.Duration `protobuf:"bytes,2,opt,name=idle,proto3" json:"idle,omitempty"`
	// If the execution or IO timeouts are exceeded, the worker will try to
	// gracefully terminate the task and return any existing logs. However,
	// tasks may be hard-frozen in which case this process will fail. This
	// timeout specifies how long to wait for a terminated task to shut down
	// gracefully (e.g. via SIGTERM) before we bring down the hammer (e.g.
	// SIGKILL on *nix, CTRL_BREAK_EVENT on Windows).
	Shutdown             *duration.Duration `protobuf:"bytes,3,opt,name=shutdown,proto3" json:"shutdown,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CommandTask_Timeouts) Reset()         { *m = CommandTask_Timeouts{} }
func (m *CommandTask_Timeouts) String() string { return proto.CompactTextString(m) }
func (*CommandTask_Timeouts) ProtoMessage()    {}
func (*CommandTask_Timeouts) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{0, 2}
}
func (m *CommandTask_Timeouts) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandTask_Timeouts.Unmarshal(m, b)
}
func (m *CommandTask_Timeouts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandTask_Timeouts.Marshal(b, m, deterministic)
}
func (dst *CommandTask_Timeouts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandTask_Timeouts.Merge(dst, src)
}
func (m *CommandTask_Timeouts) XXX_Size() int {
	return xxx_messageInfo_CommandTask_Timeouts.Size(m)
}
func (m *CommandTask_Timeouts) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandTask_Timeouts.DiscardUnknown(m)
}

var xxx_messageInfo_CommandTask_Timeouts proto.InternalMessageInfo

func (m *CommandTask_Timeouts) GetExecution() *duration.Duration {
	if m != nil {
		return m.Execution
	}
	return nil
}

func (m *CommandTask_Timeouts) GetIdle() *duration.Duration {
	if m != nil {
		return m.Idle
	}
	return nil
}

func (m *CommandTask_Timeouts) GetShutdown() *duration.Duration {
	if m != nil {
		return m.Shutdown
	}
	return nil
}

// Describes the actual outputs from the task.
type CommandOutputs struct {
	// exit_code is only fully reliable if the status' code is OK. If the task
	// exceeded its deadline or was cancelled, the process may still produce an
	// exit code as it is cancelled, and this will be populated, but a successful
	// (zero) is unlikely to be correct unless the status code is OK.
	ExitCode int32 `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	// The output files. The blob referenced by the digest should contain
	// one of the following (implementation-dependent):
	//    * A marshalled DirectoryMetadata of the returned filesystem
	//    * A LUCI-style .isolated file
	Outputs              *Digest  `protobuf:"bytes,2,opt,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandOutputs) Reset()         { *m = CommandOutputs{} }
func (m *CommandOutputs) String() string { return proto.CompactTextString(m) }
func (*CommandOutputs) ProtoMessage()    {}
func (*CommandOutputs) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{1}
}
func (m *CommandOutputs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandOutputs.Unmarshal(m, b)
}
func (m *CommandOutputs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandOutputs.Marshal(b, m, deterministic)
}
func (dst *CommandOutputs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandOutputs.Merge(dst, src)
}
func (m *CommandOutputs) XXX_Size() int {
	return xxx_messageInfo_CommandOutputs.Size(m)
}
func (m *CommandOutputs) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandOutputs.DiscardUnknown(m)
}

var xxx_messageInfo_CommandOutputs proto.InternalMessageInfo

func (m *CommandOutputs) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *CommandOutputs) GetOutputs() *Digest {
	if m != nil {
		return m.Outputs
	}
	return nil
}

// Can be used as part of CompleteRequest.metadata, or are part of a more
// sophisticated message.
type CommandOverhead struct {
	// The elapsed time between calling Accept and Complete. The server will also
	// have its own idea of what this should be, but this excludes the overhead of
	// the RPCs and the bot response time.
	Duration *duration.Duration `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	// The amount of time *not* spent executing the command (ie
	// uploading/downloading files).
	Overhead             *duration.Duration `protobuf:"bytes,2,opt,name=overhead,proto3" json:"overhead,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CommandOverhead) Reset()         { *m = CommandOverhead{} }
func (m *CommandOverhead) String() string { return proto.CompactTextString(m) }
func (*CommandOverhead) ProtoMessage()    {}
func (*CommandOverhead) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{2}
}
func (m *CommandOverhead) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandOverhead.Unmarshal(m, b)
}
func (m *CommandOverhead) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandOverhead.Marshal(b, m, deterministic)
}
func (dst *CommandOverhead) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandOverhead.Merge(dst, src)
}
func (m *CommandOverhead) XXX_Size() int {
	return xxx_messageInfo_CommandOverhead.Size(m)
}
func (m *CommandOverhead) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandOverhead.DiscardUnknown(m)
}

var xxx_messageInfo_CommandOverhead proto.InternalMessageInfo

func (m *CommandOverhead) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func (m *CommandOverhead) GetOverhead() *duration.Duration {
	if m != nil {
		return m.Overhead
	}
	return nil
}

// The metadata for a file. Similar to the equivalent message in the Remote
// Execution API.
type FileMetadata struct {
	// The path of this file. If this message is part of the
	// CommandResult.output_files fields, the path is relative to the execution
	// root and must correspond to an entry in CommandTask.outputs.files. If this
	// message is part of a Directory message, then the path is relative to the
	// root of that directory.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// A pointer to the contents of the file. The method by which a client
	// retrieves the contents from a CAS system is not defined here.
	Digest *Digest `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	// If the file is small enough, its contents may also or alternatively be
	// listed here.
	Contents []byte `protobuf:"bytes,3,opt,name=contents,proto3" json:"contents,omitempty"`
	// Properties of the file
	IsExecutable         bool     `protobuf:"varint,4,opt,name=is_executable,json=isExecutable,proto3" json:"is_executable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileMetadata) Reset()         { *m = FileMetadata{} }
func (m *FileMetadata) String() string { return proto.CompactTextString(m) }
func (*FileMetadata) ProtoMessage()    {}
func (*FileMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{3}
}
func (m *FileMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileMetadata.Unmarshal(m, b)
}
func (m *FileMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileMetadata.Marshal(b, m, deterministic)
}
func (dst *FileMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileMetadata.Merge(dst, src)
}
func (m *FileMetadata) XXX_Size() int {
	return xxx_messageInfo_FileMetadata.Size(m)
}
func (m *FileMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_FileMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_FileMetadata proto.InternalMessageInfo

func (m *FileMetadata) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *FileMetadata) GetDigest() *Digest {
	if m != nil {
		return m.Digest
	}
	return nil
}

func (m *FileMetadata) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

func (m *FileMetadata) GetIsExecutable() bool {
	if m != nil {
		return m.IsExecutable
	}
	return false
}

// The metadata for a directory. Similar to the equivalent message in the Remote
// Execution API.
type DirectoryMetadata struct {
	// The path of the directory, as in [FileMetadata.path][google.devtools.remoteworkers.v1test2.FileMetadata.path].
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// A pointer to the contents of the directory, in the form of a marshalled
	// Directory message.
	Digest               *Digest  `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DirectoryMetadata) Reset()         { *m = DirectoryMetadata{} }
func (m *DirectoryMetadata) String() string { return proto.CompactTextString(m) }
func (*DirectoryMetadata) ProtoMessage()    {}
func (*DirectoryMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{4}
}
func (m *DirectoryMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DirectoryMetadata.Unmarshal(m, b)
}
func (m *DirectoryMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DirectoryMetadata.Marshal(b, m, deterministic)
}
func (dst *DirectoryMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DirectoryMetadata.Merge(dst, src)
}
func (m *DirectoryMetadata) XXX_Size() int {
	return xxx_messageInfo_DirectoryMetadata.Size(m)
}
func (m *DirectoryMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_DirectoryMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_DirectoryMetadata proto.InternalMessageInfo

func (m *DirectoryMetadata) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *DirectoryMetadata) GetDigest() *Digest {
	if m != nil {
		return m.Digest
	}
	return nil
}

// A reference to the contents of a file or a directory. If the latter, the has
// refers to the hash of a marshalled Directory message. Similar to the
// equivalent message in the Remote Execution API.
type Digest struct {
	// A string-encoded hash (eg "1a2b3c", not the byte array [0x1a, 0x2b, 0x3c])
	// using an implementation-defined hash algorithm (eg SHA-256).
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	// The size of the contents. While this is not strictly required as part of an
	// identifier (after all, any given hash will have exactly one canonical
	// size), it's useful in almost all cases when one might want to send or
	// retrieve blobs of content and is included here for this reason.
	SizeBytes            int64    `protobuf:"varint,2,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Digest) Reset()         { *m = Digest{} }
func (m *Digest) String() string { return proto.CompactTextString(m) }
func (*Digest) ProtoMessage()    {}
func (*Digest) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{5}
}
func (m *Digest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Digest.Unmarshal(m, b)
}
func (m *Digest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Digest.Marshal(b, m, deterministic)
}
func (dst *Digest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Digest.Merge(dst, src)
}
func (m *Digest) XXX_Size() int {
	return xxx_messageInfo_Digest.Size(m)
}
func (m *Digest) XXX_DiscardUnknown() {
	xxx_messageInfo_Digest.DiscardUnknown(m)
}

var xxx_messageInfo_Digest proto.InternalMessageInfo

func (m *Digest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Digest) GetSizeBytes() int64 {
	if m != nil {
		return m.SizeBytes
	}
	return 0
}

// The contents of a directory. Similar to the equivalent message in the Remote
// Execution API.
type Directory struct {
	// The files in this directory
	Files []*FileMetadata `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	// Any subdirectories
	Directories          []*DirectoryMetadata `protobuf:"bytes,2,rep,name=directories,proto3" json:"directories,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Directory) Reset()         { *m = Directory{} }
func (m *Directory) String() string { return proto.CompactTextString(m) }
func (*Directory) ProtoMessage()    {}
func (*Directory) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_205cc72904871888, []int{6}
}
func (m *Directory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Directory.Unmarshal(m, b)
}
func (m *Directory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Directory.Marshal(b, m, deterministic)
}
func (dst *Directory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Directory.Merge(dst, src)
}
func (m *Directory) XXX_Size() int {
	return xxx_messageInfo_Directory.Size(m)
}
func (m *Directory) XXX_DiscardUnknown() {
	xxx_messageInfo_Directory.DiscardUnknown(m)
}

var xxx_messageInfo_Directory proto.InternalMessageInfo

func (m *Directory) GetFiles() []*FileMetadata {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *Directory) GetDirectories() []*DirectoryMetadata {
	if m != nil {
		return m.Directories
	}
	return nil
}

func init() {
	proto.RegisterType((*CommandTask)(nil), "google.devtools.remoteworkers.v1test2.CommandTask")
	proto.RegisterType((*CommandTask_Inputs)(nil), "google.devtools.remoteworkers.v1test2.CommandTask.Inputs")
	proto.RegisterType((*CommandTask_Inputs_EnvironmentVariable)(nil), "google.devtools.remoteworkers.v1test2.CommandTask.Inputs.EnvironmentVariable")
	proto.RegisterType((*CommandTask_Outputs)(nil), "google.devtools.remoteworkers.v1test2.CommandTask.Outputs")
	proto.RegisterType((*CommandTask_Timeouts)(nil), "google.devtools.remoteworkers.v1test2.CommandTask.Timeouts")
	proto.RegisterType((*CommandOutputs)(nil), "google.devtools.remoteworkers.v1test2.CommandOutputs")
	proto.RegisterType((*CommandOverhead)(nil), "google.devtools.remoteworkers.v1test2.CommandOverhead")
	proto.RegisterType((*FileMetadata)(nil), "google.devtools.remoteworkers.v1test2.FileMetadata")
	proto.RegisterType((*DirectoryMetadata)(nil), "google.devtools.remoteworkers.v1test2.DirectoryMetadata")
	proto.RegisterType((*Digest)(nil), "google.devtools.remoteworkers.v1test2.Digest")
	proto.RegisterType((*Directory)(nil), "google.devtools.remoteworkers.v1test2.Directory")
}

func init() {
	proto.RegisterFile("google/devtools/remoteworkers/v1test2/command.proto", fileDescriptor_command_205cc72904871888)
}

var fileDescriptor_command_205cc72904871888 = []byte{
	// 711 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x95, 0x4d, 0x4f, 0xdb, 0x4a,
	0x14, 0x86, 0xe5, 0x84, 0x84, 0xe4, 0x84, 0x7b, 0xb9, 0x77, 0x0a, 0x52, 0x9a, 0x7e, 0x08, 0xa5,
	0x42, 0xa2, 0x0b, 0x1c, 0x11, 0x54, 0xf5, 0x83, 0x45, 0x55, 0x08, 0x45, 0x2c, 0x50, 0xdb, 0x51,
	0x04, 0x12, 0x9b, 0x68, 0x12, 0x1f, 0x9c, 0x11, 0x89, 0x27, 0xb2, 0xc7, 0x06, 0xba, 0xa9, 0xd4,
	0x9f, 0xd2, 0x5d, 0xd9, 0xb5, 0x7f, 0xa1, 0xea, 0xff, 0xaa, 0xe6, 0xc3, 0x26, 0x29, 0x55, 0x93,
	0x66, 0xd1, 0xdd, 0xf8, 0x65, 0xde, 0x67, 0xce, 0x9c, 0x79, 0x0f, 0x81, 0x6d, 0x5f, 0x08, 0x7f,
	0x80, 0x0d, 0x0f, 0x13, 0x29, 0xc4, 0x20, 0x6a, 0x84, 0x38, 0x14, 0x12, 0x2f, 0x44, 0x78, 0x8e,
	0x61, 0xd4, 0x48, 0xb6, 0x24, 0x46, 0xb2, 0xd9, 0xe8, 0x89, 0xe1, 0x90, 0x05, 0x9e, 0x3b, 0x0a,
	0x85, 0x14, 0x64, 0xdd, 0x98, 0xdc, 0xd4, 0xe4, 0x4e, 0x98, 0x5c, 0x6b, 0xaa, 0x3d, 0xb4, 0x6c,
	0x6d, 0xea, 0xc6, 0x67, 0x0d, 0x2f, 0x0e, 0x99, 0xe4, 0x22, 0x30, 0x98, 0xfa, 0xb7, 0x22, 0x54,
	0xf6, 0x0c, 0xb8, 0xcd, 0xa2, 0x73, 0xf2, 0x0e, 0x8a, 0x3c, 0x18, 0xc5, 0x32, 0xaa, 0x3a, 0x6b,
	0xce, 0x46, 0xa5, 0xf9, 0xdc, 0x9d, 0xe9, 0x1c, 0x77, 0x8c, 0xe1, 0x1e, 0x6a, 0x00, 0xb5, 0x20,
	0x82, 0xf0, 0x1f, 0x5e, 0x8e, 0xb0, 0x27, 0xd1, 0xeb, 0x88, 0x58, 0x6a, 0xf8, 0x82, 0x86, 0xbf,
	0x98, 0x03, 0xfe, 0xc6, 0x10, 0xe8, 0x72, 0xca, 0xb4, 0x02, 0x39, 0x81, 0x92, 0xe4, 0x43, 0x14,
	0x0a, 0x5f, 0xd0, 0xf8, 0x9d, 0x39, 0xf0, 0x6d, 0x8b, 0xa0, 0x19, 0xac, 0xf6, 0x25, 0x07, 0x45,
	0x73, 0x25, 0x72, 0x1f, 0xca, 0x2c, 0xf4, 0xe3, 0x21, 0x06, 0xba, 0x41, 0xf9, 0x8d, 0x32, 0xbd,
	0x11, 0xc8, 0x1e, 0x14, 0xce, 0xf8, 0x00, 0xa3, 0x6a, 0x6e, 0x2d, 0xbf, 0x51, 0x69, 0x6e, 0xce,
	0x78, 0x7c, 0x8b, 0xfb, 0x18, 0x49, 0x6a, 0xbc, 0xe4, 0xa3, 0x03, 0xab, 0x18, 0x24, 0x3c, 0x14,
	0x81, 0xa2, 0x76, 0x12, 0x16, 0x72, 0xd6, 0x55, 0xd4, 0xbc, 0xa6, 0x1e, 0xcd, 0xfd, 0x20, 0xee,
	0xfe, 0x0d, 0xf6, 0xd8, 0x52, 0xe9, 0x0a, 0xde, 0x16, 0xa3, 0xda, 0x4b, 0xb8, 0xf3, 0x8b, 0xcd,
	0x84, 0xc0, 0x42, 0xc0, 0x86, 0xa8, 0xa3, 0x51, 0xa6, 0x7a, 0x4d, 0x56, 0xa0, 0x90, 0xb0, 0x41,
	0x8c, 0xd5, 0x9c, 0x16, 0xcd, 0x47, 0xed, 0x15, 0x2c, 0xa6, 0xef, 0xb2, 0x92, 0x76, 0xc5, 0xf4,
	0xcb, 0x5e, 0x73, 0x0d, 0x2a, 0x1e, 0x0f, 0xb1, 0x27, 0x45, 0xc8, 0x6d, 0xc7, 0xca, 0x74, 0x5c,
	0xaa, 0x7d, 0x76, 0xa0, 0x94, 0xbe, 0x06, 0x79, 0x0a, 0x65, 0xbc, 0xc4, 0x5e, 0xac, 0x92, 0x6b,
	0x93, 0x79, 0x37, 0x6d, 0x44, 0x1a, 0x6d, 0xb7, 0x65, 0xa3, 0x4d, 0x6f, 0xf6, 0x92, 0x4d, 0x58,
	0xe0, 0xde, 0xc0, 0x54, 0xf7, 0x5b, 0x8f, 0xde, 0x46, 0x9e, 0x40, 0x29, 0xea, 0xc7, 0xd2, 0x13,
	0x17, 0x41, 0x35, 0x3f, 0xcd, 0x92, 0x6d, 0xad, 0x27, 0xf0, 0xaf, 0xed, 0x77, 0x7a, 0xeb, 0x7b,
	0xaa, 0x60, 0x2e, 0x3b, 0x3d, 0xe1, 0x99, 0x7e, 0x15, 0x68, 0x49, 0x09, 0x7b, 0xc2, 0x43, 0x72,
	0x00, 0x8b, 0xe9, 0x20, 0x98, 0xba, 0xfe, 0x30, 0x2a, 0xa9, 0xbb, 0xfe, 0x01, 0x96, 0xd3, 0x73,
	0x13, 0x0c, 0xfb, 0xc8, 0x3c, 0x75, 0x83, 0x74, 0xc4, 0xa7, 0x37, 0x2a, 0xdb, 0xaa, 0x6c, 0xc2,
	0x22, 0xa6, 0xf7, 0x2a, 0xdb, 0x5a, 0xbf, 0x76, 0x60, 0xe9, 0x35, 0x1f, 0xe0, 0x11, 0x4a, 0xe6,
	0x31, 0xc9, 0x54, 0x44, 0x46, 0x4c, 0xf6, 0xd3, 0x88, 0xa8, 0x35, 0xd9, 0x87, 0xa2, 0xa7, 0x0b,
	0x9f, 0xef, 0xb6, 0xd6, 0x4c, 0x6a, 0x50, 0xea, 0x89, 0x40, 0xea, 0xd9, 0x53, 0x6f, 0xb3, 0x44,
	0xb3, 0x6f, 0xf2, 0x08, 0xfe, 0xe1, 0x51, 0xc7, 0x3c, 0xbb, 0x8a, 0xaa, 0xfe, 0x07, 0x53, 0xa2,
	0x4b, 0x3c, 0xda, 0xcf, 0xb4, 0x7a, 0x00, 0xff, 0xb7, 0x6c, 0xc0, 0xae, 0xfe, 0x42, 0xc1, 0xf5,
	0x1d, 0x28, 0x1a, 0x45, 0x1d, 0xd2, 0x67, 0x51, 0x76, 0x88, 0x5a, 0x93, 0x07, 0x00, 0x11, 0x7f,
	0x8f, 0x9d, 0xee, 0x95, 0x44, 0x93, 0x83, 0x3c, 0x2d, 0x2b, 0x65, 0x57, 0x09, 0xf5, 0xaf, 0x0e,
	0x94, 0xb3, 0x6a, 0xc9, 0xe1, 0xf8, 0x10, 0x55, 0x9a, 0xdb, 0x33, 0x16, 0x34, 0xfe, 0x34, 0xe9,
	0xe4, 0x9d, 0xde, 0x9e, 0xbc, 0x4a, 0xf3, 0xd9, 0xcc, 0x37, 0xfc, 0xa9, 0x7f, 0x13, 0x33, 0xbb,
	0xfb, 0xdd, 0x81, 0xc7, 0x3d, 0x31, 0x9c, 0x0d, 0xb6, 0xbb, 0x4a, 0xb5, 0x7c, 0x62, 0x64, 0x1b,
	0xe4, 0xe8, 0xad, 0x73, 0x4a, 0xad, 0xdf, 0x17, 0x03, 0x16, 0xf8, 0xae, 0x08, 0xfd, 0x86, 0x8f,
	0x81, 0xce, 0x61, 0xc3, 0xfc, 0x89, 0x8d, 0x78, 0x34, 0xe5, 0xf7, 0x72, 0x67, 0x42, 0xfd, 0x94,
	0xcb, 0xd1, 0x93, 0xeb, 0xdc, 0xfa, 0x81, 0x21, 0xb7, 0x30, 0x69, 0xeb, 0xca, 0x26, 0x4a, 0x70,
	0x8f, 0xb7, 0xda, 0xca, 0xda, 0x2d, 0xea, 0xb3, 0xb6, 0x7f, 0x04, 0x00, 0x00, 0xff, 0xff, 0x3d,
	0x34, 0x39, 0xf2, 0x9a, 0x07, 0x00, 0x00,
}
