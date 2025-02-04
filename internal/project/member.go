package project

type TeamMemberID string
type TeamMember struct {
	id TeamMemberID
}

func NewTeamMember(id TeamMemberID) *TeamMember {
	return &TeamMember{id: id}
}
func (u *TeamMember) ID() TeamMemberID {
	return u.id
}
