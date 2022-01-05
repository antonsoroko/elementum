#/bin/sh

function test {
  echo "+ $@"
  "$@"
  local status=$?
  if [ $status -ne 0 ]; then
    exit $status
  fi
  return $status
}

GIT_VERSION=`cd ${GOPATH}/src/github.com/elgatito/elementum; git describe --tags`

if [ "$1" == "local" ]
then
  # This will run with local go using libtorrent-go/local-env/ locally copied dependencies compilation.
  export LOCAL_ENV=$GOPATH/src/github.com/ElementumOrg/libtorrent-go/local-env/
  export PATH=$PATH:$LOCAL_ENV/bin/
  export PKG_CONFIG_PATH=$LOCAL_ENV/lib/pkgconfig
  export SWIG_LIB=$LOCAL_ENV/share/swig/4.1.0/

  cd $GOPATH/src/github.com/elgatito/elementum
  set -e
  test go build -ldflags="-w -X github.com/elgatito/elementum/util.Version=${GIT_VERSION}" -o /var/tmp/elementum .
  test chmod +x /var/tmp/elementum
  test cp -rf /var/tmp/elementum $HOME/.kodi/addons/plugin.video.elementum/resources/bin/linux_x64/
  test cp -rf /var/tmp/elementum $HOME/.kodi/userdata/addon_data/plugin.video.elementum/bin/linux_x64/
elif [ "$1" == "sanitize" ]
then
  # This will run with local go
  cd $GOPATH
  set -e
  CGO_ENABLED=1 CGO_LDFLAGS='-fsanitize=leak -fsanitize=address' CGO_CFLAGS='-fsanitize=leak -fsanitize=address' test go build -ldflags="-w -X github.com/elgatito/elementum/util.Version=${GIT_VERSION}" -o /var/tmp/elementum github.com/elgatito/elementum
  test chmod +x /var/tmp/elementum
  test cp -rf /var/tmp/elementum $HOME/.kodi/addons/plugin.video.elementum/resources/bin/linux_x64/
  test cp -rf /var/tmp/elementum $HOME/.kodi/userdata/addon_data/plugin.video.elementum/bin/linux_x64/
elif [ "$1" == "docker" ]
then
  # This will run with docker libtorrent:linux-x64 image
  cd $GOPATH/src/github.com/elgatito/elementum
  test make linux-x64
  test cp -rf build/linux_x64/elementum $HOME/.kodi/addons/plugin.video.elementum/resources/bin/linux_x64/
  test cp -rf build/linux_x64/elementum $HOME/.kodi/userdata/addon_data/plugin.video.elementum/bin/linux_x64/
fi
