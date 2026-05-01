package wizard

type Option struct {
	ID   string
	Desc string
}

var frameworkOptions = []Option{
	{"gin", "Fast HTTP framework with middleware support"},
	{"echo", "High performance, minimalist web framework"},
	{"fiber", "Express-inspired framework built on Fasthttp"},
}

var dbOptions = []Option{
	{"postgres", "PostgreSQL relational database"},
	{"redis", "In-memory key-value store"},
	{"mysql", "MySQL relational database"},
	{"mongo", "MongoDB document database"},
	{"cassandra", "Apache Cassandra distributed database"},
	{"sqlite", "Embedded SQL database"},
}

var sdkOptions = []Option{
	{"openai", "OpenAI API client"},
	{"stripe", "Stripe payment processing"},
}

var patternOptions = []Option{
	{"worker", "Background job processing"},
}

var addonOptions = []Option{
	{"docker", "Dockerfile + docker-compose.yml"},
	{"githubactions", "GitHub Actions CI/CD workflow"},
}
