package errors

func ErrorCommand(err error, exitCode int) (string, string) {
	switch exitCode {
	case 1:
		return "error", "Abnormal end"
	case 17:
		return "error", "The parameter is incorrect."
	case 18:
		return "error", "Failed to open the file."
	case 19:
		return "error", "Failed to read the file."
	case 20:
		return "error", "Failed to write the file."
	case 21:
		return "error", "Failed to get information about the file."
	case 22:
		return "error", "Failed to rename the file."
	case 23:
		return "error", "Failed to delete theun file."
	case 24:
		return "error", "Failed to copy the file."
	case 25:
		return "error", "Failed to seek the file."
	case 26:
		return "error", "Failed to get the file position."
	case 27:
		return "error", "Failed to allocate memory."
	case 28:
		return "error", "Failed to encode to BASE64."
	case 29:
		return "error", "Failed to decode from BASE64."
	case 30:
		return "error", "Failed to encrypt to AES."
	case 31:
		return "error", "Failed to decrypt from AES."
	case 32:
		return "error", "Failed to encrypt to RSA."
	case 33:
		return "error", "Failed to decrypt from RSA."
	case 34:
		return "error", "Failed to deflate to ZIP."
	case 35:
		return "error", "Failed to inflate from ZIP."
	case 36:
		return "error", "Failed to throw exception."
	case 37:
		return "error", "Failed to generate RSA key."
	case 38:
		return "error", "Failed to encode to ARS."
	case 39:
		return "error", "Failed to decode from ARS."
	case 40:
		return "error", "There is no option specified."
	case 41:
		return "error", "The specified option is duplicated."
	case 42:
		return "error", "The value of the option is invalid."
	case 43:
		return "error", "You do not have permission to open this file."
	case 44:
		return "error", "The GUID of the file is different."
	case 45:
		return "error", "Failed to get the current time."
	case 46:
		return "error", "Failed to get the result of HTTP communication."
	case 47:
		return "error", "The result of HTTP communication is invalid."
	case 48:
		return "error", "Failed to initialize OpenSSL."
	case 49:
		return "error", "Failed to release OpenSSL."
	case 50:
		return "error", "The input file size is invalid."
	case 51:
		return "error", "Failed to encrypt to FCL file."
	case 52:
		return "error", "Failed to decrypt from FCL file."
	case 53:
		return "error", "Failed to initialize HTTP connection."
	case 54:
		return "error", "HTTP connection failed."
	case 55:
		return "error", "Failed to read HTTP connection."
	case 56:
		return "error", "Failed to write HTTP connection."
	case 57:
		return "error", "Failed to generate SSL."
	case 58:
		return "error", "Failed to set SSL."
	case 59:
		return "error", "Failed to shut down SSL."
	case 60:
		return "error", "Failed to get HTTP communication response."
	case 61:
		return "error", "HTTP response is invalid."
	case 62:
		return "error", "The format of the FCL file is invalid."
	case 63:
		return "error", "The result of HTTP communication is [NG]. Unknown error code was returned to API from a new version of the server."
	case 64:
		return "error", "The result of HTTP communication is [DEVICE_NOT_FOUND]. Device ID is invalid. This device is not registered on the FinalCode Server."
	case 65:
		return "error", "The result of HTTP communication is [ERR_EXPIRE_TRIAL]. FC server trial period has expired."
	case 66:
		return "error", "The result of HTTP communication is [ERR_EXPIRE_DATE]. License for the server (product version) has expired."
	case 67:
		return "error", "The result of HTTP communication is [INVALID_TOKEN]. User authentication token has expired or the token string is invalid."
	case 68:
		return "error", "The result of HTTP communication is [INVALID_PRESET]. Specified template ID does not exist, or the user specified in [-user:] is not authorized to use this template."
	case 69:
		return "error", "The result of HTTP communication is [INVALID_WATERMARK]. Specified watermark ID does not exist, or the user specified in [-user:] is not authorized to use this watermark."
	case 70:
		return "error", "The result of HTTP communication is [LICENSE_DELETED]. FC Server license has been deleted."
	case 71:
		return "error", "The result of HTTP communication is [LIMIT_OVER]. Exceeded the maximum number of licensed users that can be registered on FC server."
	case 72:
		return "error", "The result of HTTP communication is [NOT_LICENSED]. Invalid operation by Viewer user or user from another company."
	case 73:
		return "error", "The result of HTTP communication is [NOT_REGISTERED]."
	case 74:
		return "error", "The result of HTTP communication is [NOT_SENDMAIL]. Could not send email."
	case 75:
		return "error", "The result of HTTP communication is [USER_NOT_FOUND]. Specified user is not registered on FinalCode server."
	case 76:
		return "error", "The result of HTTP communication is [USER_NOT_LICENSED]. Encrypt by specifying Viewer user. * Not used in API 2.0"
	case 77:
		return "error", "The result of HTTP communication is [NOT_FOUND]. File ID does not exist on FinalCode server. This file may have been created on a different [FinalCode server]."
	case 78:
		return "error", "The result of HTTP communication is [WATERMARK_NOT_FOUND]. Watermark not found."
	case 79:
		return "error", "Failed to create password."
	case 80:
		return "error", "The result of HTTP communication is [DUPLICATE_META]. Duplicate meta data."
	case 81:
		return "error", "The result of HTTP communication is [INVALID_OWNER]. Specified owner is invalid."
	case 82:
		return "error", "The result of HTTP communication is [INVALID_EXTENSION]. Extension is invalid or cannot be changed."
	case 83:
		return "error", "The result of HTTP communication is [NOT_LICENSED_DOMAIN_USER]. Not a domain user."
	case 84:
		return "error", "The result of HTTP communication is [NOT_LICENSED_USERS_OPTION]. Not authorized to access user function options."
	case 85:
		return "error", "The result of HTTP communication is [POLICY_NOT_FOUND]. No policies are assigned."
	case 86:
		return "error", "The result of HTTP communication is [INVALID_USER]. User specified in [-user:] is managed under a different license."
	case 87:
		return "error", "The result of HTTP communication is [INVALID_ROLE]. User specified in [-user:] is not authorized to perform operation. * This error occurs when unauthorized users try to use commands like [-encrypt]."
	case 88:
		return "error", "The result of HTTP communication is [INVALID_POLICY]."
	case 89:
		return "error", "The result of HTTP communication is [OVER_PRESET]. Reached the maximum number of templates."
	case 90:
		return "error", "The result of HTTP communication is [DUPLICATE_PRESET]. Duplicate templates."
	case 91:
		return "error", "Reached the minimum number of policies."
	case 92:
		return "error", "Reached the maximum number of recipients."
	case 93:
		return "error", "Operation not supported by the server."
	case 94:
		return "error", "The result of HTTP communication is [IS_USE_PRESET]. Individual policies cannot be specified on policies set in the template by [-update_file_info_ex]."
	case 95:
		return "error", "Failed to initialize Fips mode."
	case 96:
		return "error", "Failed to compress folder."
	case 97:
		return "error", "Failed to decompress folder."
	case 98:
		return "error", "Failed to search file."
	case 99:
		return "error", "Failed to delete folder."
	case 100:
		return "error", "The result of HTTP communication is [INVALID_TAG]. The value of the tag is invalid."
	case 101:
		return "error", "The result of HTTP communication is [INVALID_BV_AUTH]. The browser view authentication method is invalid."
	case 102:
		return "error", "Browser view file format is invalid."
	case 103:
		return "error", "The result of HTTP communication is [DUPLICATE_TAG]. Tags are duplicated."
	case 104:
		return "error", "Invalid format of transparent secure file."
	default:
		return "error", "Unknown error"
	}
}
