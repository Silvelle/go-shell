package shell

import "strings"

type Command struct {
	Name string
	Args []string
}

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

func (c Command) IsEmpty() bool {
	return c.Name == ""
}
