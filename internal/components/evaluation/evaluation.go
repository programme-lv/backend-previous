package evaluation

type TestDownloadURLProvider interface {
	GetTestDownloadURL(testSHA256 string) (string, error)
}
