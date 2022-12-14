package db

type Map map[string]interface{}

func collectionsMap() map[string]string {
	collections := map[string]string{
		"devices":                "devices",
		"credentials":            "credentials",
		"exports":                "exports",
		"exportsRetentionPolicy": "exportsRetentionPolicy",
		"users":                  "users",
		"sessions":               "sessions",
		// "neighbors":   "neighbors",
	}
	return collections
}
