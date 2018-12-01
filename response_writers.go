package jsonapi

// TopLevelDocumentInterfaces

type (
	// DataSetter represents the interface to set an individual resource object.
	DataSetter interface {
		SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}

	// DataAppender represents the interface to set a collection resource objects.
	DataAppender interface {
		AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}

	// IdentifierSetter represents the interface to set a single resource linkage.
	IdentifierSetter interface {
		SetIdentifier(resourceType, id string) error
	}

	// IdentifierAppender represents the interface to set a collection of resource
	// linkages.
	IdentifierAppender interface {
		AppendIdentifier(resourceType, id string) error
	}

	// ErrorAppender represents the interface to append an error to the `errors`
	// array member of top level document.
	ErrorAppender interface {
		AppendError(err error)
	}

	// Includer represents the interface to append resource objects to the
	// `included` array member of top level document.
	Includer interface {
		Include(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}
)
