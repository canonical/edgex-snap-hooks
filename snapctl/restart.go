/*
Usage help for snapctl restart subcommand:

	snapctl [OPTIONS] restart [restart-OPTIONS] <service>...

	The restart command restarts the given services of the snap. If executed from
	the
	"configure" hook, the services will be restarted after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[restart command options]
			--reload     Reload the given services if they support it (see man
						systemctl for details)
*/

package snapctl
