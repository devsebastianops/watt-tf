package plugin

type Event string

// Before transformation has taken place
const EventBeforeTransform Event = "beforeTransform"

// After transformation has taken place
const EventAfterTransform Event = "afterTransform"
