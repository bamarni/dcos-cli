/*
 * DC/OS
 *
 * DC/OS API
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dcos

type CosmosServiceUpdateV1Request struct {
	AppId          string                            `json:"appId"`
	Options        map[string]map[string]interface{} `json:"options,omitempty"`
	PackageName    string                            `json:"packageName,omitempty"`
	PackageVersion string                            `json:"packageVersion,omitempty"`
	// If true any stored configuration will be ignored when producing the updated service configuration.
	Replace bool `json:"replace"`
}
