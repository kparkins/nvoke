package cmd

const BibleTextIndex = `{
	"mappings": {
	  "dynamic": true
	},
	"storedSource": {
	  "include": [
		"text"
	  ]
	}
}`

const BibleEmbeddingVectorIndex = `{
	"fields": [
		{
			"numDimensions": 1536,
			"path": "embedding",
			"similarity": "cosine",
			"type": "vector"
		}
	]
}`
