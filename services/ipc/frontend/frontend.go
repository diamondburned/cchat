// Package frontend
package frontend

import "github.com/diamondburned/cchat/services"

func init() {
	services.RegisterSource(loadPlugins)
}

func loadPlugins() []error {
	panic("Implement me")

	// d, err := ioutil.ReadDir("stuff")
	// // err check

	// for _, info := range d {
	// 	f, err := os.Open(filepath.Join("stuff", file.Name()))
	// 	// err check

	// 	services.RegisterService(f)
	// }
}
