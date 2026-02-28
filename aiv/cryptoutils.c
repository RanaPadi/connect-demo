#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <openssl/hmac.h>
#include <openssl/sha.h>
#include <openssl/ec.h>
#include <openssl/evp.h>
#include <openssl/bn.h>
#include <openssl/rand.h>
#include <openssl/hmac.h>
// Function to get the BNP256 curve group

EC_GROUP *get_ec_group_bnp256(void){


    unsigned char bnp256_p[32] = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0xF0, 0xCD, 0x46, 0xE5, 0xF2, 0x5E, 0xEE, 0x71, 0xA4,
                                0x9F, 0x0C, 0xDC, 0x65, 0xFB, 0x12, 0x98, 0x0A, 0x82, 0xD3, 0x29, 0x2D, 0xDB, 0xAE, 0xD3,
                                0x30, 0x13};

    unsigned char BNP256_ORDER[32] = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0xF0, 0xCD, 0x46, 0xE5, 0xF2, 0x5E, 0xEE, 0x71,
                                    0xA4, 0x9E, 0x0C, 0xDC, 0x65, 0xFB, 0x12, 0x99, 0x92, 0x1A, 0xF6, 0x2D, 0x53, 0x6C,
                                    0xD1, 0x0B, 0x50, 0x0D};

    int ok = 0;
    EC_GROUP *curve = NULL;
    EC_POINT *generator = NULL;
    BN_CTX *ctx = BN_CTX_new();
    BIGNUM *tmp_1 = NULL, *tmp_2 = NULL, *tmp_3 = NULL;

    unsigned char bnp256_a[1] = {0x00};
    unsigned char bnp256_b[1] = {0x03};
    unsigned char bnp256_gX_[1] = {0x01};
    unsigned char bnp256_gy_[1] = {0x02};


    if ((tmp_1 = BN_bin2bn(bnp256_p, 32, NULL)) == NULL)
        goto err;
    if ((tmp_2 = BN_bin2bn(bnp256_a, 1, NULL)) == NULL)
        goto err;
    if ((tmp_3 = BN_bin2bn(bnp256_b, 1, NULL)) == NULL)
        goto err;
    if ((curve = EC_GROUP_new_curve_GFp(tmp_1, tmp_2, tmp_3, NULL)) == NULL)
        goto err;


    generator = EC_POINT_new(curve);
    if (generator == NULL)
        goto err;
    if ((tmp_1 = BN_bin2bn(bnp256_gX_, 1, tmp_1)) == NULL)
        goto err;
    if ((tmp_2 = BN_bin2bn(bnp256_gy_, 1, tmp_2)) == NULL)
        goto err;
    if (1 != EC_POINT_set_affine_coordinates(curve, generator, tmp_1, tmp_2, ctx))
        goto err;

    if ((tmp_1 = BN_bin2bn(BNP256_ORDER, 32, tmp_1)) == NULL)
        goto err;
    BN_one(tmp_2);
    if (1 != EC_GROUP_set_generator(curve, generator, tmp_1, tmp_2))
        goto err;

    ok = 1;

    err:
    if (tmp_1)
        BN_free(tmp_1);
    if (tmp_2)
        BN_free(tmp_2);
    if (tmp_3)
        BN_free(tmp_3);
    if (generator)
        EC_POINT_free(generator);
    if (ctx)
        BN_CTX_free(ctx);
    if (!ok) {
        printf("[AIV]\t]FAILED TO CREATE BNP256 CURVE\n");
        EC_GROUP_free(curve);
        curve = NULL;
    }
    return (curve);
}

int SET_ECC_POINT(EC_POINT *res, uint8_t *pX, uint8_t *pY, EC_GROUP *ecgrp, int len){
    BIGNUM *x = BN_new();
    BIGNUM *y = BN_new();
    BN_CTX *ctx = BN_CTX_new();

    BN_bin2bn(pX, len, x);
    BN_bin2bn(pY, len, y);
    if (1 != EC_POINT_set_affine_coordinates(ecgrp, res, x, y, ctx)) {
        printf("[AIV]\t FAILED TO CREATE EC POINT FROM X,Y COORDINATES\n");
        BN_free(y);
        BN_free(x);  
        BN_CTX_free(ctx); 
        EC_GROUP_free(ecgrp);
        exit(-1);
    }
    BN_free(y);
    BN_free(x);  
    BN_CTX_free(ctx);
    return 1;
}

int point2bb(uint8_t *X, uint8_t *Y, EC_POINT *resPoint){
	
	int rc;
	EC_GROUP *ecgrp = NULL;
	BIGNUM *bn_X = BN_new();
	BIGNUM *bn_Y = BN_new();
	BN_CTX *ctx = BN_CTX_new();

	ecgrp = get_ec_group_bnp256();

	rc = EC_POINT_get_affine_coordinates(ecgrp , resPoint, bn_X, bn_Y, ctx);
	BN_bn2bin(bn_X, X);
	BN_bn2bin(bn_Y, Y);
    
    EC_GROUP_free(ecgrp);
    BN_free(bn_X);
    BN_free(bn_Y);
    BN_CTX_free(ctx);

	return rc;
}


