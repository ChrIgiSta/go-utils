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

import "testing"

func TestUiValidatorsValidateIp(t *testing.T) {
	validIpv4 := "111.22.31.254"
	invalidIpv4 := "1.0.0.256"

	if err := ValidateIp(validIpv4); err != nil {
		t.Error("validate valid ipv4 as invalid", err)
	}
	if ValidateIp(invalidIpv4) == nil {
		t.Error("validate invalid ipv4 as valid")
	}

	validIpv6 := "2002:aa12:3222:abcd:affe:cafe:adde:1"
	invalidIpv6 := "2003:aa12:3222:abcd:affe:cafe:adde"

	if err := ValidateIp(validIpv6); err != nil {
		t.Error("validate valid ipv6 as invalid", err)
	}
	if ValidateIp(invalidIpv6) == nil {
		t.Error("validate invalid ipv6 as valid")
	}
}

func TestUiValidatorsValidateUrl(t *testing.T) {
	validUrl := "https://www.myorg.org"
	invalidUrl := "http://kaiser.a/mypath"

	if err := ValidateUrl(validUrl); err != nil {
		t.Error("validate valid url as invalid", err)
	}
	if err := ValidateUrl(invalidUrl); err == nil {
		t.Error("validate invalid url as valid", err)
	}
}

func TestUiValidatorsValidateTcpUdpPort(t *testing.T) {
	validPort := "65535"
	invalidPort := "65536"
	// invalidPort2 := "0"

	if err := ValidateTcpUdpPort(validPort); err != nil {
		t.Error("validate valid port as invalid", err)
	}
	if err := ValidateTcpUdpPort(invalidPort); err == nil {
		t.Error("validate invalid port as valid", err)
	}
}

func TestUiValidatorsValidateHexString(t *testing.T) {
	validHexString := "abcdef0123456789"
	invalidHexString := "abcdef01234567890"

	if err := ValidateHexString(validHexString); err != nil {
		t.Error("validate valid hex string as invalid", err)
	}
	if err := ValidateHexString(invalidHexString); err == nil {
		t.Error("validate invalid hex string as valid", err)
	}
}

func TestUiValidatorsIPv4SubnetMask(t *testing.T) {
	validMask := "255.255.0.0"
	invalidMask := "255.0.255.0"

	if err := ValidateIPv4SubnetMask(validMask); err != nil {
		t.Error("validate valid subnetmask as invalid", err)
	}
	if err := ValidateIPv4SubnetMask(invalidMask); err == nil {
		t.Error("validate invalid subnetmask as valid", err)
	}
}

func TestUiValidatorsValidateNumber(t *testing.T) {
	validNumber := "12345678901"
	invalidNumber := "12345678a"

	if err := ValidateNumber(validNumber); err != nil {
		t.Error("validate valid number as invalid", err)
	}
	if err := ValidateNumber(invalidNumber); err == nil {
		t.Error("validate invalid number as valid", err)
	}
}

func TestUiValidatorsValidateNotEmpty(t *testing.T) {
	emptyString := ""
	notemptyString := "not Empty ;)"

	if err := ValidateNotEmpty(notemptyString); err != nil {
		t.Error("validate not empty string as empty", err)
	}
	if err := ValidateNotEmpty(emptyString); err == nil {
		t.Error("validate empty string as not empty", err)
	}
}

func TestUiValidatorsValidateMail(t *testing.T) {
	validMail := "my.org@myorg.org"
	invalidMail := "a@b.c"

	if err := ValidateMail(validMail); err != nil {
		t.Error("validate valid email as invalid", err)
	}
	if err := ValidateMail(invalidMail); err == nil {
		t.Error("validate invalid email as valid", err)
	}
}
