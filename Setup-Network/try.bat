:a
aws ec2 start-instances --instance-ids i-0107eb0b8116b6393 i-01127dad19a1dd53d i-0fc34c095d281c2e9 i-0df52417f01fa2341 i-0d949bd21218a207c i-0bd92c6b472459c6a i-0109ee191f89d9eb4
IF %ERRORLEVEL% NEQ 0 ( 
   goto a 
)