/*
Usage help for snapctl start subcommand:

	snapctl [OPTIONS] start [start-OPTIONS] <service>...

	The start command starts the given services of the snap. If executed from the
	"configure" hook, the services will be started after the hook finishes.

	Help Options:
	-h, --help           Show this help message

	[start command options]
			--enable     Enable the specified services (see man systemctl for
						details)
*/

package snapctl
