/* Encrypter.cpp */

#include "Encryptor.h"

/*
bool AesEncryptor::encrypt(const unsigned char *plaintext, int plaintext_len,
		const unsigned char *key, const unsigned char *iv,
		encrypted_message *output) {
	EVP_CIPHER_CTX *ctx;
	int len;

	if(!(ctx = EVP_CIPHER_CTX_new()))
		return true;

	// Initialize encryption
	if(1 != EVP_EncryptInit_ex(ctx, EVP_aes_256_gcm(), NULL, NULL, NULL))
		return true;

	// Set IV length
	if(1 != EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_IVLEN, 12, NULL))
		return true;

	// Initialize key and IV
	if(1 != EVP_EncryptInit_ex(ctx, NULL, NULL, key, iv))
		return true;

	// Encrypt
	if(1 != EVP_EncryptUpdate(ctx, output->ciphertext, &len, plaintext, plaintext_len))
		return true;
	output->ciphertext_len = len;

	// Finalize
	if(1 != EVP_EncryptFinal_ex(ctx, output->ciphertext + len, &len))
		return true;
	output->ciphertext_len += len;

	// Get authentication tag
	if(1 != EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_GET_TAG, 16, output->tag))
		return true;

	EVP_CIPHER_CTX_free(ctx);
	return false;

}

bool AesEncryptor::decrypt(const unsigned char* ciphertext, int ciphertext_len,
	const unsigned char* key, const unsigned char* iv,
	const unsigned char* tag, unsigned char* plaintext)
{

	EVP_CIPHER_CTX* ctx;
	int len;
	int plaintext_len;
	int ret;

	// Create and initialize context 
	if (!(ctx = EVP_CIPHER_CTX_new()))
	{
		// TODO: throw error
	}

	// Initialize decryption operation
	if (!EVP_DecryptInit_ex(ctx, EVP_aes_256_gcm(), NULL, NULL, NULL))
	{
		EVP_CIPHER_CTX_free(ctx);
		return true;

	}

	// Set IV length 
	if (!EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_IVLEN, 12, NULL))
	{
		EVP_CIPHER_CTX_free(ctx);
		return true;

	}

	// Initialize key and IV
	if (!EVP_DecryptInit_ex(ctx, NULL, NULL, key, iv))
	{
		EVP_CIPHER_CTX_free(ctx);
		return true;

	}

	// Provide the ciphertext to be decrypted 
	if (!EVP_DecryptUpdate(ctx, plaintext, &len, ciphertext, ciphertext_len))
	{
		EVP_CIPHER_CTX_free(ctx);
		return true;

	}
	plaintext_len = len;

	// Set expected authentication tag before finalizing 
	if (!EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_TAG, 16, (void*)tag))
	{
		EVP_CIPHER_CTX_free(ctx);
		return true;

	}

	 // Finalize decryption.
	 //* This will return 0 if the tag verification fails
	ret = EVP_DecryptFinal_ex(ctx, plaintext + len, &len);
	EVP_CIPHER_CTX_free(ctx);

	if (ret > 0)
	{
		//Success - tag verified 
		plaintext_len += len;
		return false;
		//return plaintext_len;

	}
	else
	{
		// Authentication failed - message was tampered with
		return true;

	}

}

// RSA
/Encrypt data using RSA public key
bool RsaEncryptor::encrypt(EVP_PKEY *public_key,
		const unsigned char *plaintext, int plaintext_len,
		unsigned char **ciphertext) 
{
	EVP_PKEY_CTX *ctx;
	size_t ciphertext_len;

	// Create context
	ctx = EVP_PKEY_CTX_new(public_key, NULL);
	if (!ctx) return true;

	// Initialize encryption
	if (EVP_PKEY_encrypt_init(ctx) <= 0) {
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Set padding (OAEP recommended) 
	if (EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_OAEP_PADDING) <= 0) {
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Determine buffer size
	if (EVP_PKEY_encrypt(ctx, NULL, &ciphertext_len, plaintext, plaintext_len) <= 0) {
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Allocate ciphertext buffer
	*ciphertext = (unsigned char*)malloc(ciphertext_len);
	if (!*ciphertext) {
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Perform encryption 
	if (EVP_PKEY_encrypt(ctx, *ciphertext, &ciphertext_len, 
				plaintext, plaintext_len) <= 0) {
		free(*ciphertext);
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	EVP_PKEY_CTX_free(ctx);
	return false;
	//return ciphertext_len;
}

bool RsaEncryptor::decrypt(EVP_PKEY* private_key,
	const unsigned char* ciphertext, int ciphertext_len,
	unsigned char** plaintext)
{
	EVP_PKEY_CTX* ctx;
	size_t plaintext_len;

	ctx = EVP_PKEY_CTX_new(private_key, NULL);
	if (!ctx) return true;

	if (EVP_PKEY_decrypt_init(ctx) <= 0)
	{
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	if (EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_OAEP_PADDING) <= 0)
	{
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Determine buffer size
	if (EVP_PKEY_decrypt(ctx, NULL, &plaintext_len, ciphertext, ciphertext_len) <= 0)
	{
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	*plaintext = (unsigned char *)malloc(plaintext_len);
	if (!*plaintext)
	{
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	// Perform decryption 
	if (EVP_PKEY_decrypt(ctx, *plaintext, &plaintext_len,
		ciphertext, ciphertext_len) <= 0)
	{
		free(*plaintext);
		EVP_PKEY_CTX_free(ctx);
		return true;

	}

	EVP_PKEY_CTX_free(ctx);
	return false;
	//return plaintext_len;

}*/



RawByteBuffer XorEncryptor::encrypt(RawByteBuffer& msg)
{
	RawByteBuffer encrypted(msg.size());
	for (size_t i = 0; i < msg.size(); i++)
	{
		encrypted[i] = msg[i] ^ key_[i % key_.size()];
	}
	return encrypted;
}

RawByteBuffer XorEncryptor::decrypt(RawByteBuffer& cipher)
{
	return encrypt(cipher); // XOR is symmetric
}
