package a24apiclient

import (
)

var C_a24apiclient_codes = map[string]map[string]map[int]string {
    "_shared_": map[string]map[int]string {
        "_codes_": map[int]string {
            200: "OK",
            204: "OK",
            401: "TOKEN_INVALID",
            403: "UNAUTHORIZED",
            429: "TOO_MANY_REQUESTS",
            500: "SYSTEM_ERROR",
        },
    },
    "dns": map[string]map[int]string {
        "delete": map[int]string {
            400: "DNS_RECORD_TO_DELETE_NOT_FOUND",
        },
        "update": map[int]string {
            400: "DNS_RECORD_TO_UPDATE_NOT_FOUND",
        },
        "create": map[int]string {
            400: "VALIDATION_ERROR",
        },
    },
    "domains": map[string]map[int]string {
        "detail": map[int]string {
            400: "OBJECT_ID_DOESNT_EXIST",
        },
    },
}
