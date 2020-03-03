:: install putty from https://www.chiark.greenend.org.uk/~sgtatham/putty/latest.html
go clean
"C:\Program Files\PuTTY\pscp.exe" -pw #password# -r %cd% root@#ip#:/root/GoWork/src/zhangmai_micro
"C:\Program Files\PuTTY\plink.exe" -batch -pw #password# root@#ip# "cd /root/GoWork/src/zhangmai_micro/apigateway/;" "sh redeploy98.sh"

