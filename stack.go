package cfnfind

import "fmt"

type Stack struct {
	Name   string
	Region string
	Status string
}

func (s Stack) String() string {
	return fmt.Sprintf("%s\t%s\t%s", s.Name, s.Region, s.Status)
}
