# snapctl
Go wrapper library for the [snapctl](https://snapcraft.io/docs/using-snapctl) tool.

Wrappers for following subcommands are implemented:

- [ ] `fde-setup-request`: Obtain full disk encryption setup request
- [ ] `fde-setup-result`: Set result for full disk encryption
- [x] `get`: The get command prints configuration and interface connection settings.                
- [ ] `is-connected`: Return success if the given plug or slot is connected, and failure otherwise   
- [ ] `reboot`: Control the reboot behavior of the system          
- [x] `restart`: Restart services    
- [x] `services`: Query the status of services      
- [x] `set`: Changes configuration options
- [ ] `set-health`: Report the health status of a snap
- [x] `start`: Start services 
- [x] `stop`: Stop services
- [ ] `system-mode`: Get the current system mode and associated details
- [x] `unset`: Remove configuration options

The commands and descriptions are from `snapctl --help`.
