package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	Version           = "1.3.6"
	VersionPrerelease = ""
	VersionMetadata   = ""
)
