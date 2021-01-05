package constants

type ApiResponseCode uint32

const (
	SUCCESS                           ApiResponseCode = 1000
	JINJA2_RENDER_FAILURE             ApiResponseCode = 1001
	GET_FILE_FAILURE                  ApiResponseCode = 1002
	COMMAND_DETACHED_START_FAILURE    ApiResponseCode = 1003
	FOLDER_ZIP_FAILURE                ApiResponseCode = 1005
	GET_FILE_FAILURE_IS_DIR           ApiResponseCode = 1006
	GET_ENV_VAR_FAILURE               ApiResponseCode = 1007
	MISSING_PARAMETER_POST            ApiResponseCode = 1008
	GET_COMMAND_DETACHED_INFO_FAILURE ApiResponseCode = 1009
	EMPTY_REQUEST_BODY_PROVIDED       ApiResponseCode = 1010
	TEST_STOP_FAILURE                 ApiResponseCode = 1011
	UPLOAD_FILE_FAILURE               ApiResponseCode = 1012
	HTTP_HEADER_NOT_PROVIDED          ApiResponseCode = 1013
	COMMAND_EXEC_FAILURE              ApiResponseCode = 1014
	EXEC_COMMAND_NOT_ALLOWED          ApiResponseCode = 1015
	UNAUTHORIZED                      ApiResponseCode = 1016
	SET_ENV_VAR_FAILURE               ApiResponseCode = 1017
	INVALID_JSON_PAYLOAD              ApiResponseCode = 1018
	NOT_IMPLEMENTED                   ApiResponseCode = 1019
	GENERAL                           ApiResponseCode = 1100
)
