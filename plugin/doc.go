// Package plugin defines the plugin API for writing thriftrw plugins.
//
// Plugins take the form of an executable which runs the Main function.
//
// 	func main() {
// 		plugin.Main(plugin.Plugin{
// 			// ...
// 		})
// 	}
package plugin
