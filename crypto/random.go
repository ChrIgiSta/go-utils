/**
 * Copyright © 2024, Staufi Tech - Switzerland
 * All rights reserved.
 *
 *   ________________________   ___ _     ________________  _  ____
 *  / _____  _  ____________/  / __|_|   /_______________  | | ___/
 * ( (____ _| |_ _____ _   _ _| |__ _      | |_____  ____| |_|_
 *  \____ (_   _|____ | | | (_   __) |     | | ___ |/ ___)  _  \
 *  _____) )| |_/ ___ | |_| | | |  | |     | | ____( (___| | | |
 * (______/  \__)_____|____/  |_|  |_|     |_|_____)\____)_| |_|
 *
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package crypto

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"time"
)

func RandCharacters(length int, capital bool) ([]byte, error) {
	start := int8('a')
	stop := int8('z')

	bytes := make([]byte, length)

	if capital {
		start = int8('A')
		stop = int8('Z')
	}

	for i := 0; i < length; i++ {
		b, err := RandInt8(start, stop)
		if err != nil {
			return nil, err
		}
		bytes[i] = byte(b)
	}

	return bytes, nil
}

func RandString(length int) (string, error) {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		b, err := RandInt8(int8('!'), int8('~'))
		if err != nil {
			return "", err
		}
		bytes[i] = byte(b)
	}

	return string(bytes), nil
}

func RandomGenerator() (rng *rand.Rand) {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng = rand.New(source)
	return
}

func RandInt64(min int64, max int64) (int64, error) {

	valRange := math.Abs(float64(max - min))

	val, err := crand.Int(crand.Reader, big.NewInt(int64(valRange)))

	return min + val.Int64(), err
}

func RandInt32(min int32, max int32) (int32, error) {
	i64, err := RandInt64(int64(min), int64(max))
	return int32(i64), err
}

func RandInt16(min int16, max int16) (int16, error) {
	i64, err := RandInt64(int64(min), int64(max))
	return int16(i64), err
}

func RandInt8(min int8, max int8) (int8, error) {
	i64, err := RandInt64(int64(min), int64(max))
	return int8(i64), err
}

func RandBool() (bool, error) {
	i64, err := RandInt64(0, 2)
	if i64 >= 1 {
		return true, err
	}
	return false, err
}

func RandFloat(max float64) float64 {

	return RandomGenerator().Float64() * max
}