int VerifySignature(uint8_t *SigR, uint8_t *SigS, uint8_t *message, int message_size, uint8_t *Public){
    BIGNUM *BN_SigS = BN_new();
    BIGNUM *BN_SigR = BN_new();
    BIGNUM *BN_Digest = BN_new();
    BIGNUM *BN_Order = BN_new();
    BIGNUM *BN_rs = BN_new();
    BIGNUM *BN_hs = BN_new();
    BN_CTX *ctx = BN_CTX_new();
    EC_POINT *Public_Point = NULL;
    EC_POINT *Generator = NULL;
    EC_POINT *Temp_Point1 = NULL;
    EC_POINT *Temp_Point2 = NULL;
    EC_POINT *Verification_Point = NULL;
    EC_GROUP *ecgrp = NULL;
    uint8_t *Digest = malloc(32);
    uint8_t *temp_x = malloc(32), *temp_y = malloc(32);
    uint8_t p1[2] = {0x01, 0x02};
    unsigned char BNP256_ORDER[32] = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0xF0, 0xCD, 0x46, 0xE5, 0xF2, 0x5E, 0xEE, 0x71,
                                    0xA4, 0x9E, 0x0C, 0xDC, 0x65, 0xFB, 0x12, 0x99, 0x92, 0x1A, 0xF6, 0x2D, 0x53, 0x6C,
                                    0xD1, 0x0B, 0x50, 0x0D};

    SHA256(message, message_size, Digest);
    BN_SigR = BN_bin2bn(SigR, 32, NULL);
    BN_SigS = BN_bin2bn(SigS, 32, NULL);
    BN_Order = BN_bin2bn(BNP256_ORDER, 32, NULL);
    BN_SigS = BN_mod_inverse(NULL, BN_SigS, BN_Order, ctx);
    if(1 != BN_mod_mul(BN_rs, BN_SigR, BN_SigS, BN_Order, ctx)){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO MULTIPLY BN UNDER MOD\n");
        return 0;
    }
    ecgrp = get_ec_group_bnp256();
    if (ecgrp == NULL){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO BNP256 CURVE\n");
        return 0;
    }
    Public_Point = EC_POINT_new(ecgrp);
    for(int i=0;i<65;i++){
        printf("%02X", Public[i]);
    }
    printf("\n");
    if(1 != EC_POINT_oct2point(ecgrp, Public_Point, Public, 65, ctx)){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO CONVERT BYTE BUFFER TO POINT\n");
        return 0;
    }
    Temp_Point1 = EC_POINT_new(ecgrp);
    if (1 != EC_POINT_mul(ecgrp, Temp_Point1, NULL, Public_Point, BN_rs, ctx)) {
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO MULTIPLY BIGNUM UNDER MOD\n");
        return 0;
    }
    BN_Digest = BN_bin2bn(Digest, 32, NULL);
    if(1 != BN_mod_mul(BN_hs, BN_Digest, BN_SigS, BN_Order, ctx)){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO MULTIPLY BN UNDER MOD\n");
        return 0;
    }
    Generator = EC_POINT_new(ecgrp);
    if(1 != SET_ECC_POINT(Generator, &p1[0], &p1[1], ecgrp, 1)){
        printf("[AIV]\t FAILED TO CREATE BNP256 GENERATOR\n");
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        EC_POINT_free(Generator);
        free(Digest);
        free(temp_x);
        free(temp_y);
        return 0;
    }
    Temp_Point2 = EC_POINT_new(ecgrp);
    if (1 != EC_POINT_mul(ecgrp, Temp_Point2, NULL, Generator, BN_hs, ctx)) {
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        EC_POINT_free(Temp_Point2);
        EC_POINT_free(Generator);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO MULTIPLY BIGNUM UNDER MOD\n");
        return 0;
    }
    Verification_Point = EC_POINT_new(ecgrp);
    if (1 != EC_POINT_add(ecgrp, Verification_Point, Temp_Point1, Temp_Point2, ctx)){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        EC_POINT_free(Temp_Point2);
        EC_POINT_free(Generator);
        EC_POINT_free(Verification_Point);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO MULTIPLY BIGNUM UNDER MOD\n");
        return 0;
    }
    if( 1 != point2bb(temp_x, temp_y, Verification_Point)){
        printf("[AIV]\tFAILED TO CONVERT EC POINT TO BB\n");
 
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        EC_POINT_free(Temp_Point2);
        EC_POINT_free(Generator);
        EC_POINT_free(Verification_Point);
        free(Digest);
        free(temp_x);
        free(temp_y);
        return 0;
    }
    if(memcmp(temp_x, SigR, 32) != 0){
        BN_CTX_free(ctx);
        BN_free(BN_SigS);
        BN_free(BN_SigR);
        BN_free(BN_Digest);
        BN_free(BN_Order);
        BN_free(BN_hs);
        BN_free(BN_rs);
        EC_GROUP_free(ecgrp);
        EC_POINT_free(Public_Point);
        EC_POINT_free(Temp_Point1);
        EC_POINT_free(Temp_Point2);
        EC_POINT_free(Generator);
        EC_POINT_free(Verification_Point);
        free(Digest);
        free(temp_x);
        free(temp_y);
        printf("[AIV]\t FAILED TO VERIFY SIGNATURE\n");
        return 0;
    }
    BN_CTX_free(ctx);
    BN_free(BN_SigS);
    BN_free(BN_SigR);
    BN_free(BN_Digest);
    BN_free(BN_Order);
    BN_free(BN_hs);
    BN_free(BN_rs);
    EC_GROUP_free(ecgrp);
    EC_POINT_free(Public_Point);
    EC_POINT_free(Temp_Point1);
    EC_POINT_free(Temp_Point2);
    EC_POINT_free(Generator);
    EC_POINT_free(Verification_Point);
    free(Digest);
    free(temp_x);
    free(temp_y);
    return 1;
}
