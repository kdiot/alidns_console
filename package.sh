# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

alidns_version=`./bin/alidns version`
echo "build version: $alidns_version"

# cross_compiles
make -f ./Makefile.cross-platform

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        target_name="alidns_${alidns_version}_${os}_${arch}"
        target_path="./packages/${target_name}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./alidns_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir -p ${target_path}/bin
            mv ./alidns_${os}_${arch}.exe ${target_path}/bin/alidns.exe
        else
            if [ ! -f "./alidns_${os}_${arch}" ]; then
                continue
            fi
            mkdir -p ${target_path}/bin
            mv ./alidns_${os}_${arch} ${target_path}/bin/alidns
        fi
        mkdir ${target_path}/conf
        cp ../sample/ddns-conf-simple.json ${target_path}/conf/
        cp ../sample/profile-simple.json ${target_path}/conf/

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${target_name}.zip ${target_name}
        else
            tar -zcf ${target_name}.tar.gz ${target_name}
        fi
        rm -rf ${target_name}
        cd ..
    done
done

cd -