## Connection
We will establish two connections between two nodes, as gRPC acts like client-server mode.
For example, we will maintain A->B connection in A's connectedNodes and B->A connection in B's connectedNodes.

```mermaid
sequenceDiagram
    participant A as A (listen at :12001)
    participant B as B (listen at :12000)
    participant C as C (listen at :12002)
    A->>B: dial grpc to :12000
    A->>A: add B (connection to :12000) to connectedNodes
    A->>B: rpc RequestNode()
    B->>A: dial grpc to :12001
    B->>B: add A (connection to :12001) to connectedNodes
    par broadcast nodes
        B->>A: rpc BroadcastNodes()
        B->>C: rpc BroadcastNodes()
    end
    par C connects to A
        C->>A: dial grpc to :12001
        C->>C: add A (connection to :12001) to connectedNodes
    and A connects to C
        A->>C: dial grpc to :12002
        A->>A: add C (connection to :12002) to connectedNodes
    end
```