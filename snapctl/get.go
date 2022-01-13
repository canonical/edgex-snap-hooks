package snapctl

import (
	"errors"
	"strings"
)

/*
$ snapctl get --help
Usage:
  snapctl [OPTIONS] get [get-OPTIONS] [:<plug|slot>] [<keys>...]

The get command prints configuration options for the current snap.

$ snapctl get username
frank

If multiple option names are provided, a document is returned:

$ snapctl get username password
{
"username": "frank",
"password": "..."
}

Nested values may be retrieved via a dotted path:

$ snapctl get author.name
frank

Values of interface connection settings may be printed with:

$ snapctl get :myplug usb-vendor
$ snapctl get :myslot path

This will return the named setting from the local interface endpoint, whether a
plug
or a slot. Returning the setting from the connected snap's endpoint is also
possible
by explicitly requesting that via the --plug and --slot command line options:

$ snapctl get :myplug --slot usb-vendor

This requests the "usb-vendor" setting from the slot that is connected to
"myplug".


Help Options:
  -h, --help              Show this help message

[get command options]
          --slot          return attribute values from the slot side of the
                          connection
          --plug          return attribute values from the plug side of the
                          connection
      -d                  always return document, even with single key
      -t                  strict typing with nulls and quoted strings

[get command arguments]
  <keys>:                 option keys
*/

type get struct {
	options    []string
	plug, slot string
	keys       []string
}


// get takes reads config options
// It returns an object for setting the CLI arguments before running the command
func Get() get {
	return get{}
}

// Keys takes one or more keys
func (cmd get) Keys(key ...string) get {
	cmd.keys = key
	return cmd
}

// Plug takes a plug name
func (cmd get) Plug(name string) get {
	cmd.plug = name
	return cmd
}

// Slot takes a slot name
func (cmd get) Slot(name string) get {
	cmd.slot = name
	return cmd
}

// Document sets the following command option:
// -d	always return document, even with single key
func (cmd get) Document() get {
	cmd.options = append(cmd.options, "-d")
	return cmd
}

// List sets the following command option:
// -l	strict typing with nulls and quoted strings
func (cmd get) List() get {
	cmd.options = append(cmd.options, "-l")
	return cmd
}

// Run executes the get command
func (cmd get) Run() (string, error) {
	if err := cmd.validate(); err != nil {
		return "", err
	}

	// construct the command args
	// [get-OPTIONS] [:<plug|slot>] [<keys>...]
	var args []string
	// options
	args = append(args, cmd.options...)
	// plug|slot, pre-validated
	if cmd.plug != "" {
		args = append(args, ":"+cmd.plug)
	}
	if cmd.slot != "" {
		args = append(args, ":"+cmd.slot)
	}
	// keys
	args = append(args, cmd.keys...)

	return run("get", args...)
}

func (cmd get) validate() error {
	if cmd.plug != "" && cmd.slot != "" {
		return errors.New("only one of plug or slot can be set")
	}
	if strings.HasPrefix(cmd.plug, ":") {
		return errors.New("plug name must not contain colon as prefix")
	}
	if strings.HasPrefix(cmd.slot, ":") {
		return errors.New("slot name must not contain colon as prefix")
	}
	return nil
}