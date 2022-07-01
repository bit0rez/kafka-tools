# kafka-tools

## Usage
### Consumer

```
  $ kafka-tools consume -t Topic -f raw | while read line; \
    do \
      echo $line | protoc -I ~/path_to_proto_libs/ --decode=package.ProtoMessage ~/path_to/message.proto; \
      echo ;\
    done
```
