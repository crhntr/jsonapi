package jsonapi

// TopLevelDocumentInterfaces
type (
	DataSetter interface {
		SetData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
	}

	DataAppender interface {
		AppendData(resourceType, id string, attributes interface{}, relationships Relationships, links Linker, meta Meta) error
	}

	IdentifierSetter interface {
		SetData(resourceType, id string) error
	}

	IdentifierAppender interface {
		AppendData(resourceType, id string) error
	}

	ErrorAppender interface {
		AppendError(err error)
	}

	Includer interface {
		Include(resourceType, id string, attributes interface{}, links Linker, meta Meta) error
	}
)
