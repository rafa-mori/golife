package property

// Metadata represents a map of metadata key-value pairs.
type Metadata map[string]interface{}

// EventMetadata represents metadata for a change event.
type EventMetadata struct {
	Timestamp string
	Source    string
	Details   map[string]interface{}
}
