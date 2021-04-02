package crypto

// PKCS7Pad() pads an byte array to be a multiple of 16
// http://tools.ietf.org/html/rfc5652#section-6.3
func PKCS7Pad(data []byte) []byte {
	dataLen := len(data)

	validLen := int(dataLen/16+1) * 16

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
	// Edge case
	if dataLen == 0 {
		return nil
	}
	// the last byte indicates the length of the padding to remove
	paddingLen := int(data[dataLen-1])

	if paddingLen == dataLen {
		return nil
	} else if paddingLen < dataLen {
		return data[:dataLen-paddingLen]
	}
	return data
}
