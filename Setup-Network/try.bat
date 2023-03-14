:a
aws ec2 start-instances --instance-ids i-0107eb0b8116b6393
IF %ERRORLEVEL% NEQ 0 ( 
   goto a 
)