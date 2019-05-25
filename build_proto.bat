@echo off
set work_path=.\config\pb
cd %work_path% 
for /R %%s in (*.proto) do (
echo %%s
protoc --gofast_out=..\..\pb %%~ns.proto
)
cd ..\..\
