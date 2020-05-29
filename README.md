# Indico Proxy common

This project holds common resources and interfaces for Indico Proxy, and its connectors. These are meant to be resources for independent connectors for Proxy.

## What is Indico Proxy

Proxy can receive files from clients as they upload, and pass these files to one or more custom backends. It can also receive files from various other sources, like watching a folder for events, as a cli-tool.

In Indico, Proxy serves as an very flexible translator from various sources to other backends.

## What is a connector

 These backends may have vastly different mechanics, from S3 to webdav and really anything. To support the backend, a small connector must be created.

The connector lives in the same binary as Indico Proxy, and receives files and metadata and other operations as high-level function calls. 

