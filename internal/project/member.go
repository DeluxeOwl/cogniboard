package project

type (
	TeamMemberID string
	TeamMember   struct {
		id   TeamMemberID
		name string
	}
)

func NewTeamMember(id TeamMemberID, name string) *TeamMember {
	return &TeamMember{id: id}
}

func (u *TeamMember) ID() TeamMemberID {
	return u.id
}

func (u *TeamMember) Name() string {
	return u.name
}
