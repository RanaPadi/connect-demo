cmake_minimum_required(VERSION 3.16)
project(VPE)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -Wno-sign-compare -O3 -Wno-comment -Wno-unused-function -Wno-write-strings -fexceptions -fpermissive -no-pie")

# Define necessary macros
add_definitions(-DTPM_POSIX)
add_definitions(-DVERBOSE)
add_definitions(-DTPM_TPM20)

# OpenSSL paths
if(NOT DEFINED OSSL_INCLUDE_PATH OR NOT DEFINED OSSL_LIB_PATH)
    include(ExternalProject)

    ExternalProject_Add(OSSL
        URL https://github.com/openssl/openssl/releases/download/openssl-3.1.2/openssl-3.1.2.tar.gz
        URL_MD5 1d7861f969505e67b8677e205afd9ff4
        BUILD_IN_SOURCE 0
        PREFIX ${CMAKE_BINARY_DIR}/deps/OPENSSL
        CONFIGURE_COMMAND ""
        BUILD_COMMAND cd <SOURCE_DIR> && ./Configure no-threads && make
        INSTALL_COMMAND ""
        BUILD_BYPRODUCTS ${CMAKE_BINARY_DIR}/deps/OPENSSL/src/OSSL/include/openssl/ssl.h
    )

    set(OSSL_INCLUDE_PATH ${CMAKE_BINARY_DIR}/deps/OPENSSL/src/OSSL/include)
    set(OSSL_LIB_PATH ${CMAKE_BINARY_DIR}/deps/OPENSSL/src/OSSL)
endif()

# Include and link directories
include_directories(${OSSL_INCLUDE_PATH})
link_directories(${OSSL_LIB_PATH})

# Status messages for debugging
message(STATUS "OPENSSL_INCLUDE_PATH: ${OSSL_INCLUDE_PATH}")
message(STATUS "OPENSSL_LIB_PATH: ${OSSL_LIB_PATH}")

# Find OpenSSL libraries
find_library(OSSL_CRYPTO NAMES crypto HINTS ${OSSL_LIB_PATH})
find_library(OSSL_SSL NAMES ssl HINTS ${OSSL_LIB_PATH})

# Status messages for debugging
message(STATUS "OSSL_CRYPTO: ${OSSL_CRYPTO}")
message(STATUS "OSSL_SSL: ${OSSL_SSL}")

# Check if OpenSSL libraries were found
if(NOT OSSL_CRYPTO OR NOT OSSL_SSL)
    message(FATAL_ERROR "Could not find OpenSSL libraries. Please ensure OSSL_INCLUDE_PATH and OSSL_LIB_PATH are set correctly.")
endif()

# Add the shared library
add_library(aiv_crypto SHARED cryptoutils.c cryptoutils.h)

# Link the OpenSSL libraries
target_link_libraries(aiv_crypto ${OSSL_CRYPTO} ${OSSL_SSL})

# Ensure the OpenSSL build completes before building the library
add_dependencies(aiv_crypto OSSL)
