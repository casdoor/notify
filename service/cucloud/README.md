# CuCloud

[CUCloud](https://www.cucloud.cn) is the cloud service brand of China Unicom

## Usage

1. Create a pair of accessKey and accessKeySecret in `https://console.cucloud.cn/console/uiam/user/${your_id}/info`
2. Create a topic

```go
// define accessKey and  accessKeySecret
accessKey := ""
secretKey := ""

// create client
cuCloudClient := New(accessKey, secretKey, "topic_name", "message title", "cloud region code", "account id", "notify type")

// send message
err := cuCloudClient.Send(context.Background(), "", "content")
```
