# protoc-gen-yaml 

protoc-gen-yaml is a Protobuf plugin to generate a yaml file that is a subset of the Profobuf file's content.

# Build 

```
go build .
```

# Run the plugin 

```
protoc --plugin protoc-gen-yaml --yaml_out=./ --yaml_opt=paths=source_relative example/example.proto
```
