:: install putty from https://www.chiark.greenend.org.uk/~sgtatham/putty/latest.html
go clean
"D:\putty\pscp.exe" -pw SZLianXiSpingLaiFA3Z97 -r %cd% root@#ip#:/root/GoWorkTest/src/zhangmai_micro
"D:\putty\plink.exe" -batch -pw SZLianXiSpingLaiFA3Z97 root@#ip# "cd /root/GoWorkTest/src/zhangmai_micro/taskservice/;" "sh redeploy97.sh"
pause
