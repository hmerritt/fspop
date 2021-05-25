package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	Version           = "0.8.2"
	VersionPrerelease = "alpha"
	VersionMetadata   = ""
)
