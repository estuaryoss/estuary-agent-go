package utils

import "github.com/estuaryoss/estuary-agent-go/src/constants"

func GetMessage() map[uint32]string {
	return map[uint32]string{
		uint32(constants.SUCCESS):                           "Success",
		uint32(constants.JINJA2_RENDER_FAILURE):             "Jinja2 render failed",
		uint32(constants.GET_FILE_FAILURE):                  "Getting file %s from the estuary agent service failed",
		uint32(constants.COMMAND_DETACHED_START_FAILURE):    "Starting command in background with id %s failed",
		uint32(constants.COMMAND_DETACHED_STOP_FAILURE):     "Stopping running detached command failed",
		uint32(constants.GET_FILE_FAILURE_IS_DIR):           "Getting %s from %s failed. It is a directory, not a file.",
		uint32(constants.GET_ENV_VAR_FAILURE):               "Getting env var %s failed.",
		uint32(constants.MISSING_PARAMETER_POST):            "Body parameter \"%s\" sent in request missing. Please include parameter. E.g. {\"parameter\", \"value\"}",
		uint32(constants.FOLDER_ZIP_FAILURE):                "Failed to zip folder %s.",
		uint32(constants.EMPTY_REQUEST_BODY_PROVIDED):       "Empty request body provided.",
		uint32(constants.UPLOAD_FILE_FAILURE):               "Failed to upload file.",
		uint32(constants.HTTP_HEADER_NOT_PROVIDED):          "Http header value not provided, '%s'",
		uint32(constants.COMMAND_EXEC_FAILURE):              "Starting command(s) failed",
		uint32(constants.UNAUTHORIZED):                      "Unauthorized",
		uint32(constants.SET_ENV_VAR_FAILURE):               "Failed to set env vars \"%s\"",
		uint32(constants.INVALID_JSON_PAYLOAD):              "Invalid json body \"%s\"",
		uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE): "Failed to get detached command info.",
		uint32(constants.INVALID_YAML_CONFIG):               "Invalid yaml config",
	}
}
