package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"testing"
	"time"

	"softwareupgrade"

	"github.com/aws/aws-lambda-go/events"
)

func createNotificationInfo(buckenName string) (result map[string]string) {
	result = make(map[string]string)
	notificationInfo := result

	notificationInfo[CBucketName] = buckenName
	notificationInfo[CBucketKey] = "enc-ssh"
	notificationInfo[CSSHUser] = "ubuntu"
	notificationInfo[CVPC] = "somevpc"
	notificationInfo[CSecurityGroup] = "somegroup"
	notificationInfo[CUpgradeCmd] = "upgrade.sh"
	notificationInfo[CRemoteUpgradeLocation] = "/home/ubuntu"

	return
}

func TestSSHExec(t *testing.T) {
	var User string
	if currentUser, err := user.Current(); err == nil && User == "" {
		User = currentUser.Username
	} else {
		t.Fatal("Unable to retrieve current user")
	}
	sshConfig := softwareupgrade.NewSSHConfig("ubuntu", "~/.ssh/quorum", "54.211.100.189")
	now := time.Now().String()
	tmp, err := ioutil.TempFile("", "")
	if err != nil || tmp == nil {
		t.Fatal("Unable to get temporary file")
	}
	Filename := "/tmp/testfile"
	defer os.Remove(Filename)
	cmd := fmt.Sprintf(`echo -n "%s" > %s`, now, Filename)
	if _, err := sshConfig.Run(cmd); err == nil {
		var byteContent []byte
		byteContent, err = ioutil.ReadFile(Filename)
		if bytes.Compare(byteContent, []byte(now)) != 0 {
			t.Fatal("Run didn't execute remotely!")
		}
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}

func TestRead(t *testing.T) {
	sshConfig := softwareupgrade.NewSSHConfig("ubuntu", "~/.ssh/quorum", "54.84.62.165")
	sshConfig.Connect()
	if err := sshConfig.CopyRemoteFileToLocalFile("/home/ubuntu/upgrade.log", "~/upgrade.log", "0000"); err != nil {
		fmt.Printf("err: %v", err)
	}
}

// In order to test this, the instance ID and the bucket name below needs to be replaced with the actual ones generated
func TestSNSHandler(t *testing.T) {
	strconv.ParseInt("123456", 10, 64)
	notificationInfo := createNotificationInfo("eximchain105-20180927011346168000000001")
	msg := AWSMessage{EC2InstanceID: "i-0391afd14eab3959c"} // instance ID
	ctx := context.Background()
	event := events.SNSEvent{Records: make([]events.SNSEventRecord, 1)}
	metadata, _ := json.Marshal(&notificationInfo)
	msg.NotificationMetadata = string(metadata)
	marshaledInfo, _ := json.Marshal(&msg)
	snsEntity := events.SNSEntity{Message: string(marshaledInfo)}
	event.Records[0].SNS = snsEntity
	SNSHandler(ctx, event)
}

func TestIntegration(t *testing.T) {
	notificationInfo := createNotificationInfo("eximchain105-20181001024305946100000003")
	notificationInfo[CNotifyBucket] = "terraform-20181001024305945900000002" // Notification Bucket ID
	notificationInfo[CNodeType] = "default"                                  // node type
	msg := AWSMessage{EC2InstanceID: "i-0b80f1e3176b84ef2"}                  // instance ID
	ctx := context.Background()
	event := events.SNSEvent{Records: make([]events.SNSEventRecord, 1)}
	metadata, _ := json.Marshal(&notificationInfo)
	msg.NotificationMetadata = string(metadata)
	marshaledInfo, _ := json.Marshal(&msg)
	snsEntity := events.SNSEntity{Message: string(marshaledInfo)}
	event.Records[0].SNS = snsEntity
	SNSHandler(ctx, event)
}
