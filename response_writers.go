package jsonapi

// TopLevelDocumentInterfaces

type (
	DataSetter interface {
		SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}

	DataAppender interface {
		AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}

	IdentifierSetter interface {
		SetIdentifier(resourceType, id string) error
	}

	IdentifierAppender interface {
		AppendIdentifier(resourceType, id string) error
	}

	ErrorAppender interface {
		AppendError(err error)
	}

	Includer interface {
		Include(resourceType, id string, attributes interface{}, relationships Relationships, links Links, meta Meta) error
	}
)
