package clcv1

import (
	"fmt"
)

/*
 * Operating System Identifiers
 * https://www.ctl.io/api-docs/v1/#server-server-object
 */
type OperatingSystem int

func (o OperatingSystem) String() string {
	switch o {
	case 2:  return "Windows 2003 32-bit"
	case 3:  return "Windows 2003 64-bit"
	case 4:  return "Windows 2008 32-bit"
	case 5:  return "Windows 2008 64-bit"
	case 6:  return "CentOS 32-bit"
	case 7:  return "CentOS 64-bit"
	case 13: return "FreeBSD 32-bit"
	case 14: return "FreeBSD 64-bit"
	case 15: return "Windows 2003 Enterprise 32-bit"
	case 16: return "Windows 2003 Enterprise 64-bit"
	case 17: return "Windows 2008 Enterprise 32-bit"
	case 18: return "Windows 2008 Enterprise 64-bit"
	case 19: return "Ubuntu 32-bit"
	case 20: return "Ubuntu 64-bit"
	case 21: return "Debian 64-bit"
	case 22: return "RedHat Enterprise Linux 64-bit"
	case 24: return "Windows 2012 64-bit"
	case 25: return "RedHat Enterprise Linux 5 64-bit"
	case 26: return "Windows 2008 Datacenter 64-bit"
	case 27: return "Windows 2012 Datacenter 64-bit"
	case 28: return "Windows 2012 R2 Datacenter 64-Bit"
	case 29: return "Ubuntu 10 32-Bit"
	case 30: return "Ubuntu 10 64-Bit"
	case 31: return "Ubuntu 12 64-Bit"
	case 32: return "CentOS 5 32-Bit"
	case 33: return "CentOS 5 64-Bit"
	case 34: return "CentOS 6 32-Bit"
	case 35: return "CentOS 6 64-Bit"
	case 36: return "Debian 6 64-Bit"
	case 37: return "Debian 7 64-Bit"
	case 38: return "RedHat 6 64-Bit"
	case 39: return "CoreOS"
	case 40: return "PXE Boot"
	case 41: return "Ubuntu 14 64-Bit"
	case 42: return "RedHat 7 64-Bit"
	case 43: return "Windows 2008 R2 Standard 64-Bit"
	case 44: return "Windows 2008 R2 Enterprise 64-Bit"
	case 45: return "Windows 2008 R2 Datacenter 64-Bit"
	default:
		return fmt.Sprintf("Unknown OS %d", o)
	}
}
