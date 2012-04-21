package state

type LeaderSort []LeaderInfo

func (l LeaderSort) Len() int {
	return len(l)
}
func (l LeaderSort) Less(i, j int) bool {
	if l[i].Worth == l[j].Worth {
		return l[i].Name < l[j].Name
	}
	return l[i].Worth > l[j].Worth
}
func (l LeaderSort) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
