package service

func treeGraphLimitNodeLabels() []string {
	return []string{"Tag", "Company", "Classification", "Application"}
}

func treeGraphLimitRelLabels() []string {
	return []string{"ATTACH_TO", "CLASSIFY_OF", "APPLICATION_OF"}
}
