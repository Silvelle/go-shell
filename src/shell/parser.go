package shell

import "strings"

// Command is a parsed input line: a command name and its arguments.
type Command struct {
	Name string
	Args []string
}

// Parse splits line on whitespace into a command name and its arguments.
// An empty or blank line yields a Command with an empty Name.
func Parse(line string) (Command, error) {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return Command{
			Args: []string{},
		}, nil
	}

	return Command{
		Name: tokens[0],
		Args: tokens[1:],
	}, nil

}

// IsEmpty reports whether the command has no name, i.e. the line was blank
func (c Command) IsEmpty() bool {
	return c.Name == ""
}
