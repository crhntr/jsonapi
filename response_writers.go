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

	// IdentitySetter represents the interface to set a single resource linkage.
	IdentitySetter interface {
		SetIdentity(resourceType, id string) error
	}

	// IdentityAppender represents the interface to set a collection of resource
	// linkages.
	IdentityAppender interface {
		AppendIdentity(resourceType, id string) error
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

	// DataCollectionSetter represents the interface to ensure top level document
	//  member `data` is encoded as an empty array when encoding an empty
	// collection. It is used interanally and is exported for mocking responses.
	DataCollectionSetter interface {
		SetDataCollection()
	}
)
