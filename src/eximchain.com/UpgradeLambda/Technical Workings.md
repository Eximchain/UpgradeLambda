How the S3BucketWatch and UpgradeLambda works
=============================================

The UpgradeLambda is a Go application that has a SNSHandler function.
An AWS lifecycle hook is attached to an autoscaling group.
When the autoscaling group launches an instances, the lifecycle hook receives the "Lifecycle Launch" notification and sends it to the UpgradeLambda.

The UpgradeLambda receives the notification metadata from the lifecycle hook which consists of these parameters:

* the autoscaling instance's
  * node type (one of bootnodes, quorum-makers, quorum-observers, quorum-validators, vault-servers),
  * id.
* the upgrade bucket id,
* the notification bucket id,
* among other information.

The UpgradeLambda then pulls the upgrade.sh script and encrypted certificate from the upgrade bucket, and decrypts the SSH cert in-memory.
Using a SSH connection, it then pushes upgrade.sh to the instance, and executes it with screen -dm (so as to work around the run-time limitation of 5 minutes on Lambda functions, as imposed by AWS), passing it the node id, node type, upgrade bucket, notification bucket IDs and 

upgrade.sh pulls upgrade.zip from the upgrade bucket, unzips it, and executes upgrade-&lt;node type&gt;.sh
It outputs the instance id, and upgrade topic ARN, so that if

After upgrade.sh finishes, it then uploads the &lt;instance id&gt;-output.log to the notification bucket, which then launches the S3BucketWatch lambda.

The S3BucketWatch lambda will then read &lt;instance id&gt;-output.log from the notification bucket and determine if the upgrade is successful or not.

Miscellaneous details
=====================
UpgradeLambda 