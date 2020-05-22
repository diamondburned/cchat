// Package services provides a global repository of cchat services. It also
// supports additional sources.
//
// Registering services
//
// To register a service, it's best to call RegisterService() in the package's
// init(). This allows for dash imports:
//
//    _ "git.sr.ht/~user/cchat-abc"
//
// Registering sources
//
// Sources are simply functions that manage other services. An example of this
// would be the plugins package. Note that only packages that can error out on
// load should do this. A package can call RegisterService() multiple times.
//
// For examples on using RegisterSource(), check the plugins package.
package services

import (
	"sync"

	"github.com/diamondburned/cchat"
)

var services []cchat.Service

// RegisterService adds a service.
func RegisterService(service ...cchat.Service) {
	services = append(services, service...)
}

var sources []func() []error
var sourceErrs []error
var sourceOnce sync.Once

// RegisterSource adds a service source. Services are expected to call
// RegisterService() on source().
func RegisterSource(source func() []error) {
	sources = append(sources, source)
}

// Get returns all services. It will also fetch the plugins from all sources.
// Future calls will not fetch the plugins again.
func Get() ([]cchat.Service, []error) {
	sourceOnce.Do(func() {
		sourceErrs = []error{} // mark as non-nil
		for _, src := range sources {
			sourceErrs = append(sourceErrs, src()...)
		}
	})

	// why are we here, just to suffer
	return services, sourceErrs
}
