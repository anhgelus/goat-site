// Package site implements Standard.site lexicons in Go.
//
// [Publication] is the implementation of site.standard.publication.
//
// [Document] is the implementation of site.standard.document.
//
// [Subscription] is the implementation of site.standard.graph.subscription.
//
// See the [xrpc] package to learn how to use them.
//
// # Extending lexicons
//
// Standard lexicons are designed to be extended.
// You can extend them by declaring a new type embedding them:
//
//	type CustomPublication struct {
//	  site.Publication
//	  // your custom fields
//	}
//
// You should not reimplement [xrpc.Record], because the embedded struct already implements it.
package site
