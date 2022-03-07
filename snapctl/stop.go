/*
Usage help for snapctl stop subcommand:

	snapctl [OPTIONS] stop [stop-OPTIONS] <service>...

	The stop command stops the given services of the snap. If executed from the
	"configure" hook, the services will be stopped after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[stop command options]
			--disable    Disable the specified services (see man systemctl for
						details)
*/

package snapctl
