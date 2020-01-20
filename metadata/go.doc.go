// Documentation for the normalized metadata-representation that Indico-Proxy uses
//
// The goal of this definition is to support a broad range of metadata that a client can produce, and then be able
// to map that data back to any of the supported backends.
//
// All metadata-fields should be defined in close cooperation with backend-providers, as well as clients.
//
// Definitions
//
// The following definitions are used throughout this document.
//
// **Bb**
//	Backend: Any server that can recieve files, and preferably, some metadata as well.
//	Client: Any software that can produce files and metadata and provide it to Proxy.
// 	External ID: A Unique identifier on the backend of a file.
// 	Parent ID: A unique identifier on the backend of a collection of files, like a folder, or an album
// 	Case ID: A unique identifier for a case on of a case that the file belong to.
//
// Clients and Backends do not need to know anything about each other, but the client needs to know about Proxy.
//
// Client-integrations
//
// Indico provides a DotNET-SDK. If you want to create a custom-client, we recommend starting by using a compatible
// library from https://tus.io and then extending it with the custom-layer that Proxy adheres to. Proxy uses the
// TUS-protocol for all uploads.
//
// Below, a basic flow for creating and uploading a new file is provided.
//
// Creating files
//
// A file must be created, with all its metadata prior to upload. The metadata will be validated internally, and
// with the backend. Any validation-failures will be returned to the client.
//
// Example of error returned:
//
//	{
//	  "type": "https://httpstatuses.com/422",
//	  "title": "Unprocessable Entity",
//	  "status": 422,
//	  "detail": "The following fields are missing: parentid",
//	  "instance": ""
//	}
//
// The size of each chunk can be set by the client, but we recommend at least 6 MB. Each chunk can optionally be
// verified by providing the checksum of the current chunk.
//
// 1. Regular https://tus.io -implementation: POSTing to `/` with metadata in header. For clients with minimal
// metadata, this method is preferred.
//
// 2. Extended by POST-ing metadata as JSON to `/create`.
//
// A unique upload-url should have been returned in the previous step. Uploads can now begin, with chunk-uploads
// preferred for large uploads. This upload-url should be stored on the client at least until the upload is completed.
// The upload-url can also be used to update metadata, if the backend supports it.
//
// Uploading file
//
//
// Info recieved from backend
//
// A backend will provide its External ID, Parent ID and Case ID (if supported) of the file, but the time of which it
// does so may change from one backend to another. The client should accommodate this. It may return in the data in the
// following events.
//
// 1. The client creates a new upload.
//
// 2. The client completes an upload
//
// 3. As a last resort, information can be available at a later time, and the client needs to accommodate this, if it
// wants to receive this information. The information will be available when the client completes an upload and the
// file is also sent to the backend. The client then needs to ask Proxy if the information is available.
//
// Information can be retrieved by using a HEAD-request to the upload-url, or to get the information as JSON,
// A GET-request can be sent to the upload-url prepended by `info`.
//
// Updating metadata
//
// The client may optionally, if the backend supports it, update the metadata. Proxy requires that the client
// sends in ALL the metadata.
//
// Updates are done in the same way as creating files, as described above, but the url is `update`
//
//
package metadata
