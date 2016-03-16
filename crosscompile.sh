# git submodule init
# git submodule update

echo "About to cross-compile ev3.mqtt.fabric.golang"
if [ "$1" = "" ]; then
        echo "You need to pass in the version number as the first parameter like ./crosscompile 0.1"
        exit
fi

rm -rf snapshot/*
mkdir snapshot

cp README.md snapshot/

echo "Building Linux amd64"
mkdir snapshot/build-$1_linux_amd64
env GOOS=linux GOARCH=amd64 go build -v -o snapshot/build-$1_linux_amd64/ev3.mqtt.fabric.golang
#tar -zcvf snapshot/serial-port-json-server-$1_linux_amd64.tar.gz snapshot/serial-port-json-server-$1_linux_amd64

echo "" 
echo "Building Linux 386"
mkdir snapshot/build-$1_linux_386
env GOOS=linux GOARCH=386 go build -v -o snapshot/build-$1_linux_386/ev3.mqtt.fabric.golang
#tar -zcvf snapshot/serial-port-json-server-$1_linux_386.tar.gz snapshot/serial-port-json-server-$1_linux_386

echo "" 
echo "Building Linux ARM (Raspi)"
mkdir snapshot/build-$1_linux_arm
env GOOS=linux GOARCH=arm go build -v -o snapshot/build-$1_linux_arm/ev3.mqtt.fabric.golang
#tar -zcvf snapshot/serial-port-json-server-$1_linux_arm.tar.gz snapshot/serial-port-json-server-$1_linux_arm

echo "" 
echo "Building Windows x32"
mkdir snapshot/build-$1_windows_386
env GOOS=windows GOARCH=386 go build -v -o snapshot/build-$1_windows_386/ev3.mqtt.fabric.golang.exe
#cd snapshot/serial-port-json-server-$1_windows_386
#zip -r ../serial-port-json-server-$1_windows_386.zip *
#cd ../..

echo "" 
echo "Building Windows x64"
mkdir snapshot/build-$1_windows_amd64
env GOOS=windows GOARCH=amd64 go build -v -o snapshot/build-$1_windows_amd64/ev3.mqtt.fabric.golang.exe
#cd snapshot/serial-port-json-server-$1_windows_amd64
#zip -r ../serial-port-json-server-$1_windows_amd64.zip *
#cd ../..

echo "" 
echo "Building Darwin x64"
mkdir snapshot/build-$1_darwin_amd64
env GOOS=darwin GOARCH=amd64 go build -v -o snapshot/build-$1_darwin_amd64/ev3.mqtt.fabric.golang
#cd snapshot/serial-port-json-server-$1_darwin_amd64
#zip -r ../serial-port-json-server-$1_darwin_amd64.zip *
#cd ../..
