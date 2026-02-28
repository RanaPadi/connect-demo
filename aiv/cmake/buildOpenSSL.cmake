include(ExternalProject)

externalproject_add(OSSL
        URL https://github.com/openssl/openssl/releases/download/openssl-3.1.2/openssl-3.1.2.tar.gz
        URL_MD5 1d7861f969505e67b8677e205afd9ff4
        BUILD_IN_SOURCE 0
        PREFIX ${CMAKE_BINARY_DIR}/deps/OPENSSL
        CONFIGURE_COMMAND ""
        BUILD_COMMAND
        cd <SOURCE_DIR> &&
        rm -f CMakeCache.txt &&
        ./Configure no-threads &&
        make
        INSTALL_COMMAND ""
        )

set(OSSL_INCLUDE_PATH ${CMAKE_BINARY_DIR}/deps/OPENSSL/src/OSSL/include)
set(OSSL_LIB_PATH ${CMAKE_BINARY_DIR}/deps/OPENSSL/src/OSSL)
