package protocol

type MemberProvider interface {
	GetMembers() []ID
}