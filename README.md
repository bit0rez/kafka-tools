# kafka-tools

## Usage
### Consumer

```
  $ kafka-tools consume -t Topic -f raw | while read line; \
    do \
      echo $line | protoc -I ~/path_to/proto_libs/ --decode=package.ProtoMessage ~/path_to/message.proto; \
      echo ;\
    done
```

### Producer

```
  $ cat ~/proto_data | protoc -I ~/path_to/proto_libs/ --encode=package.ProtoMessage ~/path_to/message.proto | \ 
  kafka-tools produce -t Topic -H header=value -b 127.1.1.2:9092
```
