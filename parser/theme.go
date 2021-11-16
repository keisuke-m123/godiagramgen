package parser

func themeAdjustment(themeName string) string {
	switch themeName {
	case "reddress-darkblue",
		"reddress-darkgreen",
		"reddress-darkorange",
		"reddress-darkred",
		"reddress-lightblue",
		"reddress-lightgreen",
		"reddress-lightorange",
		"reddress-lightred":
		return `skinparam class {
  attributeIconSize 8
}`
	default:
		return ""
	}
}
