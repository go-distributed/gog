package config

// Config describes the config of the system.
type Config struct {
	// Fanin is the nodes we allow to have
	// us in their active view.
	Fanin int
	// Fanout is the number of nodes in our
	// active view.
	Fanout int
	// AViewSize is the size of the active view.
	AViewSize int
	// PViewSize is the size of the passive view.
	PViewSize int
	// Ka is the number of nodes in active view
	// when shuffle views.
	Ka int
	// Kp is the number of nodes in passive view
	// when shuffle views.
	Kp int
}
