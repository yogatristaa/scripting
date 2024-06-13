#!/bin/bash

mkdir $WORKSPACE/.fh/.tmpr

echo "Pipeline $JOB_NAME $STATUS" >> $WORKSPACE/.fh/.tmpr/msg2slack.txt
if [ -n "$targetbranch" ] && [ "$targetbranch" == "main" ]; then
    echo "Commiter : $comitter" >> $WORKSPACE/.fh/.tmpr/msg2slack.txt
    echo "Merged By : $mergedby" >> $WORKSPACE/.fh/.tmpr/msg2slack.txt
else
    echo "Trigger By : $author" >> $WORKSPACE/.fh/.tmpr/msg2slack.txt
fi

msgText="Pipeline Logs"
messagePayload="[<$BUILD_URL|$msgText>]"
echo "$messagePayload" >> $WORKSPACE/.fh/.tmpr/msg2slack.txt