#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>
#include <openssl/ec.h>


int VerifySignature(uint8_t *SigR, uint8_t *SigS, uint8_t *message, int message_size, uint8_t *Public);

EC_GROUP *get_ec_group_bnp256(void);

int SET_ECC_POINT(EC_POINT *res, uint8_t *pX, uint8_t *pY, EC_GROUP *ecgrp, int len);

int point2bb(uint8_t *X, uint8_t *Y, EC_POINT *resPoint);