package v1alpha2

// Instructions is a list of shell commands
type Instructions []Instruction

// Instruction is a single shell command
type Instruction string

// Script formats Instruction for /bin/sh -c execution
func (ins Instruction) Script() string {
	return "set -eou pipefail\n" + string(ins)
}
