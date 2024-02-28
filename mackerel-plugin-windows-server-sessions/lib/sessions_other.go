//go:build !windows

package mpwindowsserversessions

// Do the plugin
func Do() {
	panic("The mackerel-plugin-windows-server-sessions does not work on non Windows environment, of course.")
}
