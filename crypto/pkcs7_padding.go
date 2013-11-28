package crypto

// PKCS7Pad() pads an byte array to be a multiple of 16
// http://tools.ietf.org/html/rfc5652#section-6.3
func PKCS7Pad(data []byte) []byte {
	dataLen := len(data)

	var validLen int
	if dataLen%16 == 0 {
		validLen = dataLen
	} else {
		validLen = int(dataLen/16+1) * 16
	}

	paddingLen := validLen - dataLen
	// The length of the padding is used as the byte we will
	// append as a pad.
	bitCode := byte(paddingLen)
	padding := make([]byte, paddingLen)
	for i := 0; i < paddingLen; i++ {
		padding[i] = bitCode
	}
	return append(data, padding...)
}

// PKCS7Unpad() removes any potential PKCS7 padding added.
func PKCS7Unpad(data []byte) []byte {
	dataLen := len(data)
	// the last byte indicates the length of the padding to remove
	paddingLen := int(data[dataLen-1])

	// padding length can only be between 1-15
	if paddingLen < 16 {
		return data[:dataLen-paddingLen]
	}
	return data
}
