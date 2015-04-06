# logger - A simple logging layer on top of seelog

Implements multiple log channels to a single rotated log file.  Log
levels are considered a simple hierarchy with each channel having a
single limit, below which logs to that channel are skipped.

The `makeloggerinterface` is a program to generate the `interface.go`
file to simplify its maintenance.
