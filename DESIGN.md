# A3 - gRPC Reddit

- Author - Tomy Hsieh (chiweih)
- Source code: [GitHub](https://github.com/tomy0000000/grpc-reddit)
- Read this document on [GitHub](https://github.com/tomy0000000/grpc-reddit/blob/main/DESIGN.md) for links navigation

## Protocol Buffer Definitions

See source code in [`reddit.proto`](reddit/reddit.proto)

```protobuf
message User {
  int32 id = 1;
}

message SubReddit {
  int32 id = 1;
  string name = 2;
  SubRedditState state = 3; // State should never be unspecified
  repeated string tags = 4;
}

message Post {
  int32 id = 1;
  string title = 2;
  string content = 3;
  SubReddit subReddit = 4;
  optional string videoURL = 5;
  optional string imageURL = 6;
  optional User author = 7;
  int32 score = 8;
  PostState state = 9; // State should never be unspecified
  google.type.Date publicationDate = 10;
}

message Comment {
  int32 id = 1;
  string content = 2;
  User author = 3;
  int32 score = 4;
  CommentState state = 5; // State should never be unspecified
  google.type.Date publicationDate = 6;
  ContentType parent = 7; // Parent should never be unspecified
  int32 parentID = 8;
  repeated Comment children = 9;
}
```

## Service Definitions

See source code in [`reddit.proto`](reddit/reddit.proto)

```protobuf
// The reddit service definition.
service Reddit {
  // Create a post
  rpc CreatePost (CreatePostRequest) returns (CreatePostResponse) {}

  // Upvote or downvote a Post
  rpc VotePost (VotePostRequest) returns (VotePostResponse) {}

  // Retrieve Post content
  rpc GetPost (GetPostRequest) returns (GetPostResponse) {}

  // Create a Comment
  rpc CreateComment (CreateCommentRequest) returns (CreateCommentResponse) {}

  // Upvote or downvote a Comment
  rpc VoteComment (VoteCommentRequest) returns (VoteCommentResponse) {}

  // Retrieve a Comment
  rpc GetComment (GetCommentRequest) returns (GetCommentResponse) {}

  // Retrieving a list of N most upvoted comments under a post
  rpc GetTopComments (GetTopCommentsRequest) returns (GetTopCommentsResponse) {}

  // Expand a comment branch
  rpc ExpandCommentBranch (ExpandCommentBranchRequest) returns (ExpandCommentBranchResponse) {}

  // Monitor updates to posts and comments
  rpc MonitorUpdates (stream MonitorUpdatesRequest) returns (stream MonitorUpdatesResponse) {}
}
```

| Service               | Input                        | Output                        |
| --------------------- | ---------------------------- | ----------------------------- |
| `CreatePost`          | `CreatePostRequest`          | `CreatePostResponse`          |
| `VotePost`            | `VotePostRequest`            | `VotePostResponse`            |
| `GetPost`             | `GetPostRequest`             | `GetPostResponse`             |
| `CreateComment`       | `CreateCommentRequest`       | `CreateCommentResponse`       |
| `VoteComment`         | `VoteCommentRequest`         | `VoteCommentResponse`         |
| `GetComment`          | `GetCommentRequest`          | `GetCommentResponse`          |
| `GetTopComments`      | `GetTopCommentsRequest`      | `GetTopCommentsResponse`      |
| `ExpandCommentBranch` | `ExpandCommentBranchRequest` | `ExpandCommentBranchResponse` |
| `MonitorUpdates`      | `MonitorUpdatesRequest`      | `MonitorUpdatesResponse`      |

## Storage Backend

As defined in the extra credit of implementation, SQLite is used as the storage backend.

An example database is included as [`reddit.db`](data/reddit.db)

## Implementation

* [Server](server/main.go)
* [Client](client/client.go)
* [High level function](client/main.go) and its [test](client/client_test.go)

Video demo is provided in [GitHub README](https://github.com/tomy0000000/grpc-reddit)

## Extra Credit

* Data Model
  * [Subreddit (5pts)](reddit/reddit.proto#L81-L86)
* Service design
  * [Monitor updates (5pts)](reddit/reddit.proto#L35-L36)
* Implementation
  * [Implement the server portion of the extra credit API above (5pts)](server/main.go#L162-L237)
  * [Implement the client portion of the extra credit API above (5pts).](client/client.go#L253-L300)
  * [Implement actual storage for the models using SQLite as a storage backend (10pts)](server/sqlclient.go)
* Testing
  * Postman: See screenshot in [`/img`](img)
