:a
aws ec2 start-instances --instance-ids i-01127dad19a1dd53d
IF %ERRORLEVEL% NEQ 0 ( 
   goto a 
)