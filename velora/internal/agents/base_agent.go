
package agents

type Agent interface {
	Name() string
	Run(input string) (string, error)
}
