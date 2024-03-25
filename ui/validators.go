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

package ui

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func ValidateIp(input string) (err error) {
	ip := net.ParseIP(input)
	if ip == nil {
		return errors.New("invalid ip")
	}
	if ip.To4() == nil && ip.To16() == nil {
		return errors.New("invalid ip")
	}

	return
}

func ValidateUrl(input string) (err error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return err
	}

	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") ||
		parsedURL.Host == "" {
		return errors.New("invalid url")
	}

	return
}

func ValidateTcpUdpPort(input string) (err error) {
	_, err = ValidateParseTcpUdpPort(input)
	return
}

func ValidateParseTcpUdpPort(input string) (port uint16, err error) {
	p, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	port = uint16(p)
	if p != int(port) || p < 1 {
		err = errors.New("out of range")
	}

	return
}

func ValidateHexString(input string) (err error) {
	_, err = hex.DecodeString(input)
	return
}

func ValidateFixedLenHexString(input string, length int) (err error) {
	err = ValidateHexString(input)
	if err == nil && len(input) != length {
		err = errors.New("length mismatch")
	}
	return
}

func ValidateIPv4SubnetMask(input string) (err error) {

	parts := strings.Split(input, ".")
	if len(parts) != 4 {
		return errors.New("IPv4 format. <111.222.123.212>")
	}

	var zeroDetected bool = false

	for _, partStr := range parts {
		part, err := strconv.Atoi(partStr)
		if err != nil || part < 0 || part > 255 {
			return errors.New("IPv4 octet out of range 0 .. 255")
		}

		binary := fmt.Sprintf("%08b", part)

		for i := 0; i < len(binary); i++ {
			bit := binary[i]
			if bit == '1' {
				if zeroDetected {
					return errors.New("invalid IPv4 subnet mask")
				}
			} else if bit == '0' {
				zeroDetected = true
			}
		}
	}

	return
}

func ValidateNumber(input string) (err error) {
	regex := regexp.MustCompile("^[0-9]+$")
	if regex.MatchString(input) {
		return nil
	}

	return errors.New("string is not a number")
}

func ValidateNotEmpty(input string) (err error) {
	if input == "" {
		return errors.New("empty string")
	}

	return
}

func ValidateMail(input string) (err error) {
	_, err = mail.ParseAddress(input)

	return
}
