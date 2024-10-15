// assets in locale and templates folders are converted into Go source code with go-bindata
// the data.go file in this package is auto generated through the generate-debug and generate-prod tasks in the Makefile
package assets

// these two functions stop a 'could not import' lint error
func Asset(_ string) ([]byte, error) {
	return nil, nil
}

func AssetNames() []string {
	return []string{}
}
