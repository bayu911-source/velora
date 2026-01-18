
package pkg

type Agent interface {
	Name() string
	Run(memory *MemoryManager, input string) (string, error)
}
