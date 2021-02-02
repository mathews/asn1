package asn1

import (
	"errors"
	"math/big"
)

func numberOfLeadingZeros(i int) int {
	// HD, Count leading 0's
	if i <= 0 {
		if i == 0 {
			return 3
		}
		return 0
	}
	n := 31 //int
	if i >= 1<<16 {
		n -= 16
		i >>= 16
	}
	if i >= 1<<8 {
		n -= 8
		i >>= 8
	}
	if i >= 1<<4 {
		n -= 4
		i >>= 4
	}
	if i >= 1<<2 {
		n -= 2
		i >>= 2
	}
	return n - (i >> 1)
}

func bitLengthForInt(n int) int {
	return 32 - numberOfLeadingZeros(n)
}

func bitCount(i int) int {
	// HD, Figure 5-2
	i = i - ((i >> 1) & 0x55555555)
	i = (i & 0x33333333) + ((i >> 2) & 0x33333333)
	i = (i + (i >> 4)) & 0x0f0f0f0f
	i = i + (i >> 8)
	i = i + (i >> 16)
	return i & 0x3f
}

// Miscellaneous Bit Operations

/**
 * Returns the number of bits in the minimal two's-complement
 * representation of this BigInteger, <em>excluding</em> a sign bit.
 * For positive BigIntegers, this is equivalent to the number of bits in
 * the ordinary binary representation.  For zero this method returns
 * {@code 0}.  (Computes {@code (ceil(log2(this < 0 ? -this : this+1)))}.)
 *
 * @return number of bits in the minimal two's-complement
 *         representation of this BigInteger, <em>excluding</em> a sign bit.
 */
func bitLength(signum int, mag []int) int {
	n := -1
	m := mag      //int[]
	len := len(m) //int
	if len == 0 {
		n = 0 // offset by one to initialize
	} else {
		// Calculate the bit length of the magnitude
		magBitLength := ((len - 1) << 5) + bitLengthForInt(mag[0]) //int
		if signum < 0 {
			// Check if magnitude is a power of two
			pow2 := (bitCount(mag[0]) == 1) // boolean
			for i := 1; i < len && pow2; i++ {
				pow2 = (mag[i] == 0)
			}
			if pow2 {
				n = magBitLength - 1
			} else {
				n = magBitLength
			}
		} else {
			n = magBitLength
		}
	}
	return n
}

func firstNonzeroIntNum(mag []int) int {
	fn := -2

	// Search for the first nonzero int
	var i int
	mlen := len(mag) //int
	for i = mlen - 1; i >= 0 && mag[i] == 0; i-- {
	}

	fn = mlen - i - 1

	return fn
}

/**
 * Returns the specified int of the little-endian two's complement
 * representation (int 0 is the least significant).  The int number can
 * be arbitrarily high (values are logically preceded by infinitely many
 * sign ints).
 */
func getInt(n, signum int, mag []int) int {
	if n < 0 {
		return 0
	}

	if n >= len(mag) {
		if signum < 0 {

			return 1
		}
		return 0
	}

	magInt := mag[len(mag)-n-1] //int

	var ret int

	if signum >= 0 {
		ret = magInt
	} else if n <= firstNonzeroIntNum(mag) {
		ret = -magInt
	} else {
		ret = -magInt - 1
	}

	return ret
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

/**
 * Takes an array a representing a negative 2's-complement number and
 * returns the minimal (no leading zero bytes) unsigned whose value is -a.
 */
func makePositive(a []byte, off, length int) []int {
	var keep, k int
	indexBound := off + length

	// Find first non-sign (0xff) byte of input
	for keep = off; (keep < indexBound) && (a[keep] == byte(0xff)); keep++ {
	}

	/* Allocate output array.  If all non-sign bytes are 0x00, we must
	 * allocate space for one extra output byte. */
	for k = keep; k < indexBound && a[k] == 0; k++ {
	}

	var extraByte int

	if k == indexBound {
		extraByte = 1
	} else {
		extraByte = 0
	}
	intLength := ((indexBound - keep + extraByte) + 3) >> 2
	result := make([]int, intLength)

	/* Copy one's complement of input into output, leaving extra
	 * byte (if it exists) == 0x00 */
	b := indexBound - 1
	for i := intLength - 1; i >= 0; i-- {
		result[i] = int(a[b]) & 0xff
		b = b - 1
		numBytesToTransfer := Min(3, b-keep+1)
		if numBytesToTransfer < 0 {
			numBytesToTransfer = 0
		}

		for j := 8; j <= 8*numBytesToTransfer; j += 8 {
			result[i] = result[i] | ((int(a[b]) & 0xff) << j)
			b = b - 1
		}

		// Mask indicates which bits must be complemented
		mask := -1 >> (8 * (3 - numBytesToTransfer))
		result[i] = (-result[i] - 1) & mask
	}

	// Add one to one's complement to generate two's complement
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = (int)((result[i] & LONG_MASK) + 1)
		if result[i] != 0 {
			break
		}

	}

	return result
}

func stripLeadingZeroBytes(a []byte, off, len int) []int {
	indexBound := off + len
	var keep int

	// Find first nonzero byte
	for keep = off; keep < indexBound && a[keep] == 0; keep++ {
	}

	// Allocate new array and copy relevant part of input array
	intLength := ((indexBound - keep) + 3) >> 2
	result := make([]int, intLength)
	b := indexBound - 1
	for i := intLength - 1; i >= 0; i-- {
		result[i] = int(a[b]) & 0xff
		b = b - 1
		bytesRemaining := b - keep + 1
		bytesToTransfer := Min(3, bytesRemaining)
		for j := 8; j <= (bytesToTransfer << 3); j += 8 {
			result[i] = result[i] | ((int(a[b]) & 0xff) << j)
			b = b - 1
		}

	}
	return result
}

func checkRange(mag []int) error {
	if len(mag) > MAX_MAG_LENGTH || len(mag) == MAX_MAG_LENGTH && mag[0] < 0 {
		return errors.New("BigInteger would overflow supported range")
	}
	return nil
}

const MAX_MAG_LENGTH int = 0x7fffffff/32 + 1 // (1 << 26)
const LONG_MASK int = 0xffffffff

func makeBigInt(n *big.Int) (encoder, error) {
	if n == nil {
		return nil, StructuralError{"empty integer"}
	}
	var mag []int
	var signum int
	//TODO init mag
	val := n.Bytes()

	if val[0] < 0 {
		mag = makePositive(val, 0, len(val))
		signum = -1
	} else {
		mag = stripLeadingZeroBytes(val, 0, len(val))
		if len(mag) == 0 {
			signum = 0
		}

		signum = 1
	}
	if len(mag) >= MAX_MAG_LENGTH {
		err := checkRange(mag)
		if err != nil {
			return nil, err
		}
	}

	byteLen := bitLength(signum, mag)/8 + 1 //int
	byteArray := make([]byte, byteLen)      //byte[]

	for i, bytesCopied, nextInt, intIndex := byteLen-1, 40, 0, 0; i >= 0; i-- {
		if bytesCopied == 4 {

			nextInt = getInt(intIndex, signum, mag)
			intIndex = intIndex + 1
			bytesCopied = 1
		} else {
			nextInt >>= 8
			bytesCopied++
		}
		byteArray[i] = byte(nextInt)
	}

	return bytesEncoder(byteArray), nil
}